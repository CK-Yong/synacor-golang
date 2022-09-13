package main

import (
	"fmt"
	"strconv"
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
	*/

	grid := [][]string{
		{"*", "8", "-", "1"},
		{"4", "+", "11", "*"},
		{"+", "4", "-", "18"},
		{"22", "-", "9", "*"},
	}

	steps := []step{{1, 0}, {1, 0}, {1, 0}, {0, 1}, {0, 1}, {0, 1}}

	// populate this array with permutations
	var allPermutations [][]step
	permute(steps)

	weight := -1
	for weight != 30 {
		weight = traverse(steps, grid)

		// steps = permute(steps) OR calculate all permutations beforehand and traverse them.
	}

	fmt.Printf("Found a valid route %v\n", steps)
}

type step struct {
	x int
	y int
}

func traverse(steps []step, grid [][]string) int {
	current := 22
	operation := ""
	for _, step := range steps {
		tile := grid[step.x][step.y]
		if tile == "+" || tile == "-" || tile == "*" {
			operation = tile
			continue
		}

		val, _ := strconv.Atoi(tile)

		switch operation {
		case "+":
			current += val
			break
		case "-":
			current -= val
			break
		case "*":
			current *= val
			break
		}

		if current < 0 || current > 32768 {
			// The orb breaks
			return -1
		}
	}

	return current
}
