package bloomfilter

import (
	"strings"

	"github.com/bits-and-blooms/bitset"
	"github.com/nnurry/probabilistics/v1/bitcounter"
	"github.com/nnurry/probabilistics/v1/hasher"
)

type CountingBloomFilter struct {
	capacity       uint
	hashNum        uint
	b              *bitset.BitSet           // 1-bit bitset
	counter        *bitcounter.SqBitCounter // n-bit bitset
	hashFuncName   string
	hashFunc       hasher.HashFunc64Type
	hashSchemeName string
	hashScheme     hasher.HashScheme64Type
}

func (f *CountingBloomFilter) Capacity() uint                      { return f.capacity }
func (f *CountingBloomFilter) HashNum() uint                       { return f.hashNum }
func (f *CountingBloomFilter) BitSet() *bitset.BitSet              { return f.b }
func (f *CountingBloomFilter) Counter() *bitcounter.SqBitCounter   { return f.counter }
func (f *CountingBloomFilter) HashFuncName() string                { return f.hashFuncName }
func (f *CountingBloomFilter) HashFunc() hasher.HashFunc64Type     { return f.hashFunc }
func (f *CountingBloomFilter) HashSchemeName() string              { return f.hashSchemeName }
func (f *CountingBloomFilter) HashScheme() hasher.HashScheme64Type { return f.hashScheme }

func (f *CountingBloomFilter) bitsetIndex(hashes []uint64, numHashes int, seed uint) uint {
	hash := f.hashScheme(&hashes, numHashes, seed, f.Capacity())
	return uint(hash % uint64(f.Capacity()))
}

func (f *CountingBloomFilter) Add(data []byte) *CountingBloomFilter {
	hs, hn := f.hashFunc(data)
	for seed := uint(0); seed < f.HashNum(); seed++ {
		idx := f.bitsetIndex(hs, hn, seed)
		beforeVal, afterVal, err := f.counter.Increment(idx)
		if strings.Contains(err.Error(), "invalid offset") {
			return f
		}
		if beforeVal == 0 && afterVal == 1 {
			f.b.Set(idx)
		}
	}
	return f
}

func (f *CountingBloomFilter) Remove(data []byte) *CountingBloomFilter {
	hs, hn := f.hashFunc(data)
	for seed := uint(0); seed < f.HashNum(); seed++ {
		idx := f.bitsetIndex(hs, hn, seed)
		beforeVal, afterVal, err := f.counter.Decrement(idx)
		if strings.Contains(err.Error(), "invalid offset") {
			return f
		}
		if beforeVal == 1 && afterVal == 0 {
			f.b.Clear(idx)
		}
	}
	return f
}

func (f *CountingBloomFilter) Contains(data []byte) bool {
	hs, hn := f.hashFunc(data)
	for seed := uint(0); seed < f.HashNum(); seed++ {
		idx := f.bitsetIndex(hs, hn, seed)
		if !f.b.Test(idx) {
			return false
		}
	}
	return true
}
