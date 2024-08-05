package register

// x-bit register (x != 1 && x != 2^k)
type NonStdBitRegister struct{}

func NewNonStdBitRegister(capacity, bitWidth uint) (*NonStdBitRegister, error) {
	return nil, nil
}

func (r *NonStdBitRegister) Read(offset uint) (value uint, err error)
func (r *NonStdBitRegister) Write(offset uint, value uint) (oldValue uint, err error)
func (r *NonStdBitRegister) Increment(offset uint) (before, after uint, err error)
func (r *NonStdBitRegister) Decrement(offset uint) (before, after uint, err error)
