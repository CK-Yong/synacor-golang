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
	// Program counter
	index uint16
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
	for {
		op := vm.memory[vm.index][0]

		// Prepare operands, get from registry if necessary
		operands := make([]uint16, len(vm.memory[vm.index])-1)
		for i, operand := range vm.memory[vm.index][1:] {
			val := vm.checkRegister(operand)
			operands[i] = val
		}

		switch op {
		case 0: // stop
			return nil
		case 1:
			vm.set(operands[0], operands[1])
			break
		case 6: // jmp
			vm.index = operands[0]
			break
		case 7:
			vm.jt(operands[0], operands[1])
			break
		case 8:
			vm.jf(operands[0], operands[1])
			break
		case 9:
			vm.add(operands[0], operands[1], operands[2])
			break
		case 19:
			vm.out(operands[0])
			break
		case 21: // no-op
			vm.index++
			break
		default:
			vm.index++
		}
	}
}

func (vm *VirtualMachine) checkRegister(arg uint16) uint16 {
	if arg >= 32768 && arg <= 32775 { // Is within register values
		index := arg - 32768
		return vm.register[index]
	}

	return arg
}

// assign into <a> the sum of <b> and <c> (modulo 32768)
func (vm *VirtualMachine) add(a uint16, b uint16, c uint16) {
	vm.register[a] = (b + c) % 32768
	vm.index++
}

// set register <a> to the value of <b>
func (vm *VirtualMachine) set(a uint16, b uint16) {
	vm.register[a] = b
	vm.index++
}

// if <a> is nonzero, jump to <b>
func (vm *VirtualMachine) jt(a uint16, b uint16) {
	if a != 0 {
		vm.index = b
		return
	}
	vm.index++
}

// if <a> is zero, jump to <b>. Returns the destination index
func (vm *VirtualMachine) jf(a uint16, b uint16) {
	if a == 0 {
		vm.index = b
		return
	}
	vm.index++
}

func (vm *VirtualMachine) out(arg uint16) {
	fmt.Printf("%c", arg)
	vm.index++
}
