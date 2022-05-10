package vm

const (
	// Registers, used as indecies of the register slice
	reg_R0 = iota
	reg_R1
	reg_R3
	reg_R4
	reg_R5
	reg_R6 // stack pointer
	reg_R7
	reg_PC
	reg_FLAG
	reg_PSR
)

const (
	op_BR = iota // Branch
	op_ADD
	op_LD
	op_ST
	op_JSR
	op_AND
	op_LDR
	op_STR
	op_RTI
	op_NOT
	op_LDI
	op_STI
	op_JMP
	_ // reserved
	op_LEA
	op_TRAP
)

const (
	fl_POS = 1 << iota
	fl_ZRO
	fl_NEG
)
