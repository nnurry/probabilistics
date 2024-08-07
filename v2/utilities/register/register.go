package register

import (
	"fmt"
	"math"

	"github.com/nnurry/probabilistics/v2/utilities/arch"
)

type Register interface {
	Capacity() (capacity uint)
	BitWidth() (bitWidth uint)
	MaxValue() (maxValue uint)
	Read(offset uint) (value uint, err error)
	Write(offset uint, value uint) (oldValue uint, err error)
	Increment(offset uint) (before, after uint, err error)
	Decrement(offset uint) (before, after uint, err error)
}

const (
	InvalidCapacityMsg = "invalid capacity (%v <= 0)"
)

const (
	invalidRegisterValueMsg = "invalid register value"
	ExceedRegisterValueMsg  = invalidRegisterValueMsg + " (%v > %v)"
)

const (
	invalidOffsetMsg      = "invalid offset"
	UndivisibleOffsetMsg  = invalidOffsetMsg + " (%v mod %v != 0)"
	UpperInvalidOffsetMsg = invalidOffsetMsg + " (%v > %v)"
	LowerInvalidOffsetMsg = invalidOffsetMsg + " (%v < 0)"
)

const (
	invalidBitWidth     = "invalid bit width"
	NonPositiveBitWidth = invalidBitWidth + " (%v <= 0)"
	ExceedBitWidth      = invalidBitWidth + " (%v > %v)"
)

func lastCounterOffset(r Register) uint {
	return (r.Capacity() - 1) * r.BitWidth()
}

func getLeftBitOffset(offset uint) uint {
	return offset & (arch.IntSize - 1)
}

func checkOffset(r Register, offset uint) error {
	// access inappropriate register range
	if offset%r.BitWidth() != 0 {
		return fmt.Errorf(UndivisibleOffsetMsg, offset, r.BitWidth())
	}
	// invalid upper bound
	lastCounterOffset := lastCounterOffset(r)
	if offset > lastCounterOffset {
		return fmt.Errorf(UpperInvalidOffsetMsg, offset, lastCounterOffset)
	}
	return nil
}

func PrintAll(r Register) {
	values := []uint{}
	for i := uint(0); i < r.Capacity(); i++ {
		value, _ := r.Read(i)
		values = append(values, value)
	}
	fmt.Printf(
		"%d-bit register:\tcapacity=%d (allocated=%d)\tvalues:%v\n",
		r.BitWidth(),
		r.Capacity(),
		len(values),
		values,
	)
}

func checkValueOutbound(r Register, value uint) bool {
	return value > r.MaxValue()
}

func NewRegister(capacity, bitWidth uint) (r Register, err error) {
	if bitWidth == 0 {
		return nil, fmt.Errorf(NonPositiveBitWidth, bitWidth)
	}
	if bitWidth > arch.IntSize {
		// must be smaller than word size (32 or 64 bit width)
		return nil, fmt.Errorf(ExceedBitWidth, bitWidth, arch.IntSize)
	}

	if bitWidth == 1 {
		// 1-bit register
		r, err = newBitRegister(capacity)
	} else if math.Floor(math.Log2(float64(bitWidth))) == math.Ceil(math.Log2(float64(bitWidth))) {
		// 2^k-bit register
		r, err = newStdBitRegister(capacity, bitWidth)
	} else {
		// weird registers (5-bit, 6-bit) - theoretical HLL uses this
		r, err = newNonStdBitRegister(capacity, bitWidth)
	}

	return r, err
}
