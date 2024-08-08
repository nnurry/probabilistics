package bloomfilter

import (
	"math"

	"github.com/nnurry/probabilistics/v2/utilities/hasher"
	"github.com/nnurry/probabilistics/v2/utilities/register"
)

const SquaredLn2 = math.Ln2 * math.Ln2

type ClassicBF[T hasher.HashOutType] struct {
	cap uint
	k   uint
	r   *register.BitRegister
	h   hasher.HashGenerator[T]
}

func estCap(fpr float64, elems float64) float64 {
	return math.Ceil(-1 * elems * math.Log(fpr) / SquaredLn2)
}
func estK(capacity float64, elems float64) float64 { return math.Ln2 * capacity / elems }
func ClassicBFEstimateParams(fpr float64, elems uint) (m, k uint) {
	n := float64(elems)
	mF64 := estCap(fpr, n)
	kF64 := estK(mF64, n)

	m = uint(mF64)
	k = uint(kF64)

	return m, k
}
