package register

// 2^k-bit register
type StdBitRegister struct{}

func NewStdBitRegister(capacity, bitWidth uint) (*StdBitRegister, error) {
	return nil, nil
}

func (r *StdBitRegister) Read(offset uint) (value uint, err error)
func (r *StdBitRegister) Write(offset uint, value uint) (oldValue uint, err error)
func (r *StdBitRegister) Increment(offset uint) (before, after uint, err error)
func (r *StdBitRegister) Decrement(offset uint) (before, after uint, err error)
