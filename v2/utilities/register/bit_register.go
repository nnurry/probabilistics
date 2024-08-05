package register

import (
	"fmt"
	"math"

	"github.com/nnurry/probabilistics/v2/utilities/arch"
)

// 1-bit register
type BitRegister struct {
	capacity        uint
	containers      []uint
	containerSize   uint
	totalContainers uint
}

func NewBitRegister(capacity uint) (*BitRegister, error) {
	if capacity <= 0 {
		return nil, fmt.Errorf(InvalidCapacityMsg, capacity)
	}
	containers := make([]uint, capacity)
	containerSize := uint(arch.IntSize)
	totalContainers := math.Ceil(float64(capacity) / float64(containerSize))

	register := &BitRegister{
		capacity:        capacity,
		containers:      containers,
		containerSize:   uint(containerSize),
		totalContainers: uint(totalContainers),
	}
	return register, nil
}

func (r *BitRegister) getLeftBitOffset(offset uint) uint {
	// containerSize is either 32/64 (32-bit & 64-bit word with 32/64 1-bit registers)
	// x % y, y = 2^k -> x % y = x & (y - 1)
	return offset & (r.containerSize - 1)
}

func (r *BitRegister) checkOffset(offset uint) error {
	// invalid upper bound
	if offset > r.capacity-1 {
		return fmt.Errorf(UpperInvalidOffsetMsg, offset, r.capacity-1)
	}
	return nil
}

// callable when checkOffset() != nil, otherwise fatal
func (r *BitRegister) read(offset uint) (value, helperValue uint) {
	// read whole word
	containerOffset := offset >> arch.Log2IntSize
	container := r.containers[containerOffset]
	// get bit offset in word from the left
	leftOffset := r.getLeftBitOffset(offset)
	// get bit offset in word from the right for bit truncation
	rightOffset := (r.containerSize - 1) - leftOffset
	// truncate all bits except for bit at index
	helperValue = uint(1 << rightOffset) // [[0]*a]x[[0]*b]
	value = container & helperValue
	if value != 0 {
		value = 1
	}
	return value, helperValue
}

func (r *BitRegister) Read(offset uint) (value uint, err error) {
	if err = r.checkOffset(offset); err != nil {
		return uint(0), err
	}
	value, _ = r.read(offset)
	return value, nil
}

func (r *BitRegister) Write(offset uint, value uint) (oldValue uint, err error) {
	if err = r.checkOffset(offset); err != nil {
		return 0, err
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
