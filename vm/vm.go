package vm

const (
	MEMORY_SIZE = 0xFFFF

	trapVecPtr    = 0x0000 // trap vector table
	intVecTable   = 0x0100 // interupt vector table
	supervStack   = 0x0200 // supervisor stack / operating system
	userProgram   = 0x3000 // user programs
	deviceRegAddr = 0xFE00 // device register addresses
)

type machine struct {
	memory    [MEMORY_SIZE]uint16
	registers [16]uint16
}
