package hyperloglog

import (
	"math"
	"math/bits"

	"github.com/nnurry/probabilistics/v2/utilities/hasher"
)

type ProbCounter struct {
	pMax uint64
	h    hasher.HashGenerator[uint64]
}

func (c *ProbCounter) Add(item []byte) error {
	hashes, _ := c.h.GenerateHash(item, 0, math.MaxUint64, 1)
	p := uint64(bits.TrailingZeros64(hashes[0]) + 1)
	if c.pMax < p {
		c.pMax = p
	}
	return nil
}

func (c *ProbCounter) Cardinality() uint {
	if c.pMax == 0 {
		return 0
	}
	// not zero -> 1 << pMax = 2^pMax
	return 1 << c.pMax
}
