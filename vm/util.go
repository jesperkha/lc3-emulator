package vm

// Returns a 16 bit number from the right most bits of length p.
// Example: (0xFFFF, 4) -> 0x000F
func zeroExtend(n uint16, p int) uint16 {
	x := uint16(0xFFFF >> (16 - p))
	return n & x
}

// Returns the zero-extended value but keeps the sign bit.
// Example: 1 0011 -> 1000 0000 0000 0011
func SignExtend(n uint16, length int) int16 {
	signed := 0x0001 & (n >> uint16(length-1)) // get the signed bit
	extended := zeroExtend(n, length-1)        // exclude signed bit
	if signed == 1 {
		return int16(extended) * -1
	}

	return int16(extended)
}

// Returns the number in the given bit range starting from the right.
// Example: (0xABCD, 4, 8) -> 0x00C0
func bitSequence(n uint16, start int, end int) uint16 {
	a := uint16(0xFFFF << start)
	b := uint16(0xFFFF >> end)
	return (a & b) & n
}
