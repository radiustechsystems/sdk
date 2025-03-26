package test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/radiustechsystems/sdk/go/src/radius"
)

func TestPrivateKeySigner_Address(t *testing.T) {
	privateKey := radius.GeneratePrivateKey()

	t.Run("Returns address derived from private key", func(t *testing.T) {
		signer := radius.NewPrivateKeySigner(privateKey, TestnetChainID)
		signerAddr := signer.Address()
		expectedAddr := radius.NewAddressFromPrivateKey(privateKey)
		assert.Equal(t, expectedAddr, signerAddr, "Address should match the one derived from private key")
	})
}

func TestPrivateKeySigner_Hash(t *testing.T) {
	privateKey := radius.GeneratePrivateKey()

	t.Run("Calculates transaction hash", func(t *testing.T) {
		signer := radius.NewPrivateKeySigner(privateKey, TestnetChainID)
		toAddr := radius.NewAddressFromHex("0x8ba1f109551bD432803012645Ac136ddd64DBA72")
		tx := CreateTestTransaction(toAddr)

		eip155Signer := types.NewEIP155Signer(TestnetChainID)
		expectedHash := eip155Signer.Hash(tx)

		hash := signer.Hash(tx)
		assert.Equal(t, expectedHash, hash, "Transaction hash should match expected value")
	})
}

func TestPrivateKeySigner_Sign(t *testing.T) {
	privateKey := radius.GeneratePrivateKey()

	t.Run("Signs a given message hash", func(t *testing.T) {
		signer := radius.NewPrivateKeySigner(privateKey, TestnetChainID)
		message := []byte("test message")

		signature, err := signer.Sign(message)
		require.NoError(t, err, "Sign should not return an error")
		assert.Equal(t, 65, len(signature), "Signature should be 65 bytes")

		expectedMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
		expectedMessageHash := crypto.Keccak256Hash([]byte(expectedMessage))

		pubKey, err := crypto.SigToPub(expectedMessageHash.Bytes(), signature)
		require.NoError(t, err, "Should be able to recover public key from signature")
		assert.Equal(t, crypto.PubkeyToAddress(*pubKey), crypto.PubkeyToAddress(privateKey.PublicKey),
			"Recovered address should match the signer's address")
	})
}

func TestPrivateKeySigner_SignTx(t *testing.T) {
	privateKey := radius.GeneratePrivateKey()
	address := radius.NewAddressFromPrivateKey(privateKey)

	t.Run("Signs a given transaction", func(t *testing.T) {
		signer := radius.NewPrivateKeySigner(privateKey, TestnetChainID)
		toAddr := radius.NewAddressFromHex("0x8ba1f109551bD432803012645Ac136ddd64DBA72")
		tx := CreateTestTransaction(toAddr)

		signedTx, err := signer.SignTx(tx)
		require.NoError(t, err, "SignTx should not return an error")
		assert.NotNil(t, signedTx, "SignedTransaction should not be nil")

		txV, txR, txS := signedTx.RawSignatureValues()
		assert.NotNil(t, txV, "V value should not be nil")
		assert.NotNil(t, txR, "R value should not be nil")
		assert.NotNil(t, txS, "S value should not be nil")

		expectedVMin := new(big.Int).Add(new(big.Int).Mul(TestnetChainID, big.NewInt(2)), big.NewInt(35))
		expectedVMax := new(big.Int).Add(expectedVMin, big.NewInt(1))
		assert.True(t,
			txV.Cmp(expectedVMin) == 0 || txV.Cmp(expectedVMax) == 0,
			"V value should follow EIP-155 calculation: expected %d or %d, got %d",
			expectedVMin, expectedVMax, txV,
		)

		// Verify the transaction can be recovered to the correct sender
		eip155Signer := types.NewEIP155Signer(TestnetChainID)
		recoveryID := txV.Uint64() - (TestnetChainID.Uint64()*2 + 35)
		rBytes := radius.PadBytes(txR.Bytes(), 32)
		sBytes := radius.PadBytes(txS.Bytes(), 32)

		withSigTx, err := tx.WithSignature(eip155Signer, append(append(rBytes, sBytes...), byte(recoveryID)))
		require.NoError(t, err, "Should be able to create transaction with signature")

		sender, err := types.Sender(eip155Signer, withSigTx)
		require.NoError(t, err, "Should be able to recover sender from signed transaction")

		assert.Equal(t, address, sender, "Recovered sender should match signer's address")
	})

	t.Run("SignTx with zero TestnetChainID", func(t *testing.T) {
		zeroChainID := big.NewInt(0)
		signer := radius.NewPrivateKeySigner(privateKey, zeroChainID)
		toAddr := radius.NewAddressFromHex("0x8ba1f109551bD432803012645Ac136ddd64DBA72")
		tx := CreateTestTransaction(toAddr)

		signedTx, err := signer.SignTx(tx)
		require.NoError(t, err, "SignTx should not return an error")

		txV, _, _ := signedTx.RawSignatureValues()
		assert.NotNil(t, txV, "V value should not be nil")

		expectedVMin := big.NewInt(27)
		expectedVMax := big.NewInt(28)
		assert.True(t,
			txV.Cmp(expectedVMin) == 0 || txV.Cmp(expectedVMax) == 0,
			"V value should be 27 or 28 for pre-EIP-155 signatures",
			expectedVMin, expectedVMax, txV,
		)
	})
}

