package vm

import (
	"errors"
)

var (
	errUnknownOpcode          = errors.New("unknown opcode: %d")
	errPrivilegeModeViolation = errors.New("privilege mode violation")
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

// Sets the appropriate flags based on the result value in register r
// Todo: unset flags
func (m *machine) setFlags(r uint16) {
	v := m.registers[r]
	if (v & 0x0001) == 0 {
		m.registers[reg_FLAG] |= fl_POS
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
	pcOffset := signedSequence(ins, 0, 9)

	switch opcode {
	case op_ADD:
		// ADD <dest>, <source1> <source2|immediate>
		addType := bitAt(ins, 6) // 0 for register addition, 1 for immediate value
		if addType == 0 {
			regC := bitSequence(ins, 0, 3) // second source register
			m.registers[regA] = m.registers[regB] + m.registers[regC]
		} else {
			immediate := signedSequence(ins, 0, 5)
			m.registers[regA] = m.registers[regB] + immediate
		}
		m.setFlags(regA)

	case op_AND:
		// AND <dest>, <source1> <source2|immediate>
		andType := bitAt(ins, 6) // 0 for register AND, 1 for immediate value
		if andType == 0 {
			regC := bitSequence(ins, 0, 3) // second source register
			m.registers[regA] = m.registers[regB] & m.registers[regC]
		} else {
			immediate := signedSequence(ins, 0, 5)
			m.registers[regA] = m.registers[regB] & immediate
		}
		m.setFlags(regA)

	case op_BR:
		// BR[flags] <label>
		p, z, n := bitAt(ins, 10) > 0, bitAt(ins, 11) > 0, bitAt(ins, 12) > 0 // flags
		if m.hasFlag(fl_POS) && p || m.hasFlag(fl_ZRO) && z || m.hasFlag(fl_NEG) && n {
			m.registers[reg_PC] += pcOffset
		}

	case op_JMP:
		// JMP <source> | RET
		offsetReg := regB
		if regB == 7 { // source2 is 111 for return
			offsetReg = reg_R7 // R7 holds return instruction after subroutine call
		}
		m.registers[reg_PC] = m.registers[offsetReg]

	case op_JSR:
		// JSR  <label> | JSRR <source>
		m.registers[reg_R7] = m.registers[reg_PC] // set return point to current pc at R7
		if bitAt(ins, 12) == 1 {                  // add offset to pc
			offset := signedSequence(ins, 0, 11)
			m.registers[reg_PC] += offset
		} else {
			m.registers[reg_PC] = m.registers[regB]
		}

	case op_LD:
		// LD <dest>, <label>
		m.registers[regA] = m.memory[m.registers[reg_PC]+pcOffset]
		m.setFlags(regA)

	case op_LDI:
		// LDI <dest>, <label>
		memPtr := m.memory[m.registers[reg_PC]+pcOffset] // load loc is at mem[pc+offset]
		m.registers[regA] = m.memory[memPtr]             // load value
		m.setFlags(regA)

	case op_LDR:
		// LDR <dest>, <base>, <offset>
		offset := signedSequence(ins, 0, 6)  // get offset
		memPtr := m.registers[regB] + offset // memptr is value of baseR + offset
		m.registers[regA] = m.memory[memPtr]
		m.setFlags(regA)

	case op_LEA:
		// LEA <dest>, <label>
		m.registers[regA] = m.registers[reg_PC] + pcOffset // store target mem address
		m.setFlags(regA)

	case op_NOT:
		// NOT <dest>, <source>
		m.registers[regA] = ^m.registers[regB]
		m.setFlags(regA)

	case op_RTI:
		// RTI
		if m.registers[reg_PSR]&0x8000 == 1 {
			return errorf(errPrivilegeModeViolation)
		}
		// else privilege mode and condition restored after interrupt
		m.registers[reg_PC] = m.memory[m.registers[reg_R6]] // R6 is the SSP
		temp := m.memory[m.registers[reg_R6]+1]
		m.registers[reg_R6] += 2
		m.registers[reg_PSR] = temp

	case op_ST:
		// ST <source>, <label>
		memPtr := m.registers[reg_PC] + pcOffset
		m.memory[memPtr] = m.registers[regA]

	case op_STI:
		// STI <source>, <label>
		memPtr := m.memory[m.registers[reg_PC]+pcOffset]
		m.memory[memPtr] = m.registers[regA]

	case op_STR:
		// STR <source>, <base>, <offset>
		memPtr := m.registers[regB] + signedSequence(ins, 0, 6)
		m.memory[memPtr] = m.registers[regA]

	default:
		return errorf(errUnknownOpcode, opcode)
	}

	return nil
}
