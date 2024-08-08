package bloomfilter

import (
	"github.com/nnurry/probabilistics/v2/utilities/hasher"
	"github.com/nnurry/probabilistics/v2/utilities/register"
)

type ClassicBFBuilder[T hasher.HashOutType] struct {
	cap uint
	k   uint
	h   hasher.HashGenerator[T]
}

func NewClassicBFBuilder[T hasher.HashOutType]() *ClassicBFBuilder[T] {
	defaultCap, defaultK := ClassicBFEstimateParams(0.01, 10000)
	defaultHasher, _ := hasher.NewHashGenerator[T]("murmur3", 64, 128, "extended-double-hashing")
	return &ClassicBFBuilder[T]{defaultCap, defaultK, *defaultHasher}
}

func (b *ClassicBFBuilder[T]) SetCap(cap uint) *ClassicBFBuilder[T] {
	b.cap = cap
	return b
}

func (b *ClassicBFBuilder[T]) SetHashNum(k uint) *ClassicBFBuilder[T] {
	b.k = k
	return b
}

func (b *ClassicBFBuilder[T]) SetHashGenerator(hashFamily string, platformBit uint, outputBit uint, generateMethod string) *ClassicBFBuilder[T] {
	hashGenerator, err := hasher.NewHashGenerator[T](hashFamily, platformBit, outputBit, generateMethod)
	if err != nil {
		return b
	}
	b.h = *hashGenerator
	return b
}

func (b *ClassicBFBuilder[T]) Build() *ClassicBF[T] {
	r, _ := register.NewRegister(b.cap, 1)
	bf := &ClassicBF[T]{
		cap: b.cap,
		k:   b.k,
		r:   r.(*register.BitRegister),
		h:   b.h,
	}
	return bf
}
