package test

import (
	"fmt"
	"testing"

	"github.com/nnurry/probabilistics/v2/membership/bloomfilter"
)

func TestClassicBloomCreate(t *testing.T) {
	bf := bloomfilter.NewClassicBFBuilder[uint64]().Build()
	typeName := fmt.Sprintf("%T", bf)
	fmt.Println("type of bloom filter:", typeName)
}
