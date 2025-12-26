package hash

import "github.com/dchest/siphash"

type Hash interface {
	Compute(data []byte) uint32
}

type SipHash struct {
	seed1 uint64
	seed2 uint64
}

func NewSipHash(seed1, seed2 uint64) SipHash {
	return SipHash{
		seed1: seed1,
		seed2: seed2,
	}
}

func (sh SipHash) Compute(data []byte) uint32 {
	return uint32(siphash.Hash(sh.seed1, sh.seed2, data))
}
