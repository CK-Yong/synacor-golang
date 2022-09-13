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
- Moved cache generation out of loop, and used a 'static' template to copy an initialized array over the cache every
run (~30s).
- Improving function using induced function
```
where a = r0, b = r1, c = r7
f(0,b) = b + 1
f(a,0) = f(a-1, c)
f(a,b) = f(a-1, f(a, b-1))

inductions (with some help from the internet)
f(1,b) = b + c + 1
f(2,b) = (2c + 1) + b * (c + 1) --> (~7s)
f(3,b) = f(3, b - 1)(c + 1) + (2c + 1)
f(4,1) = f(3, f(3,c))
```
From here we can generate a more efficient function
```
f(3,0) = f(2,c)
f(3,b) = f(3, b - 1)(c + 1) + (2c + 1) // Add this value to the array, plug it into the value, then add to array.
f(4,1) = f(3, f(3,c)) // After adding f(3,b) to the array, check if this calculated value equals 6
```
This was not much faster than just replacing f2 as the more efficient function. Rough runtime was (~6 seconds)
*/

package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	start := time.Now()
	// Start with c = 1, as 0 is already the default teleporter state
	cache := make([]uint16, 32768)
	for c := uint16(1); c < 32768; c++ {
		// Add f(3,0) = f(2,c) as the base case
		cache[0] = f2(c, c)

		// f(4,1) = f(3, f(3,c))
		// calculate f(3,b) for every value of b
		for b := uint16(1); b < 32768; b++ {
			// f(3,b) = f(3, b - 1)(c + 1) + (2c + 1)
			cache[b] = (cache[b-1]*(c+1) + 2*c + 1) % 32768
		}

		// f(4,1) = f(3, f(3,c))
		// The result is the f(3,c)-th spot
		result := cache[cache[c]]

		fmt.Printf("(4, 1, %v) = %v\n", c, result)
		if result == 6 {
			fmt.Printf("Register 7: %v.\n", c)
			break
		}
	}
	fmt.Printf("Finished in %s.", time.Since(start))
	os.Exit(0)
}

// f(2,b) = (2c + 1) + b * (c + 1)
func f2(b uint16, c uint16) uint16 {
	return (2*c + 1 + b*(c+1)) % 32768
}
