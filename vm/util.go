package vm

import "fmt"

// Shorthand for fmt.Errorf() to not have 10 mile long lines
func errorf(err error, a ...any) error {
	return fmt.Errorf("[ ERROR ]"+err.Error(), a...)
}

// Logs formatted message
func logf(msg string, a ...any) {
	fmt.Printf("[ LOG ] %s\n", fmt.Sprintf(msg, a...))
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
	sign := bitAt(n, length)
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

// Shorthand for sign-extended bit sequence
func signedSequence(n uint16, start int, length int) uint16 {
	return signExtend(bitSequence(n, start, length), length)
}

// Returns the bit at the given pos, starting at 1 from the right
func bitAt(n uint16, pos int) uint16 {
	return 0x0001 & (n >> (pos - 1))
}
