package vm

import (
	"testing"
)

type TestOpcode struct {
	name         string
	instructions []uint16
	resultReg    uint16
	expect       uint16
}

func TestOperations(t *testing.T) {
	m := machine{
		memory: make([]uint16, 256),
	}

	testCases := []TestOpcode{
		{
			"ADD",
			[]uint16{
				0b0001000000100101, // ADD R0, R0  5
				0b0001000000111111, // ADD R0, R0 -1
			},
			reg_R0,
			4,
		},
		{
			"AND",
			[]uint16{
				0b0001000000101101, // ADD R0, R0 0b10101
				0b0101001000111111, // AND R0, R0 0b11111
			},
			reg_R0,
			0x000D,
		},
		{
			"BR",
			[]uint16{
				0b0001000000000000, // ADD R0, R0 R0
				0b0000010000000101, // BRz 5
			},
			reg_PC,
			7,
		},
		{
			"JMP",
			[]uint16{
				0b0001000000101000, // ADD R0, R0 8
				0b1100000000000000, // JMP R0
			},
			reg_PC,
			8,
		},
	}

	for _, c := range testCases {
		m.registers[reg_PC] = 0
		copy(m.memory, c.instructions)           // copy instructions into memory
		copy(m.registers[:], make([]uint16, 16)) // clear registers

		for i := 0; i < len(c.instructions); i++ {
			if err := m.executeInstruction(); err != nil {
				t.Errorf("operation %s failed at instruction %d, err %s", c.name, i+1, err)
			}
		}

		value := m.registers[c.resultReg]
		if value != c.expect {
			t.Errorf("operation %s expected %d, got %d", c.name, c.expect, value)
		}
	}
}
