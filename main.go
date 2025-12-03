package main

import (
	day1 "AoC_2025/src/day01"
	day2 "AoC_2025/src/day02"
	day3 "AoC_2025/src/day03"
	"AoC_2025/src/utils"
	"fmt"
)

var mode utils.Mode = utils.Both // Can be Real, Test, or Both

func main() {
	fmt.Println("Starting tests")

	fmt.Println("\n### Day 1 ###")
	utils.Run(day1.Day1_part1, "01", "1", mode, 1139, 3)
	utils.Run(day1.Day1_part2, "01", "2", mode, 6684, 6)

	fmt.Println("\n### Day 2 ###")
	utils.Run(day2.Day2_part1, "02", "1", mode, 8576933996, 1227775554)
	utils.Run(day2.Day2_part2, "02", "2", mode, 25663320831, 4174379265)

	fmt.Println("\n### Day 3 ###")
	utils.Run(day3.Day3_part1, "03", "1", mode, 17430, 357)
	utils.Run(day3.Day3_part2, "03", "2", mode, 171975854269367, 3121910778619)

	utils.Check_all()
}
