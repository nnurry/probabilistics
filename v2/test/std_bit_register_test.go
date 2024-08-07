package test

import (
	"fmt"
	"log"
	"testing"

	"github.com/nnurry/probabilistics/v2/utilities/register"
)

func TestStdBitRegisterCreate(t *testing.T) {
	var r register.Register
	var err error

	r, err = register.NewRegister(64, 4)
	if err != nil {
		log.Fatal("can't create 4-bit register:", err)
	}

	fmt.Printf(
		"4-bit register: capacity=%d, width=%d, max value=%d\n",
		r.Capacity(),
		r.BitWidth(),
		r.MaxValue(),
	)
}

func TestStdBitRegisterBasic(t *testing.T) {
	r, _ := register.NewRegister(25, 4)
	for i := 0; i < int(r.Capacity()); i++ {
		fmt.Print("r.Write(uint(i), uint(i)):")
		fmt.Println(r.Write(uint(i), uint(i)%(1<<3)))
	}
	register.PrintAll(r)
}

func TestStdBitRegisterNormal(t *testing.T) {
	var r register.Register
	var err error

	r, err = register.NewRegister(16, 4)
	if err != nil {
		log.Fatal("can't create 4-bit register:", err)
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
	r, err = register.NewRegister(16, 8)
	if err != nil {
		log.Fatal("can't create 8-bit register:", err)
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
