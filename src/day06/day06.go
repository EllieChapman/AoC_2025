package day06

import (
	"AoC_2025/src/utils"
	"strings"
)

func Day6_part1(input []string) int {
	sums := utils.Map(utils.Transpose(utils.Map(input, strings.Fields)), func(ss []string) sum { return sum{ss[len(ss)-1], utils.Map(ss[0:len(ss)-1], utils.Atoi)} })
	return utils.Sum(utils.Map(sums, doSum))
}

func Day6_part2(input []string) int {
	transposed := utils.Transpose(utils.Map(input[0:len(input)-1], func(s string) []string { return strings.Split(s, "") }))
	joined := utils.Map(transposed, myJoin)
	ops := strings.Fields(input[len(input)-1])
	sums := parseRec(joined, ops, []sum{})
	return utils.Sum(utils.Map(sums, doSum))
}

type sum struct {
	op   string
	nums []int
}

func doSum(s sum) int {
	switch s.op {
	case "+":
		return utils.Sum(s.nums)
	case "*":
		return utils.Mul(s.nums)
	default:
		panic("unexpected operator:")
	}
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

func parseRec(ss []string, ops []string, sums []sum) []sum {
	newSum := sum{ops[0], []int{}}
	for pos, s := range ss {
		if s == "" {
			sums = append(sums, newSum)
			return parseRec(ss[pos+1:], ops[1:], sums)
		} else {
			newSum.nums = append(newSum.nums, utils.Atoi(s))
		}
	}
	return append(sums, newSum)
}
