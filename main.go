package main

import (
	"AoC_2025/src/day01"
	"AoC_2025/src/day02"
	"AoC_2025/src/day03"
	"AoC_2025/src/day04"
	"AoC_2025/src/day05"
	"AoC_2025/src/day06"
	"AoC_2025/src/day07"
	"AoC_2025/src/utils"
	"fmt"
)

var mode utils.Mode = utils.Both // Can be Real, Test, or Both

func main() {
	fmt.Println("Starting tests")

	fmt.Println("\n### Day 1 ###")
	utils.Run(day01.Day1_part1, "01", "1", mode, 1139, 3)
	utils.Run(day01.Day1_part2, "01", "2", mode, 6684, 6)

	fmt.Println("\n### Day 2 ###")
	utils.Run(day02.Day2_part1, "02", "1", mode, 8576933996, 1227775554)
	utils.Run(day02.Day2_part2, "02", "2", mode, 25663320831, 4174379265)

	fmt.Println("\n### Day 3 ###")
	utils.Run(day03.Day3_part1, "03", "1", mode, 17430, 357)
	utils.Run(day03.Day3_part2, "03", "2", mode, 171975854269367, 3121910778619)

	fmt.Println("\n### Day 4 ###")
	utils.Run(day04.Day4_part1, "04", "1", mode, 1467, 13)
	utils.Run(day04.Day4_part2, "04", "2", mode, 8484, 43)

	fmt.Println("\n### Day 5 ###")
	utils.Run(day05.Day5_part1, "05", "1", mode, 811, 3)
	utils.Run(day05.Day5_part2, "05", "2", mode, 338189277144473, 14)

	fmt.Println("\n### Day 6 ###")
	utils.Run(day06.Day6_part1, "06", "1", mode, 4405895212738, 4277556)
	utils.Run(day06.Day6_part2, "06", "2", mode, 7450962489289, 3263827)

	fmt.Println("\n### Day 7 ###")
	utils.Run(day07.Day7_part1, "07", "1", mode, 1609, 21)
	utils.Run(day07.Day7_part2, "07", "2", mode, 12472142047197, 40)

	utils.Check_all()
}
