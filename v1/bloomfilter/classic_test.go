package bloomfilter_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/nnurry/probabilistics/v1/bloomfilter"
)

var testParameters = []struct {
	P float64
	N uint
}{
	{0.1, 1000},
	{0.2, 1000},
	{0.3, 1000},
	{0.4, 1000},
	{0.5, 1000},
	{0.01, 1000},
	{0.02, 1000},
	{0.03, 1000},
	{0.04, 1000},
	{0.05, 1000},
	{0.001, 1000},
	{0.002, 1000},
	{0.003, 1000},
	{0.004, 1000},
	{0.005, 1000},

	{0.1, 2000},
	{0.2, 2000},
	{0.3, 2000},
	{0.4, 2000},
	{0.5, 2000},
	{0.01, 2000},
	{0.02, 2000},
	{0.03, 2000},
	{0.04, 2000},
	{0.05, 2000},
	{0.001, 2000},
	{0.002, 2000},
	{0.003, 2000},
	{0.004, 2000},
	{0.005, 2000},

	{0.1, 5000},
	{0.2, 5000},
	{0.3, 5000},
	{0.4, 5000},
	{0.5, 5000},
	{0.01, 5000},
	{0.02, 5000},
	{0.03, 5000},
	{0.04, 5000},
	{0.05, 5000},
	{0.001, 5000},
	{0.002, 5000},
	{0.003, 5000},
	{0.004, 5000},
	{0.005, 5000},

	{0.1, 10000},
	{0.2, 10000},
	{0.3, 10000},
	{0.4, 10000},
	{0.5, 10000},
	{0.01, 10000},
	{0.02, 10000},
	{0.03, 10000},
	{0.04, 10000},
	{0.05, 10000},
	{0.001, 10000},
	{0.002, 10000},
	{0.003, 10000},
	{0.004, 10000},
	{0.005, 10000},

	{0.1, 50000},
	{0.2, 50000},
	{0.3, 50000},
	{0.4, 50000},
	{0.5, 50000},
	{0.01, 50000},
	{0.02, 50000},
	{0.03, 50000},
	{0.04, 50000},
	{0.05, 50000},
	{0.001, 50000},
	{0.002, 50000},
	{0.003, 50000},
	{0.004, 50000},
	{0.005, 50000},

	{0.1, 100000},
	{0.2, 100000},
	{0.3, 100000},
	{0.4, 100000},
	{0.5, 100000},
	{0.01, 100000},
	{0.02, 100000},
	{0.03, 100000},
	{0.04, 100000},
	{0.05, 100000},
	{0.001, 100000},
	{0.002, 100000},
	{0.003, 100000},
	{0.004, 100000},
	{0.005, 100000},

	{0.1, 200000},
	{0.2, 200000},
	{0.3, 200000},
	{0.4, 200000},
	{0.5, 200000},
	{0.01, 200000},
	{0.02, 200000},
	{0.03, 200000},
	{0.04, 200000},
	{0.05, 200000},
	{0.001, 200000},
	{0.002, 200000},
	{0.003, 200000},
	{0.004, 200000},
	{0.005, 200000},

	{0.1, 500000},
	{0.2, 500000},
	{0.3, 500000},
	{0.4, 500000},
	{0.5, 500000},
	{0.01, 500000},
	{0.02, 500000},
	{0.03, 500000},
	{0.04, 500000},
	{0.05, 500000},
	{0.001, 500000},
	{0.002, 500000},
	{0.003, 500000},
	{0.004, 500000},
	{0.005, 500000},

	{0.1, 1000000},
	{0.2, 1000000},
	{0.3, 1000000},
	{0.4, 1000000},
	{0.5, 1000000},
	{0.01, 1000000},
	{0.02, 1000000},
	{0.03, 1000000},
	{0.04, 1000000},
	{0.05, 1000000},
	{0.001, 1000000},
	{0.002, 1000000},
	{0.003, 1000000},
	{0.004, 1000000},
	{0.005, 1000000},
}

func printDebug(fp float64, n, capacity, hashNum uint) {
	fmt.Printf("(%v, %d) ->\t m = %d (bits), k = %d\n", fp, n, capacity>>3>>10, hashNum)
}

func TestCreate(t *testing.T) {
	capacity, hashNum := bloomfilter.ClassicBloomEstimateParameters(0.001, 10000)
	hashName := "fakerSTKSacombank" // deliberately wrong name
	clsBloom := (bloomfilter.
		NewClassicBloomFilterBuilder().
		SetCapacity(capacity).
		SetHashNum(hashNum).
		SetHashFunc(hashName).
		SetHashScheme("dummyName").
		Build())

	if clsBloom.HashFuncName() != "murmur3_128" {
		log.Fatal("you are fake, should be murmur3_128: ", clsBloom.HashFuncName())
	}

	if clsBloom.HashSchemeName() != "enhanced_double_hashing" {
		log.Fatal("you are fake, should be enhanced_double_hashing: ", clsBloom.HashSchemeName())
	}

	testData := []string{
		"this is a test string",
		"this is second test string",
		"what is this?",
		"lalalalallala",
	}

	for _, str := range testData {
		byteStr := []byte(str)
		clsBloom.Add(byteStr)
		if !clsBloom.Contains(byteStr) {
			log.Fatal("should be in the filter")
		}
		log.Println("Found", str)
	}

	if clsBloom.Contains([]byte("this should not be in the filter")) {
		log.Fatal("this item is not added into the filter but you said yes -> false negative -> wrong")
	}
	log.Println("aight you good bro")
}

func TestEstimateParameters(t *testing.T) {
	var capacity, hashNum uint
	maxK := 0
	for _, testParameter := range testParameters {
		fp := testParameter.P
		n := testParameter.N
		capacity, hashNum = bloomfilter.ClassicBloomEstimateParameters(fp, n)
		if maxK < int(hashNum) {
			maxK = int(hashNum)
		}
		printDebug(fp, n, capacity, hashNum)
	}
	log.Println("max K =", maxK)
}
