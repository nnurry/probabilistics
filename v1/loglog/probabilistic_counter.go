package loglog

import (
	"math/bits"

	"github.com/nnurry/probabilistics/v1/hasher"
)

type ProbabilisticCounter struct {
	pMax     uint
	hashFunc hasher.HashFunc64Type
}

func HomemadeCountTrailingZeroes(item uint64) (p uint) {
	// dissect item to count trailing zeroes
	for item != 0 {
		if item&1 != 0 {
			// seen 1 -> stop shifting
			break
		}
		// haven't seen 1 -> increment the counter
		p++
		item >>= 1
	}
	return p
}

func NewProbabilisticCounter() *ProbabilisticCounter {
	return &ProbabilisticCounter{0, hasher.GetHashFunc64("murmur3_128")}
}

func (h *ProbabilisticCounter) PMax() uint {
	return h.pMax
}

func (h *ProbabilisticCounter) addHash(item uint64) {
	p := uint(bits.TrailingZeros64(item) + 1)
	// update max trailing zeroes
	if p > h.pMax {
		h.pMax = p
	}
}

func (h *ProbabilisticCounter) Add(item []byte) error {
	hashes, _ := h.hashFunc(item)
	// take first hash as we aren't gonna use this algorithm anyway
	h.addHash(hashes[0])
	return nil
}

func (h *ProbabilisticCounter) Cardinality() uint {
	if h.pMax == 0 {
		return 0
	}
	// not zero -> 1 << pMax = 2^pMax
	return 1 << h.pMax
}
