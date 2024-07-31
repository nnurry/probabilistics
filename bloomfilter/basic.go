package bloomfilter

import (
	"math"

	"github.com/bits-and-blooms/bitset"
	"github.com/nnurry/probabilistics/hasher"
)

type BasicBloomFilter struct {
	capacity uint
	hashNum  uint
	b        *bitset.BitSet
	hashName string
	hashFunc hasher.HashFunc64Type
}

const SquaredLn2 = math.Ln2 * math.Ln2

func NewBasicBloomFilter(capacity, hashNum uint, hashName string, hashFunc hasher.HashFunc64Type, b *bitset.BitSet) *BasicBloomFilter {
	f := &BasicBloomFilter{
		capacity: capacity,
		hashNum:  hashNum,
		hashName: hashName,
		hashFunc: hashFunc,
	}

	if b != nil {
		f.b = b
	} else {
		f.b = bitset.New(capacity)
	}

	return f
}

func estimateCapacity(falsePositive float64, elements float64) float64 {
	return math.Ceil(-1 * elements * math.Log(falsePositive) / SquaredLn2)
}

func estimateHashNum(capacity float64, elements float64) float64 {
	return math.Ln2 * capacity / elements
}

// Andrii Gakhov - PDSA book (section 2.1 - page 29)
func BasicBloomEstimateParameters(falsePositive float64, elements uint) (capacity, hashNum uint) {
	n := float64(elements)
	m := estimateCapacity(falsePositive, n)
	k := estimateHashNum(m, n)

	capacity = uint(m)
	hashNum = uint(k)
	return capacity, hashNum
}

func (f *BasicBloomFilter) Capacity() uint {
	return f.capacity
}

func (f *BasicBloomFilter) HashNum() uint {
	return f.hashNum
}

func (f *BasicBloomFilter) BitSet() *bitset.BitSet {
	return f.b
}

func (f *BasicBloomFilter) HashFuncName() string {
	return f.hashName
}

func (f *BasicBloomFilter) HashFunc() hasher.HashFunc64Type {
	return f.hashFunc
}

func (f *BasicBloomFilter) bitsetIndex(hashes []uint64, numHashes int, seed uint) uint {
	hash := hasher.EnhancedDoubleHashing(&hashes, numHashes, seed, f.Capacity())
	return uint(hash % uint64(f.Capacity()))
}

func (f *BasicBloomFilter) Add(data []byte) *BasicBloomFilter {
	hs, hn := f.hashFunc(data)
	for seed := uint(0); seed < f.HashNum(); seed++ {
		idx := f.bitsetIndex(hs, hn, seed)
		f.BitSet().Set(idx)
	}
	return f
}

func (f *BasicBloomFilter) Contains(data []byte) bool {
	hs, hn := f.hashFunc(data)
	for seed := uint(0); seed < f.HashNum(); seed++ {
		idx := f.bitsetIndex(hs, hn, seed)
		if !f.BitSet().Test(idx) {
			return false
		}
	}
	return true
}
