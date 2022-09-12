package main

import (
	"fmt"
	VirtualMachine "github.com/ckyong/synacor/vm"
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

	vm, err := VirtualMachine.LoadDebugger(file)
	if err != nil {
		panic(err)
	}

	err = vm.Run()

	if err != nil {
		panic("Error occurred during execution" + err.Error())
	}
}
