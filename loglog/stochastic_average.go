package loglog

import (
	"math"
	"math/bits"

	"github.com/nnurry/probabilistics/bitcounter"
	"github.com/nnurry/probabilistics/hasher"
)

type StochAvgProbabilisticCounter struct {
	buckets  *bitcounter.SqBitCounter
	kBit     uint
	hashFunc hasher.HashFunc64Type
}

func NewStochAvgProbabilisticCounter(kBit uint, log2CounterRange uint) (*StochAvgProbabilisticCounter, error) {
	// take 3 k-bit -> 8 3-bit combinations -> 2^3
	// take x k-bit -> 2^x k-bit combinations
	counter, err := bitcounter.NewSqBitCounter(1<<kBit, log2CounterRange)
	if err != nil {
		return nil, err
	}
	return &StochAvgProbabilisticCounter{
		counter,
		kBit,
		hasher.GetHashFunc64("murmur3_128"),
	}, nil
}

func (h *StochAvgProbabilisticCounter) addHash(item uint64) error {
	// take 1st k-bit from uint64
	bucketIndex := item >> (64 - h.kBit)
	// take the rest
	item <<= uint64(h.kBit)
	// add to bucket
	p := uint64(bits.TrailingZeros64(item) + 1)

	// within counter limit -> read counter
	if p < ((1 << h.buckets.BitRange()) - 1) {
		pMax, err := h.buckets.Read(uint(bucketIndex))
		if err != nil {
			return err
		}
		// current counter value (pMax) < p -> write back to counter
		if pMax < p {
			err := h.buckets.Write(uint(bucketIndex), p)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *StochAvgProbabilisticCounter) Add(item []byte) error {
	hashes, _ := h.hashFunc(item)
	// take first hash as we aren't gonna use this algorithm anyway
	err := h.addHash(hashes[0])
	if err != nil {
		return err
	}
	return nil
}

func (h *StochAvgProbabilisticCounter) getBucketsPmax() []float64 {
	pMaxes := make([]float64, h.buckets.Capacity())
	for i := uint(0); i < h.buckets.Capacity(); i++ {
		counterVal, _ := h.buckets.Read(i)
		pMaxes = append(pMaxes, float64(counterVal))
	}
	return pMaxes
}

func (h *StochAvgProbabilisticCounter) Cardinality() uint {
	avgPMax := float64(0)

	for _, pMax := range h.getBucketsPmax() {
		avgPMax += pMax
	}

	avgPMax /= float64(h.buckets.Capacity())

	// E = 2^A * m
	return uint(math.Pow(2, avgPMax) * float64(h.buckets.Capacity()))
}
