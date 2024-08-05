package loglog

import "math"

type LogLog struct {
	sapc   *StochAvgProbabilisticCounter
	alphaM float64
}

const PiSquare float64 = math.Pi * math.Pi
const Ln2Square float64 = math.Ln2 * math.Ln2

func LogLogAlphaM(capacity uint) float64 {
	// alphaM = alphaInf = 0.39701 (lim to infinity calculated)
	alphaM := 0.39701
	if capacity >= 64 {
		alphaM -= (2*PiSquare + Ln2Square) / float64(48*capacity)
	}
	return alphaM
}

func NewLogLog(kBit uint, log2CounterRange uint) (*LogLog, error) {
	sapc, err := NewStochAvgProbabilisticCounter(kBit, log2CounterRange)
	if err != nil {
		return nil, err
	}
	alphaM := LogLogAlphaM(sapc.buckets.Capacity())

	return &LogLog{sapc, alphaM}, nil
}

func (h *LogLog) Add(data []byte) error {
	err := h.sapc.Add(data)
	return err
}

func (h *LogLog) Cardinality() uint {
	return uint(h.alphaM * float64(h.sapc.Cardinality()))
}
