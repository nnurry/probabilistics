package test

import (
	"fmt"
	"log"
	"testing"

	"github.com/nnurry/probabilistics/v2/utilities/hasher"
)

func TestHashGeneratorCreate(t *testing.T) {
	g, _ := hasher.NewHashGenerator[uint64]("murmur3", 64, 128, "extended-double-hashing")
	data, err := g.GenerateHash([]byte("sample"), uint64(13), 17, 4)
	if err != nil {
		log.Fatal("damn", err)
	}
	fmt.Println("data:", data)
}
