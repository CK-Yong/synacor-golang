package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
Illustration of the grid
  * ─── 8 ─── - ─── 1		3
  │     │     │     │
  │     │     │     │
  4 ─── + ───11 ─── *		2
  │     │     │     │
  │     │     │     │
  + ─── 4 ─── - ───18		1
  │     │     │     │
  │     │     │     │
 22 ─── - ─── 9 ─── *		0

  0		1	  2     3
Constraints:
- The weight cannot be lower than 0
- The weight cannot be higher than 32768
- If we arrive at 22, the orb needs to be picked up again
- If we arrive at 1, the orb shatters
- We need to find the shortest route that allows for this
*/

func main() {
	/*
		We can create a path by determining whether we increase the x coordinate or y coordinate by 1.
		For the shortest route, we need to have a path that in total increases x by 3, and y by 3 (otherwise we can't reach the end)
		Generate a base route, evaluate whether it adds up to 30, and if not, permute the steps and try it.

		todo: Is there operator precedence??
	*/

	grid := [][]string{ // Mirror the grid so we make it a bit easier to use the coordinates.
		{"22", "-", "9", "*"},
		{"+", "4", "-", "18"},
		{"4", "+", "11", "*"},
		{"*", "8", "-", "1"},
	}

	// base
	steps := []step{{1, 0}, {1, 0}, {1, 0}, {0, 1}, {0, 1}, {0, 1}, {1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	// Step in X direction added

	// populate this array with permutations
	permutations := [][]step{steps}
	permute(&permutations, steps, len(steps))

	//addedX := false
	//addedY := false

	cache := map[string]int{}

	//for
	{
		for _, route := range permutations {
			weight := traverse(cache, route, grid)
			if weight == 30 {
				fmt.Printf("Found a valid route %v\n", route)
				os.Exit(0)
			} else {
				fmt.Printf("%v: %v\n", route, weight)
			}
		}

		// todo: generate search tree and traverse it
	}
}

func permute(permutations *[][]step, steps []step, index int) {
	// Use Heap's algorithm to generate permutations
	if index == 1 {
		*permutations = append(*permutations, clone(steps))
	} else {
		permute(permutations, clone(steps), index-1)
		for i := 0; i < index-1; i++ {
			var swapped []step
			if index%2 == 0 { // even
				swapped = swap(clone(steps), i, index-1)
			} else {
				clone(steps)
				swapped = swap(clone(steps), 0, index-1)
			}
			permute(permutations, swapped, index-1)
		}
	}
}

func clone(array []step) []step {
	newArray := make([]step, len(array))
	copy(newArray, array)
	return newArray
}

func swap(array []step, a int, b int) []step {
	clone := array
	temp := clone[a]
	clone[a] = clone[b]
	clone[b] = temp
	return clone
}

type step struct {
	x int
	y int
}

func traverse(cache map[string]int, steps []step, grid [][]string) int {
	weight := 22
	operation := ""
	loc := step{0, 0}

	if val, ok := cache[serialize(steps)]; ok {
		return val
	}

	for _, change := range steps {
		loc.x += change.x
		loc.y += change.y

		if loc.x < 0 || loc.x > 3 || loc.y < 0 || loc.y > 3 {
			cache[serialize(steps)] = -1
			return -1 // out of bounds
		}

		tile := grid[loc.y][loc.x]
		if tile == "+" || tile == "-" || tile == "*" {
			operation = tile
			continue
		}

		val, _ := strconv.Atoi(tile)

		switch operation {
		case "+":
			weight += val
			break
		case "-":
			weight -= val
			break
		case "*":
			weight *= val
			break
		}

		if weight <= 0 || weight > 32768 {
			// The orb breaks
			cache[serialize(steps)] = -1
			return -1
		}

		if loc.x == 3 && loc.y == 3 {
			cache[serialize(steps)] = weight
			return weight // finished
		}

		if loc.x == 0 && loc.y == 0 {
			cache[serialize(steps)] = -1
			return -1 // orb shattered
		}
	}

	return weight
}

func serialize(steps []step) string {
	result := strings.Builder{}
	for _, item := range steps {
		result.WriteByte(byte(item.x))
		result.WriteByte(byte(item.y))
	}
	return result.String()
}
