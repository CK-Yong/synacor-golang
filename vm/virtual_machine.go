package VirtualMachine

import (
	"fmt"
	"os"
)

type VirtualMachine struct {
	// 15 bit address space memory
	memory [32768]uint16
	// Register set with 8 slots
	register [8]uint16
	// Unbounded stack
	stack Stack
	// Program counter
	index  uint16
	opArgs map[uint16]uint16
}

type Stack struct {
	inner []uint16
}

func (stack *Stack) push(arg uint16) {
	stack.inner = append(stack.inner, arg)
}

func (stack *Stack) pop() uint16 {
	val := stack.inner[len(stack.inner)-1]
	stack.inner = stack.inner[:len(stack.inner)-1]
	return val
}

// checks whether address refers to the VM registry, and writes it either to the registry or the corresponding memory address.
func (vm *VirtualMachine) write(address uint16, val uint16) {
	if index, ok := tryGetRegistryAddress(address); ok {
		vm.register[index] = vm.tryGetRegistryValue(val)
	} else {
		vm.memory[address] = vm.tryGetRegistryValue(val)
	}
}

func Load(file *os.File) (*VirtualMachine, error) {
	vm := VirtualMachine{
		memory:   [32768]uint16{},
		register: [8]uint16{},
		stack:    Stack{inner: []uint16{}},
		opArgs: map[uint16]uint16{
			0:  0,
			1:  2,
			2:  1,
			3:  1,
			4:  3,
			5:  3,
			6:  1,
			7:  2,
			8:  2,
			9:  3,
			10: 3,
			11: 3,
			12: 3,
			13: 3,
			14: 2,
			15: 2,
			16: 2,
			17: 1,
			18: 0,
			19: 1,
			20: 1,
			21: 0,
		},
	}
	err := vm.load(file)

	if err != nil {
		return nil, err
	}

	return &vm, nil
}

func (vm *VirtualMachine) Run() error {
	defer func() {
		fmt.Printf("Fault index: %v\n", vm.index)
	}()

	for {
		op := vm.memory[vm.index]

		// Prepare operands, get from registry if necessary
		var operands []uint16
		if vm.opArgs[op] > 0 {
			operands = vm.memory[vm.index+1 : vm.index+vm.opArgs[op]+1]
		}

		switch op {
		case 0: // stop
			return nil
		case 1:
			vm.set(operands[0], operands[1])
			break
		case 2:
			vm.push(operands[0])
			break
		case 3:
			vm.pop(operands[0])
			break
		case 4:
			vm.eq(operands[0], operands[1], operands[2])
			break
		case 5:
			vm.gt(operands[0], operands[1], operands[2])
		case 6: // jmp
			vm.jmp(operands[0])
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
		case 12:
			vm.and(operands[0], operands[1], operands[2])
			break
		case 13:
			vm.or(operands[0], operands[1], operands[2])
			break
		case 14:
			vm.not(operands[0], operands[1])
			break
		case 17:
			vm.call(operands[0])
			break
		case 19:
			vm.out(operands[0])
			break
		case 21: // no-op
			vm.index++
			break
		default:
			fmt.Printf("Unknown operation: %v at index %v\n", op, vm.index)
			vm.index++
		}
	}
}

// returns the registry address, or the default value (passed in value) if arg is not a registry address.
func tryGetRegistryAddress(arg uint16) (uint16, bool) {
	if arg >= 32768 && arg <= 32775 { // Is within register values
		return arg - 32768, true
	}

	return arg, false
}

// returns the value of the registry, or the default value (passed in value) if arg does not point to a registry address.
func (vm *VirtualMachine) tryGetRegistryValue(arg uint16) uint16 {
	index, isRegistry := tryGetRegistryAddress(arg)

	if isRegistry {
		return vm.register[index]
	}

	return arg
}

// set register <a> to the value of <b>
func (vm *VirtualMachine) set(index uint16, b uint16) {
	index, _ = tryGetRegistryAddress(index)
	vm.register[index] = vm.tryGetRegistryValue(b)
	vm.index += 3
}

func (vm *VirtualMachine) push(a uint16) {
	vm.stack.push(vm.tryGetRegistryValue(a))
	vm.index += 2
}

// remove the top element from the stack and write it into <a>; empty stack = error
func (vm *VirtualMachine) pop(a uint16) {
	val := vm.stack.pop()
	vm.write(a, val)
	vm.index += 2
}

// set <a> to 1 if <b> is equal to <c>; set it to 0 otherwise
func (vm *VirtualMachine) eq(a uint16, b uint16, c uint16) {
	if vm.tryGetRegistryValue(b) == vm.tryGetRegistryValue(c) {
		vm.write(a, 1)
	}
	vm.index += 4
}

// set <a> to 1 if <b> is greater than <c>; set it to 0 otherwise
func (vm *VirtualMachine) gt(a uint16, b uint16, c uint16) {
	if vm.tryGetRegistryValue(b) > vm.tryGetRegistryValue(c) {
		vm.write(a, 1)
	} else {
		vm.write(a, 0)
	}
	vm.index += 4
}

// jump to <a>
func (vm *VirtualMachine) jmp(a uint16) {
	newValue := vm.tryGetRegistryValue(a)
	vm.index = newValue
}

// if <a> is nonzero, jump to <b>
func (vm *VirtualMachine) jt(a uint16, b uint16) {
	if vm.tryGetRegistryValue(a) != 0 {
		vm.jmp(vm.tryGetRegistryValue(b))
		return
	}
	vm.index += 3
}

// if <a> is zero, jump to <b>. Returns the destination index
func (vm *VirtualMachine) jf(a uint16, b uint16) {
	if vm.tryGetRegistryValue(a) == 0 {
		vm.jmp(vm.tryGetRegistryValue(b))
		return
	}
	vm.index += 3
}

// assign into <a> the sum of <b> and <c> (modulo 32768)
func (vm *VirtualMachine) add(a uint16, b uint16, c uint16) {
	index, _ := tryGetRegistryAddress(a)
	vm.register[index] = (vm.tryGetRegistryValue(b) + vm.tryGetRegistryValue(c)) % 32768
	vm.index += 4
}

// stores into <a> the bitwise and of <b> and <c>
func (vm *VirtualMachine) and(a uint16, b uint16, c uint16) {
	and := vm.tryGetRegistryValue(b) & vm.tryGetRegistryValue(c)
	vm.write(a, vm.tryGetRegistryValue(and))
	vm.index += 4
}

// stores into <a> the bitwise or of <b> and <c>
func (vm *VirtualMachine) or(a uint16, b uint16, c uint16) {
	or := vm.tryGetRegistryValue(b) | vm.tryGetRegistryValue(c)
	vm.write(a, vm.tryGetRegistryValue(or))
	vm.index += 4
}

// stores 15-bit bitwise inverse of <b> in <a>
func (vm *VirtualMachine) not(a uint16, b uint16) {
	not := ^vm.tryGetRegistryValue(b)
	vm.write(a, vm.tryGetRegistryValue(not>>1))
	vm.index += 3
}

// write the address of the next instruction to the stack and jump to <a>
func (vm *VirtualMachine) call(a uint16) {
	vm.stack.push(vm.index + 2)
	vm.jmp(a)
}

func (vm *VirtualMachine) out(a uint16) {
	fmt.Printf("%c", vm.tryGetRegistryValue(a))
	vm.index += 2
}
