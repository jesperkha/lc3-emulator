package vm

import "fmt"

// Shorthand for fmt.Errorf() to not have 10 mile long lines
func errorf(err error, a ...any) error {
	return fmt.Errorf(err.Error(), a...)
}

// Returns a 16 bit number from the right most bits of length p.
// Example: (0xFFFF, 4) -> 0x000F
func zeroExtend(n uint16, p int) uint16 {
	x := uint16(0xFFFF >> (16 - p))
	return n & x
}

// Replicates the most significant bit until 16 bits long.
// Example: (0x00F0, 8) -> 0xFFF0
func signExtend(n uint16, length int) uint16 {
	sign := bitAt(n, length-1)
	if sign == 0 {
		return zeroExtend(n, length)
	}

	return (0xFFFF << (length - 1)) | n
}

// Returns the number in the given bit range starting from the right.
// Example: (0xABCD, 4, 8) -> 0x000C
func bitSequence(n uint16, start int, length int) uint16 {
	a := ^uint16(0xFFFF << length)
	b := n >> (start)
	return a & b
}

// Returns the bit at the given pos (15 <- 0)
func bitAt(n uint16, pos int) uint16 {
	return 0x0001 & (n >> pos)
}
