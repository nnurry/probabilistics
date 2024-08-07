package test

import (
	"fmt"
	"log"
	"testing"

	"github.com/nnurry/probabilistics/v2/utilities/register"
)

func TestBitRegisterCreate(t *testing.T) {
	var r register.Register
	var err error

	_, err = register.NewRegister(64, 0)
	if err.Error() != "invalid bit width (0 <= 0)" {
		log.Fatal("can't create 1-bit register:", err)
	}

	r, err = register.NewRegister(64, 1)
	if err != nil {
		log.Fatal("can't create 1-bit register:", err)
	}

	fmt.Printf(
		"1-bit register: capacity=%d, width=%d, max value=%d\n",
		r.Capacity(),
		r.BitWidth(),
		r.MaxValue(),
	)
}

func TestBitRegisterBasic(t *testing.T) {
	var r register.Register
	var err error

	r, err = register.NewRegister(10000, 1)
	if err != nil {
		log.Fatal("can't create 1-bit register:", err)
	}

	fmt.Println(r.Increment(72))
	fmt.Println(r.Read(72))
	fmt.Println(r.Read(73))
	fmt.Println(r.Read(10000))
	fmt.Println(r.Write(10000, 1))
	fmt.Println(r.Write(9999, 1))
	fmt.Println(r.Write(9999, 2))
}
