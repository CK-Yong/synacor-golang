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

	memDump := vm.DumpMemory()

	for _, values := range memDump {
		for _, value := range values {
			fmt.Printf("%v ", value)
		}
		fmt.Println()
	}
}
