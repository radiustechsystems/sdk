package radius

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Address = common.Address

func NewAddress(b []byte) Address {
	return common.BytesToAddress(b)
}

func NewAddressFromHex(h string) Address {
	return common.HexToAddress(h)
}

func NewAddressFromPrivateKey(privateKey *ecdsa.PrivateKey) Address {
	if privateKey == nil {
		return Address{}
	}
	return crypto.PubkeyToAddress(privateKey.PublicKey)
}
