package vm

import (
	"fmt"
)

const (
	MEMORY_SIZE = 0xFF

	trapVecPtr    = 0x0000 // trap vector table
	intVecTable   = 0x0100 // interupt vector table
	supervStack   = 0x0200 // supervisor stack / operating system
	userProgram   = 0x3000 // user programs
	deviceRegAddr = 0xFE00 // device register addresses
)

type VirtualMachine interface {
	ExecuteInstruction() error
}

type machine struct {
	memory    [MEMORY_SIZE]uint16
	registers [16]uint16
}

func NewMachine() VirtualMachine {
	m := &machine{}
	// 1 5 -1
	// 0001 101 100000001
	m.memory[0] = 0b0001101100000001
	return m
}

// Fetches the next instruction from memory and executes it
func (m *machine) ExecuteInstruction() error {
	// Fetch instruction and increment pc
	ins := m.memory[m.registers[reg_PC]]
	m.registers[reg_PC]++

	// Parse instruction
	opcode := bitSequence(ins, 12, 4)
	regA := bitSequence(ins, 9, 3) // destination register (mostly)
	regB := bitSequence(ins, 6, 3) // source/base register (mostly)
	pcOffset := signExtend(bitSequence(ins, 0, 9), 9)

	fmt.Println(regA, regB, pcOffset)

	switch opcode {
	case op_ADD:
		// Bit [5] is 0 for register addition, 1 for immediate value
		addType := bitAt(ins, 5)
		result := uint16(0) // keep result for flag check after
		if addType == 0 {
			regC := bitSequence(ins, 0, 3) // second source reg
			result = m.registers[regB] + m.registers[regC]
			m.registers[regA] = result
		} else {
			// immediate := signExtend(bitSequence(ins, 0, 5), 5)
			// result := m.registers[regB] + immediate
		}
	}

	return nil
}
