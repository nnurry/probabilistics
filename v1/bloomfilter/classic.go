package bloomfilter

import (
	"math"

	"github.com/bits-and-blooms/bitset"
	"github.com/nnurry/probabilistics/v1/hasher"
)

type ClassicBloomFilter struct {
	capacity       uint
	hashNum        uint
	b              *bitset.BitSet
	hashFuncName   string
	hashFunc       hasher.HashFunc64Type
	hashSchemeName string
	hashScheme     hasher.HashScheme64Type
}

const SquaredLn2 = math.Ln2 * math.Ln2

func NewClassicBloomFilter(
	capacity,
	hashNum uint,
	hashFuncName string,
	hashFunc hasher.HashFunc64Type,
	hashSchemeName string,
	hashScheme hasher.HashScheme64Type,
	b *bitset.BitSet,
) *ClassicBloomFilter {
	f := &ClassicBloomFilter{
		capacity:       capacity,
		hashNum:        hashNum,
		hashFuncName:   hashFuncName,
		hashFunc:       hashFunc,
		hashSchemeName: hashSchemeName,
		hashScheme:     hashScheme,
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
func ClassicBloomEstimateParameters(falsePositive float64, elements uint) (capacity, hashNum uint) {
	n := float64(elements)
	m := estimateCapacity(falsePositive, n)
	k := estimateHashNum(m, n)

	capacity = uint(m)
	hashNum = uint(k)
	return capacity, hashNum
}

func (f *ClassicBloomFilter) Capacity() uint                      { return f.capacity }
func (f *ClassicBloomFilter) HashNum() uint                       { return f.hashNum }
func (f *ClassicBloomFilter) BitSet() *bitset.BitSet              { return f.b }
func (f *ClassicBloomFilter) HashFuncName() string                { return f.hashFuncName }
func (f *ClassicBloomFilter) HashFunc() hasher.HashFunc64Type     { return f.hashFunc }
func (f *ClassicBloomFilter) HashSchemeName() string              { return f.hashSchemeName }
func (f *ClassicBloomFilter) HashScheme() hasher.HashScheme64Type { return f.hashScheme }

func (f *ClassicBloomFilter) bitsetIndex(hashes []uint64, numHashes int, seed uint) uint {
	hash := f.hashScheme(&hashes, numHashes, seed, f.Capacity())
	return uint(hash % uint64(f.Capacity()))
}

func (f *ClassicBloomFilter) Add(data []byte) *ClassicBloomFilter {
	hs, hn := f.hashFunc(data)
	for seed := uint(0); seed < f.HashNum(); seed++ {
		idx := f.bitsetIndex(hs, hn, seed)
		f.BitSet().Set(idx)
	}
	return f
}

func (f *ClassicBloomFilter) Contains(data []byte) bool {
	hs, hn := f.hashFunc(data)
	for seed := uint(0); seed < f.HashNum(); seed++ {
		idx := f.bitsetIndex(hs, hn, seed)
		if !f.BitSet().Test(idx) {
			return false
		}
	}
	return true
}
