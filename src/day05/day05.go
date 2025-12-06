package day05

import (
	"AoC_2025/src/utils"
	"slices"
	"strings"
)

func Day5_part1(input []string) int {
	ranges, ids := parse(input)
	// utils.MapV(ids, checkId, ranges) // Why didn't it like [rangeID] as C?
	return utils.Sum(utils.Map(ids, func(i int) int { return checkId(i, ranges) }))
}

func Day5_part2(input []string) int {
	ranges, _ := parse(input)
	squashed := squash(ranges)
	return utils.Sum(utils.Map(squashed, func(r rangeID) int { return 1 + r.end - r.start }))
}

type rangeID struct {
	start int
	end   int
}

func parse(lines []string) ([]rangeID, []int) {
	for pos, line := range lines {
		if len(line) == 0 {
			ranges := utils.Map(lines[0:pos], parseRangeId)
			ids := utils.Map(lines[pos+1:], utils.Atoi)
			return ranges, ids
		}
	}
	panic("no blank line found")
}

func parseRangeId(line string) rangeID {
	is := utils.Map(strings.Split(line, "-"), utils.Atoi)
	return rangeID{is[0], is[1]}
}

// return 1 if fresh
func checkId(id int, ranges []rangeID) int {
	// ranges := a[0].([]rangeID)
	for _, r := range ranges {
		if id >= r.start && id <= r.end {
			return 1
		}
	}
	return 0
}

func squash(ranges []rangeID) []rangeID {
	// pairwise try to combine, if can do and recurse. if no more can be combined we are done
	for x_pos, xr := range ranges {
		for y_pos, yr := range ranges { // NB wll check mroe than need, can clean up. Can make even more sure that ypos is always buigger if ypos start iterating from xpos+1
			if x_pos != y_pos {
				b, newR := combine(xr, yr)
				if b {
					newRs := slices.Delete(slices.Delete(ranges, y_pos, y_pos+1), x_pos, x_pos+1)
					newRs = append(newRs, newR)
					return squash(newRs)
				}
			}
		}
	}
	// no more can be combined
	return ranges
}

func combine(a rangeID, b rangeID) (bool, rangeID) {
	// if can combine do and return combined
	if !(a.end < b.start || b.end < a.start) {
		return true, rangeID{slices.Min([]int{a.start, b.start}), slices.Max([]int{a.end, b.end})}
	}
	return false, rangeID{}
}
