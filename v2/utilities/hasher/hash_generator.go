package hasher

type HashGenerator[T HashOutType] struct {
	hashFunction   HashFunction[T]
	hashFamily     string
	platformBit    uint
	outputBit      uint
	generateMethod string
}

func NewHashGenerator[T HashOutType](hashFamily string, platformBit uint, outputBit uint, generateMethod string) (*HashGenerator[T], error) {
	hashFunction, err := NewHashFunction[T](hashFamily, platformBit, outputBit)
	if err != nil {
		return nil, err
	}
	hashGenerator := &HashGenerator[T]{
		hashFunction:   hashFunction,
		hashFamily:     hashFamily,
		platformBit:    platformBit,
		outputBit:      outputBit,
		generateMethod: generateMethod,
	}

	return hashGenerator, nil
}

func (g *HashGenerator[T]) GenerateHash(data []byte, seed interface{}, hashCeil uint, times uint) ([]T, error) {
	output := make([]T, times)
	seedT := seed.(T)
	hashCeilT := T(hashCeil)
	hashes, err := g.hashFunction(data, seed)
	if err != nil {
		return nil, err
	}

	output = append(output, hashes[0])

	if len(hashes) < 2 {
		// only standard
		for i := uint(1); i < times; i++ {
			hashes, err = g.hashFunction(data, seedT+T(i))
			if err != nil {
				return nil, err
			}
			output = append(output, hashes[0])
		}
		return output, nil
	}

	if g.generateMethod == "extended-double-hashing" {
		// http://www.peterd.org/pcd-diss.pdf
		// Adaptive Approximate State Storage
		// 6.5.4 Enhanced double hashing

		for i := uint(1); i < times; i++ {
			newSeedT := seedT + T(i)
			hashes[0] = (hashes[0] + hashes[1]) % hashCeilT
			hashes[1] = (hashes[1] + newSeedT) % hashCeilT
			output = append(output, hashes[0])
		}
		return output, nil
	} else if g.generateMethod == "kirsch-mitzenmacher" {
		// Kirsch-Mitzenmacher for accomodating variable-sized hash slice (just made it up, don't know if it holds valid)
		seedT += 3
		for i := uint(0); i < times; i++ {
			finalHash := hashes[0]
			newSeedT := seedT + T(i)
			powerSeedT := newSeedT
			for j := 1; j < len(hashes); j++ {
				finalHash += (powerSeedT * hashes[i]) % hashCeilT
				powerSeedT *= newSeedT
			}
			output = append(output, finalHash)
		}
		return output, nil
	}
	// standard: k-hash functions -> hash k-times with different seed
	output = append(output, hashes[1:]...)
	i := uint(len(hashes))
	for i < times-1 {
		hashes, err = g.hashFunction(data, seedT+T(i))
		if err != nil {
			return nil, err
		}
		for _, hash := range hashes {
			output = append(output, hash)
			i++
			if i == times-1 {
				return output, nil
			}
		}
	}
	return output, nil
}
