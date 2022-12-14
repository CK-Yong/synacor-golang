package main

import (
	"bufio"
	"fmt"
	expect "github.com/Netflix/go-expect"
	"log"
	"os"
	"os/exec"
	_ "os/exec"
	"path/filepath"
	"time"
)

func main() {
	console, err := expect.NewConsole(expect.WithStdout(os.Stdout))

	if err != nil {
		log.Fatal(err)
	}
	defer console.Close()

	cmd := exec.Command("synacorvm_linux")
	cmd.Stdin = console.Tty()
	cmd.Stdout = console.Tty()
	cmd.Stderr = console.Tty()

	go func() {
		console.ExpectEOF()
	}()

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	filePath := filepath.Join("./tools/autoplay/autopath.txt")

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

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		time.Sleep(time.Millisecond)
		console.SendLine(scanner.Text())
	}

	if err != nil {
		log.Fatal(err)
	}

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		console.SendLine(input.Text())
	}
}
