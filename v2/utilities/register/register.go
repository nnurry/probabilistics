package register

import (
	"fmt"
	"math"

	"github.com/nnurry/probabilistics/v2/utilities/arch"
)

type Register interface {
	Read(uint) (uint, error)
	Write(uint, uint) (uint, error)
	Increment(uint) (uint, uint, error)
	Decrement(uint) (uint, uint, error)
}

const (
	InvalidCapacityMsg = "invalid capacity (%v <= 0)"
)

const (
	invalidOffsetMsg      = "invalid offset"
	UpperInvalidOffsetMsg = invalidOffsetMsg + " (%v > %v)"
	LowerInvalidOffsetMsg = invalidOffsetMsg + " (%v < 0)"
)

const (
	invalidBitWidth = "invalid bit width"
	ExceedBitWidth  = invalidBitWidth + " (%v > %v)"
)

func NewRegister(capacity, bitWidth uint) (r Register, err error) {
	if bitWidth > arch.IntSize {
		// must be smaller than word size (32 or 64 bit width)
		return nil, fmt.Errorf(ExceedBitWidth, bitWidth, arch.IntSize)
	}
	if bitWidth == 1 {
		// 1-bit register
		register, err := NewBitRegister(capacity)
		return register, err
	}
	if math.Floor(math.Log2(float64(bitWidth))) == math.Ceil(math.Log2(float64(bitWidth))) {
		// 2^k-bit register
		register, err := NewStdBitRegister(capacity, bitWidth)
		return register, err
	}

	register, err := NewNonStdBitRegister(capacity, bitWidth)
	return register, err
}
