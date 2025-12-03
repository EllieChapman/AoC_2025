package day2

import (
	"AoC_2025/src/utils"
	"fmt"
	"strings"
)

func Day3_part1(input []string) int {
	banks := utils.Map(input, splitBatteries)
	joltages := utils.Map(banks, getMaxJoltage2Wrapper)
	return utils.Sum(joltages)
}

func Day3_part2(input []string) int {
	banks := utils.Map(input, splitBatteries)
	joltages := utils.Map(banks, getMaxJoltage12Wrapper)
	return utils.Sum(joltages)
}

func splitBatteries(bank string) []int {
	return utils.Map(strings.Split(bank, ""), utils.Atoi)
}

func getMaxJoltage2Wrapper(bank []int) int {
	return getMaxJoltageN(bank, 2, 0)
}

func getMaxJoltage12Wrapper(bank []int) int {
	return getMaxJoltageN(bank, 12, 0)
}

// EBC how can I directly use utils.Map with getMaxJoltageN?
// Currently can't because it needs other parameters, not just the one being mapped over, so need wrapper functions
func getMaxJoltageN(bank []int, numLeftToFind int, found int) int {
	next, pos := findBiggest(bank[0 : len(bank)-numLeftToFind+1])
	newFound := found*10 + next
	if numLeftToFind == 1 {
		// we are done, just found the last digit
		return newFound
	} else {
		// more digits to find, recurse
		return getMaxJoltageN(bank[pos+1:], numLeftToFind-1, newFound)
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
	fmt.Println("should never hit here")
	return 0, 0
}
