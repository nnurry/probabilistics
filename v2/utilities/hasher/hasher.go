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

type HashFunction struct {
	f             func([]byte, interface{}) (interface{}, error)
	outAssertType string
}
type HashAttribute struct {
	platformBit uint
	outputBit   uint
}
type HashConfigurations map[HashAttribute]HashFunction

var supportedHashFunctions = map[string]HashConfigurations{
	"murmur3": murmur3HashConfigs,
}

func NewHashFunction(family string, platformBit uint, outputBit uint) (HashFunction, error) {
	var hashFunc HashFunction
	hashConf, ok := supportedHashFunctions[family]
	// no family match -> throw error
	if !ok {
		return hashFunc, fmt.Errorf(NoMatchingHashFamilyMsg, family)
	}
	hashAttr := HashAttribute{platformBit, outputBit}
	hashFunc, ok = hashConf[hashAttr]
	// hash function of input config not found -> throw error
	if !ok {
		return hashFunc, fmt.Errorf(InvalidHashFuncConfigMsg, family, platformBit, outputBit)
	}
	return hashFunc, nil
}
