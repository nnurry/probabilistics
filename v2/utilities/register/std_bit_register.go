package register

import (
	"fmt"
	"math"

	"github.com/nnurry/probabilistics/v2/utilities/arch"
)

// 2^k-bit register
type StdBitRegister struct {
	capacity          uint
	bitWidth          uint
	maxValue          uint
	containers        []uint
	containerCapacity uint
	totalContainers   uint
}

func newStdBitRegister(capacity, bitWidth uint) (*StdBitRegister, error) {
	if capacity <= 0 {
		return nil, fmt.Errorf(InvalidCapacityMsg, capacity)
	}
	// 64 / 2^k (log2(IntSize) > k)
	containerCapacity := uint(arch.IntSize / bitWidth)
	totalContainers := uint(math.Ceil(float64(capacity) / float64(containerCapacity)))
	containers := make([]uint, totalContainers)

	register := &StdBitRegister{
		capacity:          capacity,
		bitWidth:          bitWidth,
		maxValue:          (1 << bitWidth) - 1,
		containers:        containers,
		containerCapacity: containerCapacity,
		totalContainers:   totalContainers,
	}
	return register, nil
}

// callable when checkOffset() != nil, otherwise fatal
func (r *StdBitRegister) read(offset uint) (value uint) {
	// read whole word
	containerOffset := offset >> arch.Log2IntSize
	container := r.containers[containerOffset]
	// get bit offset in word from the left
	leftOffset := getLeftBitOffset(offset)
	// push (k-bit) at left offset to LSB
	lastCounterOffset := arch.IntSize - r.bitWidth
	container >>= (lastCounterOffset - leftOffset)
	// truncate upper [0-lastCounterOffset] bits
	// 2^k = 1[0]{k} -> 2^k - 1 = 0{1}{k-1}
	// & 2^k - 1 = keep k LSB
	value = container & ((1 << r.bitWidth) - 1)
	return value
}

// callable when checkOffset() != nil, otherwise fatal
func (r *StdBitRegister) write(offset, value uint) error {
	containerOffset := offset >> arch.Log2IntSize
	container := r.containers[containerOffset]

	leftOffset := getLeftBitOffset(offset)

	// upper: truncate n-LSB for upper part
	// 1st way: upper = container >> (arch.IntSize - 1 - (leftOffset - 1)) << (arch.IntSize - 1 - (leftOffset - 1))
	// 2nd way: upper = container & ^((1 << (arch.IntSize - (leftOffset - 1))) - 1)
	// NOTE: i will do a bit differently by sparing last k bit-width LSBs for writing new value
	// lower: truncate (n - k)-MSB for lower part

	var upper, lower uint

	lowerStartOffset := leftOffset + r.bitWidth
	if lowerStartOffset == r.bitWidth {
		// [][0-3][4-63]
		upper = value << (arch.IntSize - r.bitWidth)
		lower = container << r.bitWidth >> r.bitWidth
	} else if lowerStartOffset == arch.IntSize {
		// [0-59][60-63][]
		upper = (container >> r.bitWidth << r.bitWidth) + value
		lower = 0
	} else {
		// [0-(k-1)][k-(k+width-1)][(k+width)-63]
		upper = (container >> (arch.IntSize - leftOffset) << r.bitWidth) + value
		upper <<= (arch.IntSize - leftOffset - r.bitWidth)
		lower = container << (leftOffset + r.bitWidth) >> (leftOffset + r.bitWidth)
	}

	// write upper | lower to original container
	r.containers[containerOffset] = upper | lower
	return nil
}

func (r *StdBitRegister) Capacity() (capacity uint) {
	capacity = r.capacity
	return capacity
}

func (r *StdBitRegister) BitWidth() (bitWidth uint) {
	bitWidth = r.bitWidth
	return bitWidth
}

func (r *StdBitRegister) MaxValue() (maxValue uint) {
	maxValue = r.maxValue
	return maxValue
}

func (r *StdBitRegister) Read(offset uint) (value uint, err error) {
	offset *= r.bitWidth
	if err = checkOffset(r, offset); err != nil {
		return 0, err
	}
	value = r.read(offset)
	return value, nil
}

func (r *StdBitRegister) Write(offset uint, newValue uint) (oldValue uint, err error) {
	offset *= r.bitWidth
	if err = checkOffset(r, offset); err != nil {
		return 0, err
	}
	if checkValueOutbound(r, newValue) {
		return 0, fmt.Errorf(ExceedRegisterValueMsg, newValue, r.maxValue)
	}
	oldValue = r.read(offset)
	if oldValue == newValue {
		// same value -> no need to write
		return oldValue, nil
	}

	err = r.write(offset, newValue)
	return oldValue, err
}
func (r *StdBitRegister) Increment(offset uint) (before, after uint, err error) {
	offset *= r.bitWidth
	if err = checkOffset(r, offset); err != nil {
		return 0, 0, err
	}
	before = r.read(offset)
	after = before + 1
	if before > after {
		// overflow
		return 0, 0, fmt.Errorf("integer overflow")
	}
	if checkValueOutbound(r, after) {
		return 0, 0, fmt.Errorf(ExceedRegisterValueMsg, after, r.maxValue)
	}

	err = r.write(offset, after)
	return before, after, err
}
func (r *StdBitRegister) Decrement(offset uint) (before, after uint, err error) {
	offset *= r.bitWidth
	if err = checkOffset(r, offset); err != nil {
		return 0, 0, err
	}
	before = r.read(offset)
	after = before - 1
	if before < after {
		// underflow
		return 0, 0, fmt.Errorf("integer underflow")
	}

	err = r.write(offset, after)
	return before, after, err
}
