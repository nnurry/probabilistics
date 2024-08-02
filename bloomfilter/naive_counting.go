package bloomfilter

import (
	"math"

	"github.com/bits-and-blooms/bitset"
	"github.com/nnurry/probabilistics/hasher"
)

type NaiveCountingBloomFilter struct {
	capacity       uint
	hashNum        uint
	b              *bitset.BitSet
	counter        []uint8
	hashFuncName   string
	hashFunc       hasher.HashFunc64Type
	hashSchemeName string
	hashScheme     hasher.HashScheme64Type
}

func (f *NaiveCountingBloomFilter) Capacity() uint {
	return f.capacity
}

func (f *NaiveCountingBloomFilter) HashNum() uint {
	return f.hashNum
}

func (f *NaiveCountingBloomFilter) BitSet() *bitset.BitSet {
	return f.b
}

func (f *NaiveCountingBloomFilter) Counter() *[]uint8 {
	return &f.counter
}

func (f *NaiveCountingBloomFilter) HashFuncName() string {
	return f.hashFuncName
}

func (f *NaiveCountingBloomFilter) HashFunc() hasher.HashFunc64Type {
	return f.hashFunc
}

func (f *NaiveCountingBloomFilter) HashSchemeName() string {
	return f.hashSchemeName
}

func (f *NaiveCountingBloomFilter) HashScheme() hasher.HashScheme64Type {
	return f.hashScheme
}

func (f *NaiveCountingBloomFilter) bitsetIndex(hashes []uint64, numHashes int, seed uint) uint {
	hash := f.hashScheme(&hashes, numHashes, seed, f.Capacity())
	return uint(hash % uint64(f.Capacity()))
}

func (f *NaiveCountingBloomFilter) Add(data []byte) *NaiveCountingBloomFilter {
	hs, hn := f.hashFunc(data)
	for seed := uint(0); seed < f.HashNum(); seed++ {
		idx := f.bitsetIndex(hs, hn, seed)
		if f.counter[idx] < math.MaxUint8 {
			// allow to increment counter
			f.counter[idx]++
		}
		if f.counter[idx] == 1 {
			// bit position flipped to 1
			f.b.Set(idx)
		}
	}
	return f
}

func (f *NaiveCountingBloomFilter) Remove(data []byte) *NaiveCountingBloomFilter {
	hs, hn := f.hashFunc(data)
	for seed := uint(0); seed < f.HashNum(); seed++ {
		idx := f.bitsetIndex(hs, hn, seed)
		if f.counter[idx] > 0 {
			// allow to decrement counter
			f.counter[idx]--
		}
		if f.counter[idx] == 0 {
			// bit position flipped to 0
			f.b.Clear(idx)
		}
	}
	return f
}

func (f *NaiveCountingBloomFilter) Contains(data []byte) bool {
	hs, hn := f.hashFunc(data)
	for seed := uint(0); seed < f.HashNum(); seed++ {
		idx := f.bitsetIndex(hs, hn, seed)
		if !f.b.Test(idx) {
			return false
		}
	}
	return true
}
