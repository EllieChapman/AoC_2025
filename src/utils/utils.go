package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

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

// func Check_answer[A comparable](ans A, expect A) bool {
// 	if ans == expect {
// 		fmt.Println("Test Passed")
// 		return true
// 	} else {
// 		fmt.Println("Test Failed, expected:", expect, "received:", ans)
// 		suite_passed = false
// 		return false
// 	}
// }

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

func Map[A, B any](as []A, f func(A) B) []B {
	result := make([]B, len(as))
	for i, a := range as {
		result[i] = f(a)
	}
	return result
}

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

/// Dont need as no error handling, use directly
// func Stringify(i int) string {
// 	s := strconv.Itoa(i)
// 	return s
// }

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
