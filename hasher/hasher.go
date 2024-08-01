package hasher

type HashFunc64Type func([]byte) ([]uint64, int)
type HashFunc32Type func([]byte) ([]uint32, int)

type HashScheme64Type func(hs *[]uint64, hn int, seed, capacity uint) uint64
type HashScheme32Type func(hs *[]uint32, hn int, seed, capacity uint) uint32

var hashScheme64 = map[string]HashScheme64Type{
	"enhanced_double_hashing": Enhanced64BitDoubleHashing,
	"modified_kirmit":         Modified64BitKirMit,
}

var hashFunc64 = map[string]HashFunc64Type{
	"murmur3_256": MurmurHash256,
	"murmur3_128": MurmurHash128,
}

func GetHashFunc64(hName string) HashFunc64Type {
	return hashFunc64[hName]
}

func GetHashScheme64(sName string) HashScheme64Type {
	return hashScheme64[sName]
}
