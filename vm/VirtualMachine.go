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

	program := toExecutionSet(file)

	for _, operands := range program {
		op := operands[0]

		switch op {
		case 0:
			return nil
		case 19:
			out(operands[1])
		}
	}

	return nil
}

func out(arg uint16) {
	fmt.Printf("%c", arg)
}

func toExecutionSet(file *os.File) [][]uint16 {
	program := make([][]uint16, 0)
	var args []uint16
	for {
		op, err := readInt(file)

		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println(err)
				break
			} else {
				panic(err)
			}
		}

		switch op {
		case 0:
			args = nil
			break
		case 1:
			args, _ = readInts(file, 2)
			break
		case 2:
			args, _ = readInts(file, 1)
			break
		case 3:
			args, _ = readInts(file, 1)
			break
		case 4:
			args, _ = readInts(file, 3)
			break
		case 5:
			args, _ = readInts(file, 3)
			break
		case 6:
			args, _ = readInts(file, 1)
			break
		case 7:
			args, _ = readInts(file, 2)
			break
		case 8:
			args, _ = readInts(file, 2)
			break
		case 9:
			args, _ = readInts(file, 3)
			break
		case 10:
			args, _ = readInts(file, 3)
			break
		case 11:
			args, _ = readInts(file, 3)
			break
		case 12:
			args, _ = readInts(file, 3)
			break
		case 13:
			args, _ = readInts(file, 3)
			break
		case 14:
			args, _ = readInts(file, 2)
			break
		case 15:
			args, _ = readInts(file, 2)
			break
		case 16:
			args, _ = readInts(file, 2)
			break
		case 17:
			args, _ = readInts(file, 1)
			break
		case 18:
			args = nil
			break
		case 19:
			args, _ = readInts(file, 1)
			break
		case 20:
			args, _ = readInts(file, 1)
			break
		case 21:
			args = nil
			break
		default:
			args = nil
			break
		}

		program = append(program, createOpSlice(op, args...))
	}

	return program
}

func createOpSlice(op uint16, args ...uint16) []uint16 {
	result := []uint16{op}

	for _, arg := range args {
		result = append(result, arg)
	}

	return result
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

// Reads multiple integers and returns them as an array
func readInts(file *os.File, n int) ([]uint16, error) {
	ret := make([]uint16, n)
	for i := 0; i < n; i++ {
		result, err := readInt(file)

		if err != nil {
			return []uint16{}, err
		}

		ret[i] = result
	}

	return ret, nil
}
