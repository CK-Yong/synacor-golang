package main

import (
	"fmt"
	"github.com/ckyong/synacor/vm"
	"os"
	"path/filepath"
)

func main() {
	filePath := filepath.Join("./resources/challenge.bin")

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Could not close file: %v", err)
		}
	}(file)

	vm, err := VirtualMachine.Load(file)

	opArgs := map[uint16]uint16{
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
	}

	opName := map[uint16]string{
		0:  "halt",
		1:  "set",
		2:  "push",
		3:  "pop",
		4:  "eq",
		5:  "gt",
		6:  "jmp",
		7:  "jt",
		8:  "jf",
		9:  "add",
		10: "mult",
		11: "mod",
		12: "and",
		13: "or",
		14: "not",
		15: "rmem",
		16: "wmem",
		17: "call",
		18: "ret",
		19: "out",
		20: "in",
		21: "noop",
	}

	for {
		if int(vm.Index) == len(vm.Memory) {
			return
		}

		op := vm.Memory[vm.Index]

		if op > 21 {
			fmt.Printf("%v: %v\n", vm.Index, op)
			vm.Index++
			continue
		}

		fmt.Printf("%v: %v ", vm.Index, opName[op])

		var operands []uint16
		if opArgs[op] > 0 {
			operands = vm.Memory[vm.Index+1 : vm.Index+opArgs[op]+1]

			for i := uint16(0); i < opArgs[op]; i++ {
				fmt.Printf("%v ", operands[i])
			}
		}

		fmt.Println()

		vm.Index += opArgs[op] + 1
	}
}
