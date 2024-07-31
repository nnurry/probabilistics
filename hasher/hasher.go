package hasher

type HashFunc64Type func([]byte) ([]uint64, int)
type HashFunc32Type func([]byte) ([]uint32, int)

var hashFunc64 = map[string]HashFunc64Type{
	"murmur3_256": MurmurHash256,
	"murmur3_128": MurmurHash128,
}

func GetHashFunc64(hName string) HashFunc64Type {
	return hashFunc64[hName]
}
