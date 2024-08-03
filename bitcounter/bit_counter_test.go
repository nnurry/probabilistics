package bitcounter_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/nnurry/probabilistics/bitcounter"
)

func TestCreate(t *testing.T) {
	counter, err := bitcounter.NewSqBitCounter(64, 4)
	if err != nil {
		log.Fatal("failed to create counter:", err)
	}

	fmt.Printf("counter.Capacity() = %v\n", counter.Capacity())
	fmt.Printf("counter.BitRange() = %v\n", counter.BitRange())
	fmt.Printf("counter.TotalContainers() = %v\n", counter.TotalContainers())
	fmt.Printf("counter.ContainerSize() = %v\n", counter.ContainerSize())

	for idx := uint(0); idx < 1<<8+1; idx++ {
		fmt.Println("-------Iteration", idx+1)
		if _, _, err := counter.Increment(idx); err != nil {
			fmt.Printf("can't increment the filter: %v\n", err)
			break
		}
		fmt.Println("-----------------------------------")
	}
}
