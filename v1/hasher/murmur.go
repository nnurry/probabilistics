package hasher

func MurmurHash256(data []byte) ([]uint64, int) {
	var d digest128
	hash1, hash2, hash3, hash4 := d.sum256(data)
	return []uint64{hash1, hash2, hash3, hash4}, 4
}

// extract from from sum256() in murmur_bnb.go
func MurmurHash128(data []byte) ([]uint64, int) {
	var d digest128
	d.h1, d.h2 = 0, 0
	// Process as many bytes as possible.
	d.bmix(data)
	// We have enough to compute the first two 64-bit numbers
	length := uint(len(data))
	tail_length := length % block_size
	tail := data[length-tail_length:]
	hash1, hash2 := d.sum128(false, length, tail)
	return []uint64{hash1, hash2}, 2
}
