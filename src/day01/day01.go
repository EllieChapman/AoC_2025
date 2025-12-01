package day1

import (
	"AoC_2025/src/utils"
	"fmt"
)

var lines []string = utils.ReadLines("src/day01/input.txt")

func Day1_part1() int {

	rotations := utils.Map(lines, parseLine)
	positions := doClicks(rotations)

	return countZeros(positions)
}

func Day1_part2() int {

	rotations := utils.Map(lines, parseLine)
	oneClickRotations := utils.Flatten(utils.Map(rotations, squash))
	positions := doClicks(oneClickRotations)

	return countZeros(positions)
}

type rotation struct {
	direction string
	clicks    int
}

func parseLine(line string) rotation {
	r := rotation{}
	r.direction = line[0:1]
	if r.direction != "R" && r.direction != "L" {
		fmt.Print("unknown direction", r.direction)
	}
	r.clicks = utils.Atoi(line[1:])
	return r
}

func doClicks(rs []rotation) []int {
	positions := []int{}
	current := 50
	for _, value := range rs {
		current = doClick(current, value)
		positions = append(positions, current)
	}
	return positions
}

func doClick(current int, r rotation) int {
	var pos int
	if r.direction == "R" {
		pos = current + r.clicks
	} else {
		pos = current - r.clicks
	}

	// get within 0 to 99 again
	pos = fixRange(pos)

	return pos
}

func fixRange(pos int) int {
	for pos > 99 {
		pos = pos - 100
	}
	for pos < 0 {
		pos = pos + 100
	}
	return pos
}

func countZeros(is []int) int {
	total := 0
	for _, value := range is {
		if value == 0 {
			total++
		}
	}
	return total
}

func squash(r rotation) []rotation {
	rs := []rotation{}
	for r.clicks > 0 {
		rs = append(rs, rotation{r.direction, 1})
		r.clicks--
	}
	return rs
}
