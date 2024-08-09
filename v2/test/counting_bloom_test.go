package test

import (
	"fmt"
	"testing"

	"github.com/nnurry/probabilistics/v2/membership/bloomfilter"
	"github.com/nnurry/probabilistics/v2/utilities/register"
)

func testCountingBloomHelperBasic(fpr float64, elems uint, populationRatio float64, generateMethod string) {
	testFp := fpr
	testN := elems
	testHashGenerateMethod := generateMethod
	realTestN := uint(float64(testN) / float64(populationRatio))

	builder := bloomfilter.NewCountingBFBuilder[uint64]()
	optM, optK := bloomfilter.ClassicBFEstimateParams(testFp, testN)

	bitR, _ := register.NewRegister(optM, 1)
	countR, _ := register.NewRegister(optM, 4)

	builder = builder.
		SetCap(optM).
		SetHashNum(optK).
		SetBitRegister(bitR.(*register.BitRegister)).
		SetCountRegister(countR).
		SetHashGenerator("murmur3", 64, 128, testHashGenerateMethod)
	bf := builder.Build()

	typeName := fmt.Sprintf("%T", bf)
	fmt.Println("type of bloom filter:", typeName)
	fmt.Printf(
		"fpr = %.2f %%, n = %v, real n = %d / %f%% = %d\nm = %d, k = %d, gen method = %s\n",
		testFp*100,
		testN,
		testN, populationRatio*100, realTestN,
		optM,
		optK,
		testHashGenerateMethod,
	)

	data := [][]byte{}

	for i := uint(0); i < realTestN; i++ {
		value := []byte(fmt.Sprintf("data %b", i))
		data = append(data, value)
	}
	fmt.Printf("prepared %d test elements\n", realTestN)

	for i := range data[:testN] {
		bf.Add(data[i])
	}
	fmt.Printf("added %d test elements\n", testN)

	expectedFalseCount := realTestN - testN
	expectedFalsePerc := float64(expectedFalseCount) * 100 / float64(realTestN)
	falseCount := 0
	tp, fp, tn, fn := 0, 0, 0, 0

	for i := range data {
		ok := bf.Contains(data[i])
		if uint(i) < testN {
			// checking added data
			if ok {
				// added and found -> true positive
				tp++
			} else {
				// added but not found -> false negative
				fn++
				falseCount++
			}
		} else {
			// checking unadded data
			if ok {
				// not added but found -> false positive
				fp++
			} else {
				// not added and not found -> true negative
				tn++
				falseCount++
			}
		}
	}
	falsePerc := float64(falseCount*100.0) / float64(realTestN)

	pos, neg := register.GetBitNums(bitR)
	loadFactor := float64(testN) * 100 / float64(bf.Cap())
	bitLoadFactor := float64(pos) * 100 / float64(pos+neg)

	fmt.Printf("checked %d test elements\n", realTestN)

	fmt.Printf("load factor = %.2f %% (%d / %d) \n", loadFactor, testN, bf.Cap())
	fmt.Printf("bit load factor = %.2f %% (%d / %d) \n", bitLoadFactor, pos, pos+neg)

	fmt.Printf("false count: %v (%.2f %%)\n", falseCount, falsePerc)
	fmt.Printf("expected false count: %v (%.2f %%)\n", expectedFalseCount, expectedFalsePerc)

	fmt.Printf("true/false P = %v / %v, true/false N = %v / %v\n",
		tp, fp,
		tn, fn,
	)
}

func TestCountingBloomCreate(t *testing.T) {
	bf := bloomfilter.NewCountingBFBuilder[uint64]().Build()
	typeName := fmt.Sprintf("%T", bf)
	fmt.Println("type of bloom filter:", typeName)
}

func TestCountingBloomBasic(t *testing.T) {
	testFp := 0.1
	testN := uint(4 * 100000)
	populationRatio := 1 / 30.0
	fmt.Printf("\n\n---------test for standard---------\n\n")
	testCountingBloomHelperBasic(testFp, testN, populationRatio, "standard")
	fmt.Printf("\n\n---------test for extended double hashing---------\n\n")
	testCountingBloomHelperBasic(testFp, testN, populationRatio, "extended-double-hashing")
	fmt.Printf("\n\n---------test for kirsch-mitzenmacher---------\n\n")
	testCountingBloomHelperBasic(testFp, testN, populationRatio, "kirsch-mitzenmacher")
}
