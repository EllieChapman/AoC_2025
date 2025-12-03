package day2

import (
	"AoC_2025/src/utils"
	"strings"
)

func Day3_part1(input []string) int {
	banks := utils.Map(input, func(bank string) []int { return utils.Map(strings.Split(bank, ""), utils.Atoi) })
	joltages := utils.Map(banks, func(b []int) int { return getMaxJoltageN(b, 2, 0) })
	return utils.Sum(joltages)
}

func Day3_part2(input []string) int {
	banks := utils.Map(input, func(bank string) []int { return utils.Map(strings.Split(bank, ""), utils.Atoi) })
	joltages := utils.Map(banks, func(b []int) int { return getMaxJoltageN(b, 12, 0) })
	return utils.Sum(joltages)
}

func getMaxJoltageN(bank []int, numLeftToFind int, found int) int {
	next, pos := findBiggest(bank[0 : len(bank)-numLeftToFind+1])
	newFound := found*10 + next
	if numLeftToFind == 1 {
		return newFound // we are done, just found the last digit
	} else {
		return getMaxJoltageN(bank[pos+1:], numLeftToFind-1, newFound) // More digits to find, recurse
	}
}

func findBiggest(bank []int) (int, int) {
	for target := 9; target >= 0; target-- {
		for pos, value := range bank {
			if value == target {
				return target, pos
			}
		}
	}
	panic("Didn't match digit of any value")
}
