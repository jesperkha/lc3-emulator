package vm

import (
	"errors"
)

var (
	errUnknownOpcode = errors.New("unknown opcode: %d")
)

const (
	MEMORY_SIZE = 0xF // Debug

	trapVecPtr    = 0x0000 // trap vector table
	intVecTable   = 0x0100 // interupt vector table
	supervStack   = 0x0200 // supervisor stack / operating system
	userProgram   = 0x3000 // user programs
	deviceRegAddr = 0xFE00 // device register addresses
)

type VirtualMachine interface {
}

type machine struct {
	memory    []uint16
	registers [16]uint16
}

func NewMachine() VirtualMachine {
	return &machine{
		memory: make([]uint16, MEMORY_SIZE),
	}
}

// Sets the appropriate flags based on the result value v
// Todo: unset flags
func (m *machine) setFlags(v uint16) {
	if (v & 0x0001) == 0 {
		m.registers[reg_FLAG] |= fl_PTY
	}
	if v == 0 {
		m.registers[reg_FLAG] |= fl_ZRO
	}
	if (v & 0x8000) == 1 {
		m.registers[reg_FLAG] |= fl_NEG
	}
}

// Returns true if the flag is set
func (m *machine) hasFlag(v uint16) bool {
	return m.registers[reg_FLAG]&v != 0
}

// Fetches the next instruction from memory and executes it
func (m *machine) executeInstruction() error {
	// Fetch instruction and increment pc
	ins := m.memory[m.registers[reg_PC]]
	m.registers[reg_PC]++

	// Parse instruction
	opcode := bitSequence(ins, 12, 4)
	regA := bitSequence(ins, 9, 3) // destination register (mostly)
	regB := bitSequence(ins, 6, 3) // source/base register (mostly)
	pcOffset := signExtend(bitSequence(ins, 0, 9), 9)

	switch opcode {
	case op_ADD: // ADD <dest>, <source1> <source2|immediate>
		addType := bitAt(ins, 6) // 0 for register addition, 1 for immediate value
		result := uint16(0)
		if addType == 0 {
			regC := bitSequence(ins, 0, 3) // second source register
			result = m.registers[regB] + m.registers[regC]
		} else {
			immediate := signExtend(bitSequence(ins, 0, 5), 5)
			result = m.registers[regB] + immediate
		}
		m.registers[regA] = result
		m.setFlags(result)

	case op_AND: // AND <dest>, <source1> <source2|immediate>
		andType := bitAt(ins, 6) // 0 for register AND, 1 for immediate value
		result := uint16(0)
		if andType == 0 {
			regC := bitSequence(ins, 0, 3) // second source register
			result = m.registers[regB] & m.registers[regC]
		} else {
			immediate := signExtend(bitSequence(ins, 0, 5), 5)
			result = m.registers[regB] & immediate
		}
		m.registers[regA] = result
		m.setFlags(result)

	case op_BR: // BR[flags] <label>
		p, z, n := bitAt(ins, 10) > 0, bitAt(ins, 11) > 0, bitAt(ins, 12) > 0 // flags
		if m.hasFlag(fl_POS) && p || m.hasFlag(fl_ZRO) && z || m.hasFlag(fl_NEG) && n {
			m.registers[reg_PC] += pcOffset
		}

	case op_JMP: // JMP <source> or RET
		offsetReg := regB
		if regB == 7 { // source2 is 111 for return
			offsetReg = reg_R7 // R7 holds return instruction after subroutine call
		}
		m.registers[reg_PC] = m.registers[offsetReg]

	default:
		return errorf(errUnknownOpcode, opcode)
	}

	return nil
}
