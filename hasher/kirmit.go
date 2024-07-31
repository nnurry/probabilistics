package hasher

// Kirsch-Mitzenmacher = KirMit

// Kirsch-Mitzenmacher for accomodating variable-sized hash slice (just made it up, don't know if it holds valid)
func Modified64BitKirMit(hs *[]uint64, hn int, seed, capacity uint) uint64 {
	capacity64 := uint64(capacity)
	paddedSeed64 := uint64(seed + 3)
	if hn == 1 {
		return (*hs)[0] + paddedSeed64
	}
	finalHash := (*hs)[0]
	paddedPowerSeed64 := paddedSeed64
	for i := 1; i < hn; i++ {
		finalHash += (paddedPowerSeed64 * (*hs)[i]) % capacity64
		paddedPowerSeed64 *= paddedSeed64
	}
	return finalHash
}

// my version of KirMit's algorithm on n 64-bit hashes
func NaiveNHash64BitKirMit(hs *[]uint64, hn int, seed uint, capacity uint) uint64 {
	seed64 := uint64(seed)
	capacity64 := uint64(capacity)

	hLen := uint64(hn)

	halfHLen := hLen / 2

	// x % 2^n = x & (2^n - 1)
	seed64Mod2 := seed64 & 1
	hLenMod2 := hLen & 1

	// pick 1st-half index of hash slice
	firstHashIdx := seed64 % halfHLen
	// pick 2nd-half index of hash slice
	// initial value = 1st index, then further mix it
	secondHashIdx := firstHashIdx
	// mix asymmetry
	secondHashIdx += seed64Mod2
	// limit index range = [halfHLen, hLen-1]
	// case when hLen is odd, last hash may never be used -> mix even/odd
	secondHashIdx = halfHLen + (hLenMod2 & seed64Mod2) + (secondHashIdx % halfHLen)

	return (*hs)[firstHashIdx] + seed64*(*hs)[secondHashIdx]%capacity64
}
