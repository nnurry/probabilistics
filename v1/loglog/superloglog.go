package loglog

import (
	"math"
	"sort"
)

type SuperLogLog struct {
	sapc   *StochAvgProbabilisticCounter
	alphaM float64
}

func NewSuperLogLog(kBit uint, log2CounterRange uint) (*SuperLogLog, error) {
	sapc, err := NewStochAvgProbabilisticCounter(kBit, log2CounterRange)
	if err != nil {
		return nil, err
	}
	alphaM := LogLogAlphaM(sapc.buckets.Capacity())

	return &SuperLogLog{sapc, alphaM}, nil
}

func (h *SuperLogLog) Add(data []byte) error {
	err := h.sapc.Add(data)
	return err
}

func (h *SuperLogLog) Cardinality() uint {
	pMaxes := h.sapc.getBucketsPmax()
	sort.Float64s(pMaxes)

	retainRange := int(math.Round(float64(len(pMaxes)) * 0.7))

	avg70PMax := float64(0)

	for _, pMax := range pMaxes[:retainRange] {
		avg70PMax += pMax
	}

	avg70PMax /= float64(retainRange)

	cardinality := math.Pow(2, avg70PMax) * float64(retainRange)
	return uint(h.alphaM * cardinality)
}
