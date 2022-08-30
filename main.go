package main

import (
	"encoding/binary"
	"fmt"
	"github.com/ckyong/synacor/VirtualMachine"
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

	bytes, err := getBytesFromFile(file, err)
	if err != nil {
		panic("Could not read file: " + err.Error())
	}

	vm := VirtualMachine.Initialize()
	vm.Execute(bytes)
}

func getBytesFromFile(file *os.File, err error) (*[]byte, error) {
	empty := &[]byte{}
	// Get file size, read it into the buffer
	stats, statsErr := file.Stat()
	if statsErr != nil {
		return empty, err
	}

	bytes := make([]byte, stats.Size())

	err = binary.Read(file, binary.LittleEndian, bytes)

	if err != nil {
		fmt.Printf("Could not read file to buffer: %v", err)
		return empty, err
	}
	return &bytes, nil
}
