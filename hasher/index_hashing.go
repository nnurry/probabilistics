package hasher

// http://www.peterd.org/pcd-diss.pdf
// Adaptive Approximate State Storage
// 6.5.4 Enhanced double hashing
func EnhancedDoubleHashing(hs *[]uint64, hn int, seed, capacity uint) uint64 {
	seed64 := uint64(seed)
	capacity64 := uint64(capacity)
	if seed == 0 {
		return (*hs)[0]
	}
	(*hs)[0] = ((*hs)[0] + (*hs)[1]) % capacity64
	(*hs)[1] = ((*hs)[1] + seed64) % capacity64
	return (*hs)[0]
}
