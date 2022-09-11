package VirtualMachine

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// todo:
// Improvements:
// 	- Make switch/case into a hashmap of functions
//	- Make registry into a hashmap starting from 32768 to 32775 (so there's no need for conversions)
// 	- Add unit tests for all commands

type VirtualMachine struct {
	// 15 bit address space Memory
	Memory [32768]uint16 `json:"memory"`
	// Register set with 8 slots
	Register [8]uint16 `json:"register"`
	// Unbounded Stack
	Stack Stack `json:"stack"`
	// Program counter
	Index       uint16 `json:"index"`
	opArgs      map[uint16]uint16
	inputBuffer []byte
}

type Stack struct {
	inner []uint16
}

func (stack *Stack) push(arg uint16) {
	stack.inner = append(stack.inner, arg)
}

type EmptyStackError struct{}

func (err *EmptyStackError) Error() string {
	return "Stack is empty, application should halt"
}

// Returns the value from the top of the stack. If a value is returned from the stack, also returns true, otherwise false.
func (stack *Stack) pop() (uint16, error) {
	if len(stack.inner) == 0 {
		return 0, &EmptyStackError{}
	}

	val := stack.inner[len(stack.inner)-1]
	stack.inner = stack.inner[:len(stack.inner)-1]
	return val, nil
}

// checks whether address refers to the VM registry, and writes it either to the registry or the corresponding Memory address.
func (vm *VirtualMachine) write(address uint16, val uint16) {
	if index, ok := tryGetRegistryAddress(address); ok {
		vm.Register[index] = vm.tryGetRegistryValue(val)
	} else {
		vm.Memory[address] = vm.tryGetRegistryValue(val)
	}
}

func Load(file *os.File) (*VirtualMachine, error) {
	vm := VirtualMachine{
		Memory:   [32768]uint16{},
		Register: [8]uint16{},
		Stack:    Stack{inner: []uint16{}},
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
		inputBuffer: []byte{},
	}
	err := vm.load(file)

	if err != nil {
		return nil, err
	}

	return &vm, nil
}

