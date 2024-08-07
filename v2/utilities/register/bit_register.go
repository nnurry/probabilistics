package register

import (
	"fmt"
	"math"

	"github.com/nnurry/probabilistics/v2/utilities/arch"
)

// 1-bit register
type BitRegister struct {
	capacity          uint
	containers        []uint
	containerCapacity uint
	totalContainers   uint
}

func newBitRegister(capacity uint) (*BitRegister, error) {
	if capacity <= 0 {
		return nil, fmt.Errorf(InvalidCapacityMsg, capacity)
	}
	containerCapacity := uint(arch.IntSize)
	totalContainers := uint(math.Ceil(float64(capacity) / float64(containerCapacity)))
	containers := make([]uint, totalContainers)

	register := &BitRegister{
		capacity:          capacity,
		containers:        containers,
		containerCapacity: containerCapacity,
		totalContainers:   totalContainers,
	}
	return register, nil
}

// callable when checkOffset() != nil, otherwise fatal
func (r *BitRegister) read(offset uint) (value, helperValue uint) {
	// read whole word
	containerOffset := offset >> arch.Log2IntSize
	container := r.containers[containerOffset]
	// get bit offset in word from the left
	leftOffset := getLeftBitOffset(offset)
	// get bit offset in word from the right for bit truncation
	rightOffset := (r.containerCapacity - 1) - leftOffset
	// truncate all bits except for bit at index
	helperValue = uint(1 << rightOffset) // [[0]*a]x[[0]*b]
	value = container & helperValue
	if value != 0 {
		value = 1
	}
	return value, helperValue
}

func (r *BitRegister) Capacity() (capacity uint) {
	capacity = r.capacity
	return capacity
}

func (r *BitRegister) BitWidth() (bitWidth uint) {
	bitWidth = 1
	return bitWidth
}

func (r *BitRegister) MaxValue() (maxValue uint) {
	maxValue = 1
	return maxValue
}

func (r *BitRegister) Read(offset uint) (value uint, err error) {
	if err = checkOffset(r, offset); err != nil {
		return 0, err
	}
	value, _ = r.read(offset)
	return value, nil
}

func (r *BitRegister) Write(offset uint, value uint) (oldValue uint, err error) {
	if err = checkOffset(r, offset); err != nil {
		return 0, err
	}

	if checkValueOutbound(r, value) {
		return 0, fmt.Errorf(ExceedRegisterValueMsg, value, 1)
	}

	oldValue, helperValue := r.read(offset)

	if oldValue == value {
		// same value -> no need to write
		return oldValue, nil
	}

	// read whole word
	containerOffset := offset >> arch.Log2IntSize
	container := r.containers[containerOffset]
	if value == 1 {
		// set bit
		// xx0xx | 00100 = xx1xx
		container |= helperValue
	} else {
		// clear bit
		// xx1xx & 11011 = xx0xx
		container &= helperValue
	}

	r.containers[containerOffset] = container
	return oldValue, nil
}

func (r *BitRegister) Increment(offset uint) (before, after uint, err error) {
	before, err = r.Write(offset, 1)
	return before, 1, err
}

func (r *BitRegister) Decrement(offset uint) (before, after uint, err error) {
	before, err = r.Write(offset, 0)
	return before, 1, err
}
