package day04

import (
	"AoC_2025/src/utils"
)

func Day4_part1(input []string) int {
	coords, m := parseCoords(input)
	return utils.Sum(utils.Map(coords, func(c coord) int { return testCoord(c, m) }))
}

func Day4_part2(input []string) int {
	coords, m := parseCoords(input)
	return countRemoves(0, coords, m)
}

type coord struct {
	x int
	y int // this in inverted from normal graphs, bigger numbers more down
}

// func countRemoves(countRemoved int, cs []coord, m map[coord]string) int {
// 	for _, c := range cs {
// 		if testCoord(c, m) == 1 {
// 			m[c] = "."
// 			return countRemoves(countRemoved+1, cs, m)
// 		}
// 	}
// 	return countRemoved
// }

func countRemoves(countRemoved int, cs []coord, m map[coord]string) int {
	removing := []coord{}
	for _, c := range cs {
		if testCoord(c, m) == 1 {
			removing = append(removing, c)
		}
	}
	if len(removing) == 0 {
		// we are done, this round could not remove any more
		return countRemoved
	} else {
		countRemoved += len(removing)
		for _, c := range removing {
			m[c] = "."
		}
		return countRemoves(countRemoved, cs, m)
	}
}

func parseCoords(lines []string) ([]coord, map[coord]string) {
	cs := []coord{}
	m := map[coord]string{}
	for yPos, line := range lines {
		for xPos, value := range line {
			newC := coord{xPos, yPos}
			cs = append(cs, newC)
			m[newC] = string(value)
		}
	}
	return cs, m
}

func testCoord(c coord, m map[coord]string) int {
	if m[c] == "." {
		return 0
	}

	adjacentCoords := generateAdjacent(c)
	numRollsAjacent := utils.Sum(utils.Map(adjacentCoords, func(ajc coord) int {
		if m[ajc] == "@" {
			return 1
		} else {
			return 0
		}
	}))
	if numRollsAjacent > 3 {
		return 0 // cannot be reached/removed
	} else {
		return 1 // can be reached/removed
	}
}

// EBC improve this
func generateAdjacent(c coord) []coord {
	x, y := c.x, c.y
	return []coord{
		{x + 1, y + 1},
		{x + 1, y},
		{x + 1, y - 1},
		{x, y + 1},
		{x - 1, y + 1},
		{x - 1, y},
		{x - 1, y - 1},
		{x, y - 1},
	}
}
