package main

import (
	day1 "AoC_2025/src/day01"
	"AoC_2025/src/utils"
	"fmt"
)

func main() {
	fmt.Println("Starting tests")

	fmt.Println("Day 1 Part 1")
	utils.Check_answer(day1.Day1_part1(), 1139)
	fmt.Println("Day 1 Part 2")
	utils.Check_answer(day1.Day1_part2(), 6684)

	utils.Check_all()
}
