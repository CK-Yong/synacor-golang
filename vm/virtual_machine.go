package VirtualMachine

import (
	"fmt"
	"os"
)

type VirtualMachine struct {
	// 15 bit address space memory
	memory [32768][]uint16
	// Register set with 8 slots
	register [8]uint16
	// Unbounded stack
	stack []uint16
}

func Load(file *os.File) (*VirtualMachine, error) {
	vm := VirtualMachine{
		memory:   [32768][]uint16{},
		register: [8]uint16{},
		stack:    make([]uint16, 32),
	}
	err := vm.load(file)

	if err != nil {
		return nil, err
	}

	return &vm, nil
}

func (vm *VirtualMachine) Run() error {
	index := 0
	for {
		op := vm.memory[index][0]

		// Prepare operands, get from registry if necessary
		operands := make([]uint16, len(vm.memory[index])-1)
		for i, operand := range vm.memory[index][1:] {
			val := vm.checkRegister(operand)
			operands[i] = val
		}

		switch op {
		case 0: // stop
			return nil
		case 6: // jmp
			// New index - 1 (zero based)
			index = int(operands[0])
			break
		case 8:
			index = jf(index, operands[0], operands[1])
			break
		case 9:
			vm.add(operands[0], operands[1], operands[2])
			index++
			break
		case 19:
			out(operands[0])
			index++
			break
		case 21: // no-op
			index++
			break
		default:
			index++
		}
	}
}

func (vm *VirtualMachine) checkRegister(arg uint16) uint16 {
	if arg > 32768 && arg <= 32775 { // Is within register values
		index := arg - 32768
		return vm.register[index]
	}

	return arg
}

// assign into <a> the sum of <b> and <c> (modulo 32768)
func (vm *VirtualMachine) add(a uint16, b uint16, c uint16) {
	vm.register[a] = (b + c) % 32768
}

/*
if <a> is zero, jump to <b>. Returns the destination index
*/
func jf(currentIndex int, a uint16, b uint16) int {
	if a == 0 {
		return int(b)
	}
	return currentIndex + 1 // return the next index
}

func out(arg uint16) {
	fmt.Printf("%c", arg)
}
