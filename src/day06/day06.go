package day06

import (
	"AoC_2025/src/utils"
	"strings"
)

func Day6_part1(input []string) int {
	sums := utils.Transpose(utils.Map(input, strings.Fields))
	return utils.Sum(utils.Map(sums, doSum))
}

func Day6_part2(input []string) int {
	transposed := utils.Transpose(utils.Map(input[0:len(input)-1], func(s string) []string { return strings.Split(s, "") }))
	sums := extract(transposed, strings.Fields(input[len(input)-1]))
	return utils.Sum(utils.Map(sums, doSum))
}

func doSum(sum []string) int {
	nums := utils.Map(sum[0:len(sum)-1], utils.Atoi)
	op := sum[len(sum)-1]
	switch op {
	case "+":
		return utils.Sum(nums)
	case "*":
		return utils.Mul(nums)
	default:
		panic("unexpected operator:")
	}
}

// [[1    ] [2 4  ] [3 5 6] [     ] [3 6 9] [2 4 8] [8    ] [     ] [  3 2] [5 8 1] [1 7 5] [     ] [6 2 3] [4 3 1] [    4]]
// ->
func extract(ss [][]string, ops []string) [][]string {
	newss := utils.Map(ss, func(s []string) string { return myJoin(s) })
	res := [][]string{}
	outer := 0
	for _, s := range newss {
		words := strings.Fields(s)
		if len(words) == 0 {
			outer++
		}
		if len(words) == 1 {
			if outer >= len(res) {
				res = append(res, []string{})
				res[outer] = append(res[outer], ops[outer])
			}
			res[outer] = append(res[outer], s)
		}
		if len(words) > 1 {
			panic("here")
		}
	}
	return utils.Map(res, utils.Reverse)
}

func myJoin(ss []string) string {
	res := ""
	for _, s := range ss {
		if s != " " {
			res = res + s
		}
	}
	return res
}
