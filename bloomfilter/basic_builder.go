package bloomfilter

import (
	"github.com/bits-and-blooms/bitset"
	"github.com/nnurry/probabilistics/hasher"
)

type BasicBloomFilterBuilder struct {
	capacity       uint
	hashNum        uint
	hashFuncName   string
	hashFunc       hasher.HashFunc64Type // we only use 64-bit hash function here
	hashSchemeName string
	hashScheme     hasher.HashScheme64Type // hash scheme for 64-bit hash values
	b              *bitset.BitSet
}

func NewBasicBloomFilterBuilder() *BasicBloomFilterBuilder {
	defaultCapacity := uint(10000)
	defaultHashNum := uint(10)
	defaultHashName := "murmur3_128"
	defaultHashSchemeName := "enhanced_double_hashing"
	defaultHashFunc := hasher.GetHashFunc64(defaultHashName)
	defaultHashScheme := hasher.GetHashScheme64(defaultHashSchemeName)
	return &BasicBloomFilterBuilder{
		defaultCapacity,
		defaultHashNum,
		defaultHashName,
		defaultHashFunc,
		defaultHashSchemeName,
		defaultHashScheme,
		nil,
	}
}

func (f *BasicBloomFilterBuilder) SetCapacity(capacity uint) *BasicBloomFilterBuilder {
	f.capacity = capacity
	return f
}

func (f *BasicBloomFilterBuilder) SetHashNum(hashNum uint) *BasicBloomFilterBuilder {
	f.hashNum = hashNum
	return f
}

func (f *BasicBloomFilterBuilder) SetHashFunc(hashFuncName string) *BasicBloomFilterBuilder {
	hashFunc := hasher.GetHashFunc64(hashFuncName)
	if hashFunc != nil {
		f.hashFuncName = hashFuncName
		f.hashFunc = hashFunc
	}
	return f
}

func (f *BasicBloomFilterBuilder) SetHashScheme(hashSchemeName string) *BasicBloomFilterBuilder {
	hashScheme := hasher.GetHashScheme64(hashSchemeName)
	if hashScheme != nil {
		f.hashSchemeName = hashSchemeName
		f.hashScheme = hashScheme
	}
	return f
}

func (f *BasicBloomFilterBuilder) SetBitSet(b *bitset.BitSet) *BasicBloomFilterBuilder {
	f.b = b
	return f
}

func (f *BasicBloomFilterBuilder) Build() *BasicBloomFilter {
	return NewBasicBloomFilter(
		f.capacity,
		f.hashNum,
		f.hashFuncName,
		f.hashFunc,
		f.hashSchemeName,
		f.hashScheme,
		f.b,
	)
}
