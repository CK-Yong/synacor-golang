package VirtualMachine

import (
	"fmt"
	"os"
)

type VirtualMachineDebugger struct {
	inner     *VirtualMachine
	outputLog bool
}

func LoadDebugger(file *os.File) (*VirtualMachineDebugger, error) {
	vm, err := Load(file)
	if err != nil {
		return nil, err
	}

	return &VirtualMachineDebugger{
		inner: vm,
	}, nil
}

func (vm *VirtualMachineDebugger) Run() error {
	defer func() {
		fmt.Printf("Fault index: %v\n", vm.inner.Index)
	}()

	for {
		op := vm.inner.Memory[vm.inner.Index]

		// Prepare operands, get from registry if necessary
		var operands []uint16
		if vm.inner.opArgs[op] > 0 {
			operands = vm.inner.Memory[vm.inner.Index+1 : vm.inner.Index+vm.inner.opArgs[op]+1]
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
			if err := vm.pop(operands[0]); err != nil {
				return err
			}
			break
		case 4:
			vm.eq(operands[0], operands[1], operands[2])
			break
		case 5:
			vm.gt(operands[0], operands[1], operands[2])
			break
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
		case 10:
			vm.mult(operands[0], operands[1], operands[2])
			break
		case 11:
			vm.mod(operands[0], operands[1], operands[2])
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
		case 15:
			vm.rmem(operands[0], operands[1])
			break
		case 16:
			vm.wmem(operands[0], operands[1])
			break
		case 17:
			vm.call(operands[0])
			break
		case 18:
			if err := vm.ret(); err != nil {
				return err
			}
			break
		case 19:
			vm.out(operands[0])
			break
		case 20:
			vm.in(operands[0])
			break
		case 21: // no-op
			vm.inner.Index++
			break
		default:
			fmt.Printf("Unknown operation: %v at index %v\n", op, vm.inner.Index)
			vm.inner.Index++
		}
	}
}

func (vm *VirtualMachineDebugger) print(op string, args ...uint16) {
	if !vm.outputLog {
		return
	}

	fmt.Println("vmreg:", vm.inner.Register, "vmstack", vm.inner.Stack.inner)

	fmt.Printf("%v: %v ", vm.inner.Index, op)
	for _, arg := range args {
		fmt.Printf("%v ", arg)
	}
	fmt.Print("\n")
}

// set register <a> to the value of <b>
func (vm *VirtualMachineDebugger) set(a uint16, b uint16) {
	vm.print("set", a, b)
	vm.inner.set(a, b)
}

func (vm *VirtualMachineDebugger) push(a uint16) {
	vm.print("push", a)
	vm.inner.push(a)
}

// remove the top element from the stack and write it into <a>; empty stack = error
func (vm *VirtualMachineDebugger) pop(a uint16) error {
	vm.print("pop", a)
	return vm.inner.pop(a)
}

// set <a> to 1 if <b> is equal to <c>; set it to 0 otherwise
func (vm *VirtualMachineDebugger) eq(a uint16, b uint16, c uint16) {
	vm.print("eq", a, b, c)
	vm.inner.eq(a, b, c)
}

// set <a> to 1 if <b> is greater than <c>; set it to 0 otherwise
func (vm *VirtualMachineDebugger) gt(a uint16, b uint16, c uint16) {
	vm.print("gt", a, b, c)
	vm.inner.gt(a, b, c)
}

// jump to <a>
func (vm *VirtualMachineDebugger) jmp(a uint16) {
	vm.print("jmp", a)
	vm.inner.jmp(a)
}

// if <a> is nonzero, jump to <b>
func (vm *VirtualMachineDebugger) jt(a uint16, b uint16) {
	vm.print("jt", a, b)
	vm.inner.jt(a, b)
}

// if <a> is zero, jump to <b>. Returns the destination index
func (vm *VirtualMachineDebugger) jf(a uint16, b uint16) {
	vm.print("jf", a, b)
	vm.inner.jf(a, b)
}

// assign into <a> the sum of <b> and <c> (modulo 32768)
func (vm *VirtualMachineDebugger) add(a uint16, b uint16, c uint16) {
	vm.print("add", a, b, c)
	vm.inner.add(a, b, c)
}

// store into <a> the product of <b> and <c> (modulo 32768)
func (vm *VirtualMachineDebugger) mult(a uint16, b uint16, c uint16) {
	vm.print("mult", a, b, c)
	vm.inner.mult(a, b, c)
}

// store into <a> the remainder of <b> divided by <c>
func (vm *VirtualMachineDebugger) mod(a uint16, b uint16, c uint16) {
	vm.print("mod", a, b, c)
	vm.inner.mod(a, b, c)
}

// stores into <a> the bitwise and of <b> and <c>
func (vm *VirtualMachineDebugger) and(a uint16, b uint16, c uint16) {
	vm.print("and", a, b, c)
	vm.inner.and(a, b, c)
}

// stores into <a> the bitwise or of <b> and <c>
func (vm *VirtualMachineDebugger) or(a uint16, b uint16, c uint16) {
	vm.print("or", a, b, c)
	vm.inner.or(a, b, c)
}

// stores 15-bit bitwise inverse of <b> in <a>
func (vm *VirtualMachineDebugger) not(a uint16, b uint16) {
	vm.print("not", a, b)
	vm.inner.not(a, b)
}

// read memory at address <b> and write it to <a>
func (vm *VirtualMachineDebugger) rmem(a uint16, b uint16) {
	vm.print("rmem", a, b)
	vm.inner.rmem(a, b)
}

// write the value from <b> into memory at address <a>
func (vm *VirtualMachineDebugger) wmem(a uint16, b uint16) {
	vm.print("wmem", a, b)
	vm.inner.wmem(a, b)
}

// write the address of the next instruction to the stack and jump to <a>
func (vm *VirtualMachineDebugger) call(a uint16) {
	vm.print("call", a)
	vm.inner.call(a)
}

// remove the top element from the stack and jump to it; empty stack = halt
func (vm *VirtualMachineDebugger) ret() error {
	vm.print("ret")
	return vm.inner.ret()
}

// read a character from the terminal and write its ascii code to <a>
func (vm *VirtualMachineDebugger) in(a uint16) {
	vm.outputLog = true
	vm.print("in", a)
	vm.inner.in(a)
}

// write the character represented by ascii code <a> to the terminal
func (vm *VirtualMachineDebugger) out(a uint16) {
	vm.print("out", a)
	vm.inner.Index += 2
}
