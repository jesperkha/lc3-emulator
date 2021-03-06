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
		{
			"JSR",
			[]uint16{
				0b0100100000000011, // JSR 3
			},
			reg_PC,
			4,
		},
		{
			"LD",
			[]uint16{
				0b0001000000100101, // Dummy instruction
				0b0010001111111110, // LD R1, -2
			},
			reg_R1,
			0b0001000000100101,
		},
		{
			"LDI",
			[]uint16{
				0b1010000000000000, // LDI R0, 0
				0b0000000000000010, // 2 (address of next)
				0b0000000001000101, // 69
			},
			reg_R0,
			69,
		},
		{
			"LDR",
			[]uint16{
				0b0001000000100010, // ADD R0, 2
				0b0110001000000000, // LDR R1, R0 0
				0b0000000000000011, // 3
			},
			reg_R1,
			3,
		},
		{
			"LEA",
			[]uint16{
				0b1110000000000001, // LEA R0, 1
			},
			reg_R0,
			2,
		},
		{
			"NOT",
			[]uint16{
				0b0001000000101111, // ADD R0, R0 1111
				0b1001001000000000, // NOT R1, R0
			},
			reg_R1,
			0b1111111111110000,
		},
		{
			"ST",
			[]uint16{
				0b0001000000100001, // ADD R0, R0 1
				0b0011000000001000, // ST R0, 8
				0b0110001000001001, // LDR R1, R0 9
			},
			reg_R1,
			1,
		},
		{
			"STR",
			[]uint16{
				0b0001000000100001, // ADD R0, R0 1
				0b0111000001000011, // STR R0, R1 3
				0b0110000001000011, // LDR R0, R1 3
			},
			reg_R0,
			1,
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