func (vm *VirtualMachine) Run() error {
	defer func() {
		fmt.Printf("Fault index: %v\n", vm.Index)
	}()

	commands := map[uint16]func(operands []uint16) error{
		0: func(operands []uint16) error {
			return nil
		},
		1: func(operands []uint16) error {
			vm.set(operands[0], operands[1])
			return nil
		},
		2: func(operands []uint16) error {
			vm.push(operands[0])
			return nil
		},
		3: func(operands []uint16) error {
			if err := vm.pop(operands[0]); err != nil {
				return err
			}
			return nil
		},
		4: func(operands []uint16) error {
			vm.eq(operands[0], operands[1], operands[2])
			return nil
		},
		5: func(operands []uint16) error {
			vm.gt(operands[0], operands[1], operands[2])
			return nil
		},
		6: func(operands []uint16) error {
			vm.jmp(operands[0])
			return nil
		},
		7: func(operands []uint16) error {
			vm.jt(operands[0], operands[1])
			return nil
		},
		8: func(operands []uint16) error {
			vm.jf(operands[0], operands[1])
			return nil
		},
		9: func(operands []uint16) error {
			vm.add(operands[0], operands[1], operands[2])
			return nil
		},
		10: func(operands []uint16) error {
			vm.mult(operands[0], operands[1], operands[2])
			return nil
		},
		11: func(operands []uint16) error {
			vm.mod(operands[0], operands[1], operands[2])
			return nil
		},
		12: func(operands []uint16) error {
			vm.and(operands[0], operands[1], operands[2])
			return nil
		},
		13: func(operands []uint16) error {
			vm.or(operands[0], operands[1], operands[2])
			return nil
		},
		14: func(operands []uint16) error {
			vm.not(operands[0], operands[1])
			return nil
		},
		15: func(operands []uint16) error {
			vm.rmem(operands[0], operands[1])
			return nil
		},
		16: func(operands []uint16) error {
			vm.wmem(operands[0], operands[1])
			return nil
		},
		17: func(operands []uint16) error {
			vm.call(operands[0])
			return nil
		},
		18: func(operands []uint16) error {
			if err := vm.ret(); err != nil {
				return err
			}
			return nil
		},
		19: func(operands []uint16) error {
			vm.out(operands[0])
			return nil
		},
		20: func(operands []uint16) error {
			vm.in(operands[0])
			return nil
		},
		21: func(operands []uint16) error { // no-op
			vm.Index++
			return nil
		},
	}

	for {
		op := vm.Memory[vm.Index]

		// Prepare operands, get from registry if necessary
		var operands []uint16
		if vm.opArgs[op] > 0 {
			operands = vm.Memory[vm.Index+1 : vm.Index+vm.opArgs[op]+1]
		}

		err := commands[op](operands)

		if err != nil {
			return err
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
		return vm.Register[index]
	}

	return arg
}

// set Register <a> to the value of <b>
func (vm *VirtualMachine) set(index uint16, b uint16) {
	index, _ = tryGetRegistryAddress(index)
	vm.Register[index] = vm.tryGetRegistryValue(b)
	vm.Index += 3
}

func (vm *VirtualMachine) push(a uint16) {
	vm.Stack.push(vm.tryGetRegistryValue(a))
	vm.Index += 2
}

// remove the top element from the Stack and write it into <a>; empty Stack = error
func (vm *VirtualMachine) pop(a uint16) error {
	val, err := vm.Stack.pop()

	if err != nil {
		return err
	}

	vm.write(a, val)
	vm.Index += 2
	return nil
}

// set <a> to 1 if <b> is equal to <c>; set it to 0 otherwise
func (vm *VirtualMachine) eq(a uint16, b uint16, c uint16) {
	if vm.tryGetRegistryValue(b) == vm.tryGetRegistryValue(c) {
		vm.write(a, 1)
	} else {
		vm.write(a, 0)
	}

	vm.Index += 4
}

// set <a> to 1 if <b> is greater than <c>; set it to 0 otherwise
func (vm *VirtualMachine) gt(a uint16, b uint16, c uint16) {
	if vm.tryGetRegistryValue(b) > vm.tryGetRegistryValue(c) {
		vm.write(a, 1)
	} else {
		vm.write(a, 0)
	}

	vm.Index += 4
}

// jump to <a>
func (vm *VirtualMachine) jmp(a uint16) {
	newValue := vm.tryGetRegistryValue(a)
	vm.Index = newValue
}

// if <a> is nonzero, jump to <b>
func (vm *VirtualMachine) jt(a uint16, b uint16) {
	if vm.tryGetRegistryValue(a) != 0 {
		vm.jmp(b)
		return
	}
	vm.Index += 3
}

// if <a> is zero, jump to <b>. Returns the destination Index
func (vm *VirtualMachine) jf(a uint16, b uint16) {
	if vm.tryGetRegistryValue(a) == 0 {
		vm.jmp(b)
		return
	}
	vm.Index += 3
}

// assign into <a> the sum of <b> and <c> (modulo 32768)
func (vm *VirtualMachine) add(a uint16, b uint16, c uint16) {
	index, _ := tryGetRegistryAddress(a)
	vm.Register[index] = (vm.tryGetRegistryValue(b) + vm.tryGetRegistryValue(c)) % 32768
	vm.Index += 4
}

// store into <a> the product of <b> and <c> (modulo 32768)
func (vm *VirtualMachine) mult(a uint16, b uint16, c uint16) {
	index, _ := tryGetRegistryAddress(a)
	vm.Register[index] = (vm.tryGetRegistryValue(b) * vm.tryGetRegistryValue(c)) % 32768
	vm.Index += 4
}

// store into <a> the remainder of <b> divided by <c>
func (vm *VirtualMachine) mod(a uint16, b uint16, c uint16) {
	index, _ := tryGetRegistryAddress(a)
	vm.Register[index] = vm.tryGetRegistryValue(b) % vm.tryGetRegistryValue(c)
	vm.Index += 4
}

// stores into <a> the bitwise and of <b> and <c>
func (vm *VirtualMachine) and(a uint16, b uint16, c uint16) {
	and := vm.tryGetRegistryValue(b) & vm.tryGetRegistryValue(c)
	vm.write(a, and)
	vm.Index += 4
}

// stores into <a> the bitwise or of <b> and <c>
func (vm *VirtualMachine) or(a uint16, b uint16, c uint16) {
	or := vm.tryGetRegistryValue(b) | vm.tryGetRegistryValue(c)
	vm.write(a, or)
	vm.Index += 4
}

// stores 15-bit bitwise inverse of <b> in <a>
func (vm *VirtualMachine) not(a uint16, b uint16) {
	val := vm.tryGetRegistryValue(b)
	// Invert the number and use bitwise & to ensure you're storing a 15 bit number.
	val = ^val & 32767
	vm.write(a, val)

	vm.Index += 3
}

// read Memory at address <b> and write it to <a>
func (vm *VirtualMachine) rmem(a uint16, b uint16) {
	var val uint16
	if adr, isRegistry := tryGetRegistryAddress(b); isRegistry {
		val = vm.Memory[vm.Register[adr]]
	} else {
		val = vm.Memory[b]
	}

	vm.write(a, val)
	vm.Index += 3
}

// write the value from <b> into Memory at address <a>
func (vm *VirtualMachine) wmem(a uint16, b uint16) {
	vm.Memory[vm.tryGetRegistryValue(a)] = vm.tryGetRegistryValue(b)
	vm.Index += 3
}

// write the address of the next instruction to the Stack and jump to <a>
func (vm *VirtualMachine) call(a uint16) {
	vm.Stack.push(vm.Index + 2)
	vm.jmp(a)
}

// remove the top element from the Stack and jump to it; empty Stack = halt
func (vm *VirtualMachine) ret() error {
	val, err := vm.Stack.pop()

	if err != nil {
		return err
	}

	vm.jmp(val)
	return nil
}

// read a character from the terminal and write its ascii code to <a>
func (vm *VirtualMachine) in(a uint16) {
	if len(vm.inputBuffer) == 0 {
		reader := bufio.NewReader(os.Stdin)
		buffer, err := reader.ReadBytes('\n')

		if err != nil {
			fmt.Println("Could not read keyboard input.")
			panic(err)
		}

		vm.inputBuffer = buffer
	}

	// Hacks
	strVal := string(vm.inputBuffer)
	if strings.Contains(strVal, "set") {
		strVal = strings.Trim(strVal, "\n")
		integer, _ := strconv.ParseUint(strings.Split(strVal, " ")[1], 10, 16)
		vm.Register[7] = uint16(integer)
		vm.inputBuffer = []byte{}
		vm.in(a)
		return
	}
	if strings.Contains(strVal, "get") {
		fmt.Printf("R8: %v\n", vm.Register[7])
		vm.inputBuffer = []byte{}
		vm.in(a)
		return
	}

	// save state synacor_1
	if strings.Contains(strVal, "save state") {
		filePath := strings.Trim(strings.Split(strVal, " ")[2], "\n")
		vmJson, err := json.Marshal(vm)

		if err != nil {
			log.Fatal("Could not save state", err)
		}
		err = os.WriteFile(filePath, vmJson, 0777)

		if err != nil {
			log.Fatal("Could not save state", err)
		}
		fmt.Println("Saved state to", filePath)

		vm.inputBuffer = []byte{}
		vm.in(a)
		return
	}

	// load state synacor_1
	if strings.Contains(strVal, "load state") {
		filePath := strings.Trim(strings.Split(strVal, " ")[2], "\n")
		file, err := os.ReadFile(filePath)

		if err != nil {
			log.Fatal("Could not load state", err)
		}

		err = json.Unmarshal(file, vm)

		if err != nil {
			log.Fatal("Could not load state", err)
		}
		fmt.Println("State loaded from", filePath)

		vm.inputBuffer = []byte{}
		vm.in(a)
		return
	}

	// End hacks
	vm.write(a, uint16(vm.inputBuffer[0]))
	vm.inputBuffer = vm.inputBuffer[1:]

	vm.Index += 2
}

// write the character represented by ascii code <a> to the terminal
func (vm *VirtualMachine) out(a uint16) {
	fmt.Printf("%c", vm.tryGetRegistryValue(a))
	vm.Index += 2
}
