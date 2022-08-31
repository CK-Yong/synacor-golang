package VirtualMachine

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type VirtualMachine struct {
	memory     map[int16]uint16
	register   [8]uint16
	stack      []uint16
	fileOffset int
	funMap     map[byte]func(*VirtualMachine, *os.File)
}

func Initialize() VirtualMachine {
	return VirtualMachine{
		memory:     make(map[int16]uint16),
		register:   [8]uint16{},
		stack:      make([]uint16, 32),
		fileOffset: 0,
	}
}

func (vm *VirtualMachine) Execute(file *os.File) error {
	_, err := file.Seek(0, 0)
	if err != nil {
		return err
	}

	for {
		instruction, err := vm.readInt(file)

		if err != nil {
			return err
		}

		switch instruction {
		case 0:
			return nil
		case 19:
			vm.out(file)
		case 21:
			continue
		}
	}

	return nil
}

// Reads and returns the identifier for the operation
func (vm *VirtualMachine) readInt(file *os.File) (uint16, error) {
	intBuff := [2]byte{}
	bytesRead, err := io.ReadFull(file, intBuff[:])

	if err != nil {
		return 0, err
	}

	vm.fileOffset += bytesRead
	instruction := binary.LittleEndian.Uint16(intBuff[:])

	return instruction, nil
}

// Reads multiple integers and returns them as an array
func (vm *VirtualMachine) readInts(file *os.File, n int) ([]uint16, error) {
	ret := make([]uint16, n)
	for i := 0; i < n; i++ {
		result, err := vm.readInt(file)

		if err != nil {
			return []uint16{}, err
		}

		ret = append(ret, result)
	}

	return ret, nil
}

func (vm *VirtualMachine) out(file *os.File) {
	arg, _ := vm.readInt(file)
	fmt.Printf("%c", arg)
}
