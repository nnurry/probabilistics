package bitcounter

import (
	"fmt"
	"math"
)

type BitCounter struct {
	capacity        uint
	bitRange        uint
	containers      []uint64
	containerSize   uint
	totalContainers uint
}

var bitShiftFactor = uint(math.Log2(64))

func (b *BitCounter) Capacity() uint        { return b.capacity }
func (b *BitCounter) BitRange() uint        { return b.bitRange }
func (b *BitCounter) ContainerSize() uint   { return b.containerSize }
func (b *BitCounter) TotalContainers() uint { return b.totalContainers }
func (b *BitCounter) Containers() *[]uint64 { return &b.containers }

func NewBitCounter(capacity, bitRange uint) (*BitCounter, error) {
	errorString := "can't create counting set: (%s)"

	if bitRange != 2 && bitRange != 4 && bitRange != 8 && bitRange != 16 && bitRange != 32 {
		return nil, fmt.Errorf(errorString, fmt.Sprintf("bitRange = %d != 2^n", bitRange))
	}

	if capacity <= 0 {
		return nil, fmt.Errorf(errorString, fmt.Sprintf("capacity = %d <= 0", capacity))
	}

	totalBits := capacity * bitRange
	totalContainers := uint(math.Ceil(float64(totalBits) / 64))

	bitCounter := BitCounter{
		capacity:        capacity,
		bitRange:        bitRange,
		containerSize:   64 / bitRange,
		totalContainers: totalContainers,
		containers:      make([]uint64, totalContainers),
	}

	return &bitCounter, nil
}

func (b *BitCounter) checkOffset(offsetIndex uint) error {
	if offsetIndex >= b.capacity*(b.bitRange-1) {
		// invalid offset (must be a factor of b.bitRange)
		return fmt.Errorf("invalid offset (exceed limit)")
	}
	if offsetIndex%b.bitRange != 0 {
		// invalid offset (must be a factor of b.bitRange)
		return fmt.Errorf("invalid offset (%d not divisible by %d)", offsetIndex, b.bitRange)
	}
	return nil
}

func (b *BitCounter) Read(offsetIndex uint) (counterValue uint64, err error) {
	err = b.checkOffset(offsetIndex)
	counterValue = b.containers[offsetIndex>>bitShiftFactor]

	startIndex := offsetIndex % 64

	// shift to MSB
	distanceToLeft := startIndex
	counterValue <<= uint64(distanceToLeft)

	// shift to LSB
	distanceToRight := 64 - b.bitRange
	counterValue >>= uint64(distanceToRight)

	return counterValue, err
}

func (b *BitCounter) Write(offsetIndex uint, value uint64) (err error) {
	err = b.checkOffset(offsetIndex)

	startIndex := offsetIndex % 64
	endIndex := startIndex + b.bitRange - 1

	// shift value back to offset index
	value <<= uint64(63 - endIndex)

	container := b.containers[offsetIndex>>bitShiftFactor]

	upperEndIndex := startIndex - 1
	lowerStartIndex := endIndex + 1

	lowerPositive := container << uint64(lowerStartIndex) >> uint64(lowerStartIndex)
	upperPositive := container >> (63 - upperEndIndex) << (63 - upperEndIndex)

	b.containers[offsetIndex>>bitShiftFactor] = upperPositive | value | lowerPositive
	return err
}

func (b *BitCounter) update(offsetIndex uint, delta uint, isIncrement bool) (beforeValue, afterValue uint64, err error) {
	beforeValue, err = b.Read(offsetIndex)
	if err != nil {
		return uint64(0), uint64(0), err
	}
	delta64 := uint64(delta)

	if isIncrement {
		maxValue := uint64((1 << b.bitRange) - 1)
		afterValue = beforeValue + delta64
		if afterValue > maxValue {
			return beforeValue, uint64(0), fmt.Errorf(
				"incremented value exceed limit (%04b + %04b = %08b >= %04b)",
				beforeValue, delta64,
				afterValue, maxValue,
			)
		}
	} else {
		if delta64 > beforeValue {
			// cause overflow -> handle differently
			return beforeValue, uint64(0), fmt.Errorf("negative decremented value (%04b - %04b < 0)", beforeValue, delta64)
		}
		afterValue = beforeValue - delta64
	}
	beforeValue = afterValue

	b.Write(offsetIndex, afterValue)

	return beforeValue, afterValue, err
}

func (b *BitCounter) Increment(offsetIndex uint) (uint64, uint64, error) {
	return b.update(offsetIndex*b.bitRange, 1, true)
}

func (b *BitCounter) Decrement(offsetIndex uint) (uint64, uint64, error) {
	return b.update(offsetIndex*b.bitRange, 1, false)
}
