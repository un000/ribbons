package ribbons

func setBit(n uint64, pos uint64) uint64 {
	return n | (1 << pos)
}

func delBit(n uint64, pos uint64) uint64 {
	return n &^ (1 << pos)
}

func hasBit(n uint64, pos uint64) bool {
	return (n & (1 << pos)) > 0
}

func extractToggledBits(n uint64, add uint64) []uint64 {
	if n == 0 {
		return nil
	}

	toggled := make([]uint64, 0, 1)
	for i := uint64(0); i < 64; i++ {
		if hasBit(n, i) {
			toggled = append(toggled, add+i)
		}
	}

	return toggled
}
