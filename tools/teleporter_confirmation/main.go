/*
Notes:
From the debugging logs (and a few external hints), I found that there is a repeating function as soon as the
teleporter is used while the register is set to non-zero:
```
vmreg: [4 1 3 10 101 0 0 1] vmstack [6080 16 6124 1 2952 25978 3568 3599 2708 5445 3]
5489: call 6027
```

Index 6027 is the start of the repeating function. Afterwards, we see that the function returns to line 5491 (see the
stack)
```
vmreg: [4 1 3 10 101 0 0 1] vmstack [6080 16 6124 1 2952 25978 3568 3599 2708 5445 3 5491]
6027: jt 32768 6035
...
vmreg: [2 4 3 10 101 0 0 1] vmstack [6080 16 6124 1 2952 25978 3568 3599 2708 5445 3 5491 4 6056 6047 6067 2]
6054: call 6027
```

Disassembling the binary code allows us to see that this function will run until a result is found that equals 6, with
(4,1) as starting values
```
5483: set 32768 4
5486: set 32769 1
5489: call 6027
5491: eq 32769 32768 6
```

The function looks as follows:
```
6027: jt 32768 6035				if r0 != 0 { goto 6035 }
// r0 is zero
6030: add 32768 32769 1			r0 = r1 + 1
6034: ret						return r0
6035: jt 32769 6048				if r1 != 0 { goto 6048 }

// r1 is zero
6038: add 32768 32768 32767		r0--
6042: set 32769 32775			r1 = r8
6045: call 6027					recurse...
6047: ret

// Both r0 and r1 are not zero
6048: push 32768				stack.push(r0)
6050: add 32769 32769 32767		r1--
6054: call 6027					recurse...

// After returning from r0 != 0 && r1 != 0
6056: set 32769 32768			r1 = r0
6059: pop 32768					stack.pop(r0)
6061: add 32768 32768 32767		r0--
6065: call 6027					recurse...
6067: ret
```

Attempted solutions:
- Using memoization --> Note: had to reduce cache for A to 5 (~1m10s), and had to recreate it for every run
(as r7 was changing the results)
- Adding concurrency (~41s)
- Improving function using induction:
```
where a = r0, b = r1, c = r7
f(0,b) = b + 1
f(a,0) = f(a-1, c)
f(a,b) = f(a-1, f(a, b-1))

f(4,1) = f(4-1, f(4,1-1)) = f(3, f(4,0)) = f(3, f(4-1, c)) = f(3, f(3, c))
```
*/

package main

import (
	"fmt"
	"os"
	"time"
)

func main() {

	//r0 := uint16(4)
	//r1 := uint16(4)
	r7 := uint16(4)

	start := time.Now()
	for ; r7 < 32768; r7++ {
		// Always recreate cache, as r7 changes the entire scenario for each run
		cache := make([][]int, 4) // r0 never grows more than length of 5
		for i := range cache {
			cache[i] = make([]int, 32768)
			for j := range cache[i] {
				cache[i][j] = -1
			}
		}

		r7 := r7
		//go func() { // Run calculations in parallel until we find a solution
		//	result := calculate(cache, uint16(r0), uint16(r1), uint16(r7))
		//	fmt.Printf("(4, 1, %v) = %v\n", r7, result)
		//	if result == 6 {
		//		fmt.Printf("Register 7: %v.\n", r7)
		//		fmt.Printf("Finished in %s.", time.Since(start))
		//		os.Exit(0)
		//	}
		//}()

		// Alternative approach: f(3, f(3, c))
		go func() {
			num := calculate(cache, 3, r7, r7)
			result := calculate(cache, 3, num, r7)
			fmt.Printf("(4, 1, %v) = %v\n", r7, result)
			if result == 6 {
				fmt.Printf("Register 7: %v.\n", r7)
				fmt.Printf("Finished in %s.", time.Since(start))
				os.Exit(0)
			}
		}()
	}
}

func calculate(cache [][]int, a uint16, b uint16, c uint16) uint16 {
	if cache[a][b] != -1 {
		return uint16(cache[a][b])
	}

	var result uint16

	if a != 0 {
		if b != 0 {
			num := calculate(cache, a, b-1, c)
			cache[a][b-1] = int(num)
			result = calculate(cache, a-1, num, c)
			cache[a-1][num] = int(result)
			return result
		} else {
			result = calculate(cache, a-1, c, c)
			cache[a-1][c] = int(result)
			return result
		}
	} else {
		result = (b + 1) % 32768
		cache[a][b] = int(result)
		return result
	}
}
