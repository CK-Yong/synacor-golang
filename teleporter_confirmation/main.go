package main

import "fmt"

type Register struct {
	r0 uint16
	r1 uint16
	r8 uint16
}

func main() {

	r0 := 4
	r1 := 1
	r8 := 0

	result := runAckermannFunction(uint16(r0), uint16(r1), uint16(r8))
	fmt.Println(result)
}

func runAckermannFunction(a uint16, b uint16, c uint16) uint16 {
	if a != 0 {
		if b != 0 {
			return runAckermannFunction(a-1, runAckermannFunction(a, c, c), c)
		} else {
			return runAckermannFunction(a-1, 1, 0)
		}
	} else {
		return b + 1
	}
}

//func runAckermannFunction(register *Register) uint16 {
//	if register.r0 != 0 {
//		if register.r1 != 0 {
//			register.r1--
//			return runAckermannFunction(register)
//		} else {
//			register.r0--
//			register.r1 = register.r8
//			return runAckermannFunction(register)
//		}
//	} else {
//		return register.r1 + 1
//	}
//}
