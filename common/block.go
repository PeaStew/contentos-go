package common

import (
	"encoding/binary"
)

// Marshaller ...
type Marshaller interface {
	Marshall() []byte
	Unmarshall([]byte) error
}

// BlockID is a sha256 byte array, the first 2 byte is
// replaced by the block number
type BlockID struct {
	data [32]byte
}

// BlockNum returns the block num
func (bid *BlockID) BlockNum() uint64 {
	return binary.LittleEndian.Uint64(bid.data[:2])
}

// BlockHeader ...
type BlockHeader interface {
	Previous() BlockID
}

// SignedBlockHeader ...
type SignedBlockHeader interface {
	Id() BlockID
}

// SignedBlock ...
type SignedBlock interface {
	BlockHeader
	SignedBlockHeader
	Marshaller
}
