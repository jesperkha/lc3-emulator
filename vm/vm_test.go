package vm

import "testing"

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
	}

	for _, c := range testCases {
		m.registers[reg_PC] = 0
		copy(m.memory, c.instructions)
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
