package hash

// hashes := filterSize/elements * ln(2). There we assume elements=1.
// see https://en.wikipedia.org/wiki/Bloom_filter#Optimal_number_of_hash_functions
func NewHashList(filterSize uint32) []Hash {
	listLen := xLn2(filterSize)
	list := make([]Hash, listLen)

	var initialSeed int32 = int32(1337_420)
	for i := 0; i < int(listLen); i++ {
		list[i] = NewSipHash(uint64(initialSeed), uint64(initialSeed))
		initialSeed++
	}

	return list
}

// hashes := filterSize/elements * ln(2)
// see https://en.wikipedia.org/wiki/Bloom_filter#Optimal_number_of_hash_functions
func UpdateList(list []Hash, fsize, elts uint32) []Hash {
	listLen := xLn2(fsize) / elts
	if listLen < uint32(len(list)) {
		list = list[:listLen]
	}
	return list
}

func xLn2(x uint32) uint32 {
	return uint32(int64(x) * 69314 / 100000)
}
