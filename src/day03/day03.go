package day3

import (
	"AoC_2025/src/utils"
	"fmt"
	"strconv"
	"strings"
)

func Day3_part1(input []string) int {
	banks := utils.Map(input, func(bank string) []int { return utils.Map(strings.Split(bank, ""), utils.Atoi) })
	joltages := utils.MapVariadic(banks, getMaxJoltageN, 2, "0") // Nb should both be ints, but wanted to prove MapVaridic capabilities
	return utils.Sum(joltages)
}

func Day3_part2(input []string) int {
	banks := utils.Map(input, func(bank string) []int { return utils.Map(strings.Split(bank, ""), utils.Atoi) })
	joltages := utils.MapVariadic(banks, getMaxJoltageN, 12, "0")
	fmt.Println()
	return utils.Sum(joltages)
}

func getMaxJoltageN(bank []int, ii ...any) int {
	numLeftToFind := ii[0].(int)
	found := utils.Atoi(ii[1].(string))
	next, pos := findBiggest(bank[0 : len(bank)-numLeftToFind+1])
	newFound := found*10 + next
	if numLeftToFind == 1 {
		return newFound // we are done, just found the last digit
	} else {
		return getMaxJoltageN(bank[pos+1:], numLeftToFind-1, strconv.Itoa(newFound)) // More digits to find, recurse
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

// EBC way to make Map work with structs of type any

// func getMaxJoltageNWrapper(bank []int, ii II) int {
// 	return getMaxJoltageN(bank, ii.a, ii.b)
// }

// type Pair[A, B any] struct {
// 	a A
// 	b B
// }

// type II = Pair[int, int]

// call with for example
// joltages := utils.Map1(banks, getMaxJoltageNWrapper, II{2, 0})

// func getMaxJoltageN(bank []int, numLeftToFind int, found int) int {
// 	next, pos := findBiggest(bank[0 : len(bank)-numLeftToFind+1])
// 	newFound := found*10 + next
// 	if numLeftToFind == 1 {
// 		return newFound // we are done, just found the last digit
// 	} else {
// 		return getMaxJoltageN(bank[pos+1:], numLeftToFind-1, newFound) // More digits to find, recurse
// 	}
// }
