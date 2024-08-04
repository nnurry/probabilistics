package loglog

import (
	"math"
)

type HyperLogLog struct {
	sapc   *StochAvgProbabilisticCounter
	alphaM float64
}

func HyperLogLogAlphaM(capacity uint) float64 {
	var alphaM float64
	if capacity <= 16 {
		alphaM = 0.673
	} else if capacity <= 32 {
		alphaM = 0.697
	} else if capacity <= 64 {
		alphaM = 0.709
	} else {
		alphaM = 0.7213 * float64(capacity) / (1.079 + float64(capacity))
	}
	return alphaM
}

func NewHyperLogLog(kBit uint, log2CounterRange uint) (*HyperLogLog, error) {
	sapc, err := NewStochAvgProbabilisticCounter(kBit, log2CounterRange)
	if err != nil {
		return nil, err
	}

	alphaM := HyperLogLogAlphaM(sapc.buckets.Capacity())

	return &HyperLogLog{sapc, alphaM}, nil
}

func (h *HyperLogLog) Add(data []byte) error {
	err := h.sapc.Add(data)
	return err
}

func (h *HyperLogLog) Cardinality() uint {
	// formula = alphaM * m / (m / normalized harmonic mean of p max)
	nNormHarmonicAvgPMax := float64(0)

	// sigma(2^(pj))j[p0...pn-1] -> normalized harmonic mean
	for _, pMax := range h.sapc.getBucketsPmax() {
		nNormHarmonicAvgPMax += math.Pow(2, -pMax)
	}

	squaredM := float64(h.sapc.buckets.Capacity() * h.sapc.buckets.Capacity())
	nNormHarmonicAvgPMax = squaredM / nNormHarmonicAvgPMax

	return uint(h.alphaM * nNormHarmonicAvgPMax)
}
