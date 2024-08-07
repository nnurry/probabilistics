package test

import (
	"fmt"
	"log"
	"testing"

	"github.com/nnurry/probabilistics/v2/utilities/register"
)

func TestNonStdBitRegisterCreate(t *testing.T) {
	var r register.Register
	var err error

	r, err = register.NewRegister(64, 5)
	if err != nil {
		log.Fatal("can't create 5-bit register:", err)
	}

	fmt.Printf(
		"5-bit register: capacity=%d, width=%d, max value=%d\n",
		r.Capacity(),
		r.BitWidth(),
		r.MaxValue(),
	)
}

func TestNonStdBitRegisterBasic(t *testing.T) {
	r, _ := register.NewRegister(25, 3)
	for i := 0; i < int(r.Capacity()); i++ {
		fmt.Print("r.Write(uint(i), uint(i):")
		fmt.Println(r.Write(uint(i), uint(i)%(1<<3)))
	}
	register.PrintAll(r)
}

func TestNonStdBitRegisterNormal(t *testing.T) {
	var r register.Register
	var err error

	r, err = register.NewRegister(16, 5)
	if err != nil {
		log.Fatal("can't create 5-bit register:", err)
	}

	register.PrintAll(r)

	r.Write(8, 15)
	r.Read(8)
	r.Read(9)
	r.Write(9, 16)
	r.Read(9)
	r.Write(9, 7)
	r.Read(9)
	register.PrintAll(r)
	r.Increment(7)
	r.Increment(8)
	r.Increment(9)
	r.Decrement(8)
	r.Increment(8)
	r.Increment(8)
	r.Increment(8)
	r.Increment(8)
	r.Decrement(1)
	r.Decrement(0)
	r.Increment(63)
	r.Increment(64)

	r = nil
	r, err = register.NewRegister(16, 7)
	if err != nil {
		log.Fatal("can't create 7-bit register:", err)
	}

	register.PrintAll(r)

	log.Println("---------------------------------------------------------------")

	r.Write(8, 15)
	r.Read(8)
	r.Read(9)
	r.Write(9, 16)
	r.Read(9)
	r.Write(9, 7)
	r.Read(9)
	register.PrintAll(r)
	r.Increment(7)
	r.Increment(8)
	r.Increment(9)
	r.Decrement(8)
	r.Increment(8)
	r.Increment(8)
	r.Increment(8)
	r.Increment(8)
	r.Decrement(1)
	r.Decrement(0)
	r.Increment(63)
	r.Increment(64)

	register.PrintAll(r)
}
