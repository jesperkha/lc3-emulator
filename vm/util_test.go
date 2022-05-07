package vm

import (
	"testing"
)

func TestZeroExtend(t *testing.T) {
	n := uint16(0b1011000011111110)
	extended := zeroExtend(n, 5)
	expect := 30 // 11110
	if uint16(expect) != extended {
		t.Errorf("expected %d, got %d", expect, extended)
	}
}

func TestSignExtend(t *testing.T) {
	n := uint16(0b0000000000010011)
	signed := signExtend(n, 5)
	expect := -3
	if int16(expect) != signed {
		t.Errorf("expected %d, got %d", expect, signed)
	}
}

func TestBitSequence(t *testing.T) {
	n := uint16(0xABCD)
	seq := bitSequence(n, 4, 8)
	expect := 0x00C0
	if seq != uint16(expect) {
		t.Errorf("expected %d, got %d", expect, seq)
	}
}
