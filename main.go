package main

import "github.com/jesperkha/smore/vm"

// https://justinmeiners.github.io/lc3-vm/supplies/lc3-isa.pdf

func main() {
	m := vm.NewMachine()
	m.ExecuteInstruction()
}
