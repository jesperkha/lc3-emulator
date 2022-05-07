package vm

import (
	"testing"
)

func TestZeroExtend(t *testing.T) {
	n := uint16(0xFFFF)
	expect := uint16(0x000F)
	extended := zeroExtend(n, 4)
	if uint16(expect) != extended {
		t.Errorf("expected %d, got %d", expect, extended)
	}
}

func TestSignExtend(t *testing.T) {
	n := uint16(0x0010)
	expect := uint16(0xFFF0)
	extended := signExtend(n, 5)
	if expect != extended {
		t.Errorf("expected %d, got %d", expect, extended)
	}
}

func TestBitSequence(t *testing.T) {
	n := uint16(0x0F00)
	expect := uint16(0x000F)
	seq := bitSequence(n, 8, 4)
	if seq != uint16(expect) {
		t.Errorf("expected %d, got %d", expect, seq)
	}
}

func TestBitAt(t *testing.T) {
	n := uint16(0x0010)
	expect := 1
	at := bitAt(n, 5)
	if at != uint16(expect) {
		t.Errorf("expected %d, got %d", expect, at)
	}
}
