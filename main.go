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

	vm := VirtualMachine.Initialize()
	vm.Execute(file)

	if err != nil {
		panic("Error occurred during execution" + err.Error())
	}
}
