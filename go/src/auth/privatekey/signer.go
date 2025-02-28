// Package privatekey provides a Signer implementation using ECDSA private keys.
// This is the simplest approach for signing but requires careful key management.
package privatekey

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/radiustechsystems/sdk/go/src/auth"
	"github.com/radiustechsystems/sdk/go/src/common"
	"github.com/radiustechsystems/sdk/go/src/crypto"
	"github.com/radiustechsystems/sdk/go/src/providers/eth"
)

// Signer implements the auth.Signer interface using an ECDSA private key.
// This is the simplest way to sign messages and transactions, but it requires keeping the private key in memory.
// For production systems with high security requirements, consider using a custom Signer with a hardware
// security module or key management service.
type Signer struct {
	// address is the Radius address associated with this signer
	address common.Address

	// chainID is the network chain ID used for EIP-155 transaction signing
	chainID *big.Int

	// key is the ECDSA private key used for signing operations
	key *ecdsa.PrivateKey

	// signer is the underlying Ethereum signer implementation
	signer eth.Signer
}

// New creates a new Signer with the given private key.
//
// @param key The ECDSA private key to use for signing
// @param client The Radius client used to retrieve the chain ID
// @return A new Signer instance configured with the provided key and chain ID
func New(key *ecdsa.PrivateKey, client auth.SignerClient) *Signer {
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		chainID = new(big.Int)
	}

	return &Signer{
		address: crypto.PubkeyToAddress(key.PublicKey),
		chainID: chainID,
		key:     key,
		signer:  eth.NewEIP155Signer(chainID),
	}
}

// Address implements the Signer interface
// @return The Radius Address associated with the Signer
func (s *Signer) Address() common.Address {
	return s.address
}

// ChainID implements the Signer interface
// @return The Chain ID associated with the Signer
func (s *Signer) ChainID() *big.Int {
	return s.chainID
}

// Hash implements the Signer interface
// @param tx The transaction to hash
// @return The hash of the given transaction
func (s *Signer) Hash(tx *common.Transaction) common.Hash {
	ethTx := tx.EthTransaction()
	ethHash := s.signer.Hash(ethTx)
	return common.NewHash(ethHash.Bytes())
}

// SignMessage implements the Signer interface
// @param msg The message bytes to sign
// @return The signature bytes, or an error if signing fails
func (s *Signer) SignMessage(msg []byte) ([]byte, error) {
	return crypto.Sign(crypto.Keccak256(
		[]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(msg), msg)),
	), s.key)
}

// SignTransaction implements the Signer interface
// @param tx The transaction to sign
// @return The signed transaction, or an error if signing fails
func (s *Signer) SignTransaction(tx *common.Transaction) (*common.SignedTransaction, error) {
	hash := s.Hash(tx)
	sig, err := crypto.Sign(hash.Bytes(), s.key)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	v := new(big.Int).SetBytes([]byte{sig[64] + 27})
	if s.chainID.Sign() != 0 {
		v = v.Add(v, new(big.Int).Mul(s.chainID, big.NewInt(2)))
		v = v.Add(v, big.NewInt(8))
	}

	// Serialize the signed transaction
	ethTx := tx.EthTransaction()
	ethSignedTx, err := ethTx.WithSignature(s.signer, append(sig[:64], sig[64]+27))
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}
	serialized, err := ethSignedTx.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	return &common.SignedTransaction{
		Transaction: tx,
		R:           new(big.Int).SetBytes(sig[:32]),
		S:           new(big.Int).SetBytes(sig[32:64]),
		V:           v,
		Serialized:  serialized,
	}, nil
}
