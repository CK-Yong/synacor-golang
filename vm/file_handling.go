package VirtualMachine

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func (vm *VirtualMachine) DumpMemory() [32768]uint16 {
	result := [32768]uint16{}
	for address, value := range vm.memory {
		result[address] = value
	}
	return result
}

func (vm *VirtualMachine) load(file *os.File) error {
	index := 0
	for {
		num, err := readInt(file)

		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Successfully loaded file...")
				return nil
			} else {
				return err
			}
		}

		vm.memory[index] = num
		index++
	}
}

// Reads and returns the identifier for the operation
func readInt(file *os.File) (uint16, error) {
	intBuff := [2]byte{}
	_, err := io.ReadFull(file, intBuff[:])

	if err != nil {
		return 0, err
	}

	instruction := binary.LittleEndian.Uint16(intBuff[:])

	return instruction, nil
}