func TestPrivateKeySigner_VerifySignature(t *testing.T) {
	privateKey := radius.GeneratePrivateKey()

	t.Run("Returns true with valid signature", func(t *testing.T) {
		signer := radius.NewPrivateKeySigner(privateKey, TestnetChainID)
		toAddr := radius.NewAddressFromHex("0x8ba1f109551bD432803012645Ac136ddd64DBA72")
		tx := CreateTestTransaction(toAddr)
		signedTx, err := signer.SignTx(tx)
		require.NoError(t, err, "Should be able to sign transaction")

		isValid, err := signer.VerifySignature(signedTx)
		assert.NoError(t, err, "Should be able to verify signature")
		assert.True(t, isValid, "Signature should be valid")
	})

	t.Run("Returns false with wrong private key", func(t *testing.T) {
		signer := radius.NewPrivateKeySigner(privateKey, TestnetChainID)
		toAddr := radius.NewAddressFromHex("0x8ba1f109551bD432803012645Ac136ddd64DBA72")
		tx := CreateTestTransaction(toAddr)
		signedTx, err := signer.SignTx(tx)
		require.NoError(t, err, "Should be able to sign transaction")

		// Create a new private key and signer
		otherPrivateKey := radius.GeneratePrivateKey()
		otherSigner := radius.NewPrivateKeySigner(otherPrivateKey, TestnetChainID)

		isValid, err := otherSigner.VerifySignature(signedTx)
		assert.NoError(t, err, "Should be able to process verification")
		assert.False(t, isValid, "Signature should be invalid after tampering")
	})

	t.Run("Returns false and error with wrong chain ID", func(t *testing.T) {
		signer := radius.NewPrivateKeySigner(privateKey, TestnetChainID)
		toAddr := radius.NewAddressFromHex("0x8ba1f109551bD432803012645Ac136ddd64DBA72")
		tx := CreateTestTransaction(toAddr)
		signedTx, err := signer.SignTx(tx)
		require.NoError(t, err, "Should be able to sign transaction")

		// Create a new signer with a different chain ID
		otherChainID := big.NewInt(123456)
		otherSigner := radius.NewPrivateKeySigner(privateKey, otherChainID)

		isValid, err := otherSigner.VerifySignature(signedTx)
		assert.Error(t, err, "Should not be able to verify signature with different chain ID")
		assert.False(t, isValid, "Signature should be invalid after tampering")
	})
}
