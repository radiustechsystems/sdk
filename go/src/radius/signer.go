package radius

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type Signer interface {
	Address() Address
	Hash(tx *Transaction) Hash
	Sign(message []byte) ([]byte, error)
	SignTx(tx *Transaction) (*Transaction, error)
	VerifySignature(signedTx *Transaction) (bool, error)
}

type PrivateKeySigner struct {
	chainID    *big.Int
	privateKey *ecdsa.PrivateKey
	signer     types.EIP155Signer
}

func NewPrivateKeySigner(privateKey *ecdsa.PrivateKey, chainID *big.Int) Signer {
	return &PrivateKeySigner{
		chainID:    chainID,
		privateKey: privateKey,
		signer:     types.NewEIP155Signer(chainID),
	}
}

func (s *PrivateKeySigner) Address() Address {
	return NewAddressFromPrivateKey(s.privateKey)
}

func (s *PrivateKeySigner) Hash(tx *Transaction) Hash {
	return NewHash(s.signer.Hash(tx).Bytes())
}

func (s *PrivateKeySigner) Sign(message []byte) ([]byte, error) {
	prefixedMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	return crypto.Sign(crypto.Keccak256([]byte(prefixedMessage)), s.privateKey)
}

func (s *PrivateKeySigner) SignTx(tx *Transaction) (*Transaction, error) {
	hash := s.Hash(tx)
	sig, err := crypto.Sign(hash.Bytes(), s.privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	return tx.WithSignature(s.signer, sig)
}

func (s *PrivateKeySigner) VerifySignature(tx *Transaction) (bool, error) {
	txV, txR, txS := tx.RawSignatureValues()

	if txV == nil || txR == nil || txS == nil {
		return false, fmt.Errorf("missing signature components")
	}

	recoveryID := txV.Uint64() - (s.chainID.Uint64()*2 + 35)
	if recoveryID > 1 {
		return false, fmt.Errorf("invalid recovery ID: %d", recoveryID)
	}

	rBytes := PadBytes(txR.Bytes(), 32)
	sBytes := PadBytes(txS.Bytes(), 32)
	sig := append(append(rBytes, sBytes...), byte(recoveryID))

	withSigTx, err := tx.WithSignature(s.signer, sig)
	if err != nil {
		return false, fmt.Errorf("failed to create transaction with signature: %w", err)
	}

	recoveredAddr, err := types.Sender(s.signer, withSigTx)
	if err != nil {
		return false, fmt.Errorf("failed to recover signer: %w", err)
	}

	expectedAddr := NewAddressFromPrivateKey(s.privateKey)
	return recoveredAddr == expectedAddr, nil
}
