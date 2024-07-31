package bloomfilter

import (
	"github.com/bits-and-blooms/bitset"
	"github.com/nnurry/probabilistics/hasher"
)

type BasicBloomFilterBuilder struct {
	capacity uint
	hashNum  uint
	hashName string
	hashFunc hasher.HashFunc64Type // we only use 64-bit hash function here
	b        *bitset.BitSet
}

func NewBasicBloomFilterBuilder() *BasicBloomFilterBuilder {
	defaultCapacity := uint(10000)
	defaultHashNum := uint(10)
	defaultHashName := "murmur3_128"
	defaultHashFunc := hasher.GetHashFunc64(defaultHashName)
	return &BasicBloomFilterBuilder{
		defaultCapacity,
		defaultHashNum,
		defaultHashName,
		defaultHashFunc,
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

func (f *BasicBloomFilterBuilder) SetHashFunc(hashName string) *BasicBloomFilterBuilder {
	hashFunc := hasher.GetHashFunc64(hashName)
	if hashFunc != nil {
		f.hashName = hashName
		f.hashFunc = hashFunc
	}
	return f
}

func (f *BasicBloomFilterBuilder) SetBitSet(b *bitset.BitSet) *BasicBloomFilterBuilder {
	f.b = b
	return f
}

func (f *BasicBloomFilterBuilder) Build() *BasicBloomFilter {
	return NewBasicBloomFilter(f.capacity, f.hashNum, f.hashName, f.hashFunc, f.b)
}
