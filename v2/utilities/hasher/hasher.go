package hasher

import "fmt"

// errors when init hash functions
const (
	NoMatchingHashFamilyMsg  = "no matching hash family for %s"
	InvalidHashFuncConfigMsg = "invalid hash configs: (family = %v, platform bit = %v, output bit = %v)"
)

// errors in runtime
const (
	InvalidSeedTypeMsg = "invalid seed type (!= %s)"
)

// possible output type of hash function is []number, prevalently []uint64
type HashOutType interface {
	uint | uint32 | uint64
}

type HashFunction[T HashOutType] func([]byte, interface{}) ([]T, error)
type HashAttribute struct {
	hashFamily  string
	platformBit uint
	outputBit   uint
}

var unsignedInt32HashFunctions = map[HashAttribute]HashFunction[uint32]{}
var unsignedInt64HashFunctions = map[HashAttribute]HashFunction[uint64]{
	{"murmur3", 64, 128}: murmur3Hash128,
}

func NewHashFunction[T HashOutType](family string, platformBit uint, outputBit uint) (HashFunction[T], error) {
	var genericRef T
	hashAttr := HashAttribute{family, platformBit, outputBit}
	typeName := fmt.Sprintf("%T", genericRef)

	switch typeName {
	case "uint64":
		if hf, ok := unsignedInt64HashFunctions[hashAttr]; ok {
			return any(hf).(HashFunction[T]), nil
		}
	case "uint32":
		if hf, ok := unsignedInt32HashFunctions[hashAttr]; ok {
			return any(hf).(HashFunction[T]), nil
		}
	}
	return nil, fmt.Errorf(InvalidHashFuncConfigMsg, family, platformBit, outputBit)
}
