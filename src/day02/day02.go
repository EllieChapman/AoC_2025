package day2

import (
	"AoC_2025/src/utils"
	"strconv"
	"strings"
)

func Day2_part1(input []string) int {
	ranges := utils.Map(strings.Split(input[0], ","), parseRangeId)
	return utils.Sum(utils.MapVariadic(ranges, checkRange))
}

func Day2_part2(input []string) int {
	ranges := utils.Map(strings.Split(input[0], ","), parseRangeId)
	return utils.Sum(utils.Map(ranges, checkRangeComplicated))
}

type rangeID struct {
	start int
	end   int
}

func parseRangeId(s string) rangeID {
	ss := strings.Split(s, "-")
	return rangeID{
		start: utils.Atoi(ss[0]),
		end:   utils.Atoi(ss[1]),
	}
}

func checkRange(r rangeID, _ ...any) int {
	total := 0
	for i := r.start; i <= r.end; i++ {
		total = total + checkID(i)
	}
	return total
}

// EBC this is ugly converting back to a string. Can I use int directly?
func checkID(i int) int {
	s := strconv.Itoa(i)
	l := len(s)
	if l%2 == 0 {
		if s[0:(l/2)] == s[(l/2):] {
			return i
		}
	}
	// return ID if invalid, else 0
	return 0
}

func checkRangeComplicated(r rangeID) int {
	total := 0
	for i := r.start; i <= r.end; i++ {
		total = total + checkIDComplicated(i)
	}
	return total
}

func checkIDComplicated(i int) int {
	s := strconv.Itoa(i)

	// Test for repearting chunks of increasing size, starting with single characters, up to half the original string
	for chunkSize := 1; chunkSize <= len(s)/2; chunkSize++ {
		if chunkSizeRepeats(s, chunkSize) {
			return i
		}
	}

	return 0
}

func chunkSizeRepeats(s string, chunkSize int) bool {
	if len(s)%chunkSize != 0 {
		return false
	}

	chunkOne := s[0:chunkSize]
	for ii := 0; ii < len(s); ii = ii + chunkSize {
		if s[ii:ii+chunkSize] != chunkOne {
			return false
		}
	}
	return true
}
