package bloomfilter

import (
	"github.com/bits-and-blooms/bitset"
	"github.com/nnurry/probabilistics/hasher"
)

type ClassicBloomFilterBuilder struct {
	capacity       uint
	hashNum        uint
	hashFuncName   string
	hashFunc       hasher.HashFunc64Type // we only use 64-bit hash function here
	hashSchemeName string
	hashScheme     hasher.HashScheme64Type // hash scheme for 64-bit hash values
	b              *bitset.BitSet
}

func NewClassicBloomFilterBuilder() *ClassicBloomFilterBuilder {
	defaultCapacity := uint(10000)
	defaultHashNum := uint(10)
	defaultHashName := "murmur3_128"
	defaultHashSchemeName := "enhanced_double_hashing"
	defaultHashFunc := hasher.GetHashFunc64(defaultHashName)
	defaultHashScheme := hasher.GetHashScheme64(defaultHashSchemeName)
	return &ClassicBloomFilterBuilder{
		defaultCapacity,
		defaultHashNum,
		defaultHashName,
		defaultHashFunc,
		defaultHashSchemeName,
		defaultHashScheme,
		nil,
	}
}

func (f *ClassicBloomFilterBuilder) SetCapacity(capacity uint) *ClassicBloomFilterBuilder {
	f.capacity = capacity
	return f
}

func (f *ClassicBloomFilterBuilder) SetHashNum(hashNum uint) *ClassicBloomFilterBuilder {
	f.hashNum = hashNum
	return f
}

func (f *ClassicBloomFilterBuilder) SetHashFunc(hashFuncName string) *ClassicBloomFilterBuilder {
	hashFunc := hasher.GetHashFunc64(hashFuncName)
	if hashFunc != nil {
		f.hashFuncName = hashFuncName
		f.hashFunc = hashFunc
	}
	return f
}

func (f *ClassicBloomFilterBuilder) SetHashScheme(hashSchemeName string) *ClassicBloomFilterBuilder {
	hashScheme := hasher.GetHashScheme64(hashSchemeName)
	if hashScheme != nil {
		f.hashSchemeName = hashSchemeName
		f.hashScheme = hashScheme
	}
	return f
}

func (f *ClassicBloomFilterBuilder) SetBitSet(b *bitset.BitSet) *ClassicBloomFilterBuilder {
	f.b = b
	return f
}

func (f *ClassicBloomFilterBuilder) Build() *ClassicBloomFilter {
	return NewClassicBloomFilter(
		f.capacity,
		f.hashNum,
		f.hashFuncName,
		f.hashFunc,
		f.hashSchemeName,
		f.hashScheme,
		f.b,
	)
}
