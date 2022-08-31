package main

import (
	"encoding/binary"
	"fmt"
	"io"
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

	for {
		integer, err := readInt(file)
		if err != nil {
			fmt.Printf("Error: %v", err)
			break
		}
		if integer <= 21 {
			println()
		}
		print(integer)
		print(" ")
	}
}

func readInt(file *os.File) (uint16, error) {
	intBuff := [2]byte{}
	_, err := io.ReadFull(file, intBuff[:])

	if err != nil {
		return 0, err
	}

	result := binary.LittleEndian.Uint16(intBuff[:])

	return result, nil
}
