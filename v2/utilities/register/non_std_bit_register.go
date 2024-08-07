package register

import (
	"fmt"
	"math"

	"github.com/nnurry/probabilistics/v2/utilities/arch"
)

// x-bit register (x != 1 && x != 2^k)
type NonStdBitRegister struct {
	capacity        uint
	maxValue        uint
	bitWidth        uint
	containers      []uint
	totalContainers uint
}

func newNonStdBitRegister(capacity, bitWidth uint) (*NonStdBitRegister, error) {
	if capacity <= 0 {
		return nil, fmt.Errorf(InvalidCapacityMsg, capacity)
	}

	totalContainers := uint(math.Ceil(float64(capacity*bitWidth) / arch.IntSize))
	containers := make([]uint, totalContainers)

	register := &NonStdBitRegister{
		capacity:        capacity,
		maxValue:        (1 << bitWidth) - 1,
		bitWidth:        bitWidth,
		containers:      containers,
		totalContainers: totalContainers,
	}
	return register, nil
}

func (r *NonStdBitRegister) rightBitWidth(leftOffset uint) uint {
	maxOffsetExclusive := leftOffset + r.bitWidth
	if maxOffsetExclusive > arch.IntSize {
		return getLeftBitOffset(maxOffsetExclusive)
	}
	return 0
}

// callable when checkOffset() != nil, otherwise fatal
func (r *NonStdBitRegister) read(offset uint) (value uint, rightSize uint) {
	leftOffset := getLeftBitOffset(offset)
	containerOffset := offset >> arch.Log2IntSize
	container := r.containers[containerOffset]

	rightSize = r.rightBitWidth(leftOffset)

	if rightSize > 0 {
		// case 1: counter is stored in 2 containers
		leftValue := (container << rightSize) & ((1 << r.bitWidth) - 1)
		rightValue := r.containers[containerOffset+1] >> (arch.IntSize - rightSize)
		value = leftValue | rightValue
	} else {
		// case 2: counter is in 1 container
		// we can re-use StdBitRegister method but let's re-write in case of later modifications
		lastCounterOffset := arch.IntSize - r.bitWidth
		container >>= (lastCounterOffset - leftOffset)
		value = container & ((1 << r.bitWidth) - 1)
	}
	return value, rightSize
}

// callable when checkOffset() != nil, otherwise fatal
func (r *NonStdBitRegister) write(offset, value uint, rightSize uint) error {
	leftOffset := getLeftBitOffset(offset)
	containerOffset := offset >> arch.Log2IntSize
	container := r.containers[containerOffset]

	if rightSize > 0 {
		// case 1: counter is stored in 2 containers
		leftValue := value >> rightSize                   // write to LSBs of 1st container
		rightValue := value << (arch.IntSize - rightSize) // write to MSBs of 2nd container

		container = container >> (r.bitWidth - rightSize) << (r.bitWidth - rightSize)
		r.containers[containerOffset] = container | leftValue

		containerOffset++
		container = r.containers[containerOffset]

		container = container << rightSize >> rightSize
		r.containers[containerOffset] = container | rightValue
	} else {
		// case 2: counter is stored in 1 container
		var upper, lower uint
		lowerStartOffset := leftOffset + r.bitWidth

		if lowerStartOffset == r.bitWidth {
			upper = value << (arch.IntSize - r.bitWidth)
			lower = container << r.bitWidth >> r.bitWidth
		} else if lowerStartOffset == arch.IntSize {
			upper = (container >> r.bitWidth << r.bitWidth) + value
			lower = 0
		} else {
			upper = (container >> (arch.IntSize - leftOffset) << r.bitWidth) + value
			upper <<= (arch.IntSize - leftOffset - r.bitWidth)
			lower = container << (leftOffset + r.bitWidth) >> (leftOffset + r.bitWidth)
		}

		r.containers[containerOffset] = upper | lower
	}
	return nil
}

func (r *NonStdBitRegister) Capacity() (capacity uint) {
	capacity = r.capacity
	return capacity
}
func (r *NonStdBitRegister) BitWidth() (bitWidth uint) {
	bitWidth = r.bitWidth
	return bitWidth
}

func (r *NonStdBitRegister) MaxValue() (maxValue uint) {
	maxValue = r.maxValue
	return maxValue
}

func (r *NonStdBitRegister) Read(offset uint) (value uint, err error) {
	offset *= r.bitWidth
	if err = checkOffset(r, offset); err != nil {
		return 0, err
	}
	value, _ = r.read(offset)
	return value, nil
}

func (r *NonStdBitRegister) Write(offset uint, newValue uint) (oldValue uint, err error) {
	offset *= r.bitWidth
	if err = checkOffset(r, offset); err != nil {
		return 0, err
	}
	if checkValueOutbound(r, newValue) {
		return 0, fmt.Errorf(ExceedRegisterValueMsg, newValue, r.maxValue)
	}
	var rightSize uint
	oldValue, rightSize = r.read(offset)

	if oldValue == newValue {
		// same value -> no need to write
		return oldValue, nil
	}

	err = r.write(offset, newValue, rightSize)
	return oldValue, err
}
func (r *NonStdBitRegister) Increment(offset uint) (before, after uint, err error) {
	offset *= r.bitWidth
	if err = checkOffset(r, offset); err != nil {
		return 0, 0, err
	}
	var rightSize uint
	before, rightSize = r.read(offset)
	after = before + 1
	if before > after {
		// overflow
		return 0, 0, fmt.Errorf("integer overflow")
	}
	if checkValueOutbound(r, after) {
		return 0, 0, fmt.Errorf(ExceedRegisterValueMsg, after, r.maxValue)
	}

	err = r.write(offset, after, rightSize)
	return before, after, err
}

func (r *NonStdBitRegister) Decrement(offset uint) (before, after uint, err error) {
	offset *= r.bitWidth
	if err = checkOffset(r, offset); err != nil {
		return 0, 0, err
	}
	var rightSize uint
	before, rightSize = r.read(offset)
	after = before - 1
	if before < after {
		// underflow
		return 0, 0, fmt.Errorf("integer underflow")
	}

	err = r.write(offset, after, rightSize)
	return before, after, err
}
