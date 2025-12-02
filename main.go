package main

import (
	day1 "AoC_2025/src/day01"
	day2 "AoC_2025/src/day02"
	"AoC_2025/src/utils"
	"fmt"
)

func main() {
	fmt.Println("Starting tests")

	fmt.Println("Day 1 Part 1:", utils.Check_answer(day1.Day1_part1(), 1139))
	fmt.Println("Day 1 Part 2:", utils.Check_answer(day1.Day1_part2(), 6684))

	// EBC need to be able to easily swap between running test and real modes, and save off both data sets
	fmt.Println("Day 2 Part 1:", utils.Check_answer(day2.Day2_part1(), 8576933996))  // 1227775554 test, 8576933996 real
	fmt.Println("Day 2 Part 2:", utils.Check_answer(day2.Day2_part2(), 25663320831)) // 4174379265 test, 25663320831 real

	utils.Check_all()
}
