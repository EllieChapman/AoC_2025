package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Mode int // Custom type for our enum

const (
	Real Mode = iota // 0
	Test             // 1
	Both             // 2
)

func Run(f func([]string) int, day string, part string, mode Mode, realExp int, testExp int) {
	switch mode {
	case Real:
		RunReal(f, day, part, realExp)
	case Test:
		RunTest(f, day, part, testExp)
	case Both:
		RunTest(f, day, part, testExp)
		RunReal(f, day, part, realExp)
	default:
		// Handle unexpected states or provide a default action
		fmt.Println("unexpected mode", mode)
	}
}

func RunReal(f func([]string) int, day string, part string, realExp int) {
	realInput := GetInput(day, Real)
	realRes := f(realInput)
	fmt.Println("Day", day, "Part", part, "(Mode Real):", Check_answer(realRes, realExp))
}

func RunTest(f func([]string) int, day string, part string, testExp int) {
	testInput := GetInput(day, Test)
	testRes := f(testInput)
	fmt.Println("Day", day, "Part", part, "(Mode Test):", Check_answer(testRes, testExp))
}

func GetInput(day string, mode Mode) []string {
	// construct file name
	filePath := "src/day" + day + "/input.txt"
	if mode != Real {
		filePath = "src/day" + day + "/test-input.txt"
	}
	// call through to ReadLines
	return ReadLines(filePath)
}

// read line by line into memory
// all file contents is stores in lines[]
func ReadLines(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file", err)
		return nil
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func Check_answer[A comparable](ans A, expect A) string {
	if ans == expect {
		return "Test Passed"
	} else {
		suite_passed = false
		s := fmt.Sprintf("Test Failed, expected: %v, received: %v", expect, ans)
		return s
	}
}

func Check_all() bool {
	if suite_passed {
		println("\nAll tests passed!")
		return true
	} else {
		println("\nSome tests failed :(")
		return false
	}
}

var suite_passed = true

func Flatten[A any](as [][]A) []A {
	result := make([]A, 0)
	for _, a := range as {
		result = append(result, a...)
	}
	return result
}

func Atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic("Tried to convert string '" + s + "' to int")
	}
	return i
}

func Sum(s []int) int {
	total := 0
	for _, i := range s {
		total += i
	}
	return total
}

// func Reverse[A any](s []A) []A {
// 	slices.Reverse(s)
// 	return s
// }

// func AbsDiffInt(x, y int) int {
// 	if x < y {
// 		return y - x
// 	}
// 	return x - y
// }

func Map[A, B any](as []A, f func(A) B) []B {
	result := make([]B, len(as))
	for i, a := range as {
		result[i] = f(a)
	}
	return result
}

// func Map2[A, B, C, D any](as []A, f func(A, C, D) B, c C, d D) []B {
// 	result := make([]B, len(as))
// 	for i, a := range as {
// 		result[i] = f(a, c, d)
// 	}
// 	return result
// }

// func Map1[A, B, C any](as []A, f func(A, C) B, c_extraArgs C) []B {
// 	result := make([]B, len(as))
// 	for i, a := range as {
// 		result[i] = f(a, c_extraArgs)
// 	}
// 	return result
// }

func MapVariadic[A, B, C any](as []A, f func(A, ...C) B, c_extraArgs ...C) []B {
	result := make([]B, len(as))
	for i, a := range as {
		result[i] = f(a, c_extraArgs...)
	}
	return result
}
