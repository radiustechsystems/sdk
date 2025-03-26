package radius

import "github.com/ethereum/go-ethereum/common"

type Hash = common.Hash

func NewHash(b []byte) Hash {
	return common.BytesToHash(b)
}

func NewHashFromHex(h string) Hash {
	return common.HexToHash(h)
}
