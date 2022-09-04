package VirtualMachine

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func (vm *VirtualMachine) DumpMemory() [32768][]uint16 {
	result := [32768][]uint16{}
	for address, value := range vm.memory {
		result[address] = make([]uint16, len(value))
		copy(result[address], value)
	}
	return result
}

func (vm *VirtualMachine) load(file *os.File) error {
	index := 0
	var args []uint16
	for {
		op, err := readInt(file)

		if op < 21 && len(vm.memory[index]) > 0 {
			// This is an op that needs a new line before we can continue (to prevent rewriting old bytes)
			index++
		}

		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Successfully loaded file...")
				return nil
			} else {
				return err
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
			vm.memory[index] = append(vm.memory[index], op)
			continue
		}

		vm.memory[index] = createOpSlice(op, args...)
		index++
	}
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
