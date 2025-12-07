package day07

import (
	"AoC_2025/src/utils"
	"strings"
)

func Day7_part1(input []string) int {
	startBeamsPos := getStart(input[0])
	_, numSplit := applyAll(startBeamsPos, utils.Map(input[1:], func(s string) []string { return strings.Split(s, "") }))
	return numSplit
}

func Day7_part2(input []string) int {
	startBeamsPos := getStart(input[0])
	finalBeamsPos, _ := applyAll(startBeamsPos, utils.Map(input[1:], func(s string) []string { return strings.Split(s, "") }))
	totalBeams := sumMap(finalBeamsPos)
	return totalBeams
}

func apply(currentBeamsPos map[int]int, row []string) (map[int]int, int) {
	numSplit := 0
	newBeamsPos := map[int]int{}
	for beamPosKey, numBeamsInPos := range currentBeamsPos {
		if row[beamPosKey] == "^" {
			numSplit++
			if _, ok := newBeamsPos[beamPosKey-1]; ok {
				newBeamsPos[beamPosKey-1] = newBeamsPos[beamPosKey-1] + numBeamsInPos
			} else {
				newBeamsPos[beamPosKey-1] = numBeamsInPos
			}
			if _, ok := newBeamsPos[beamPosKey+1]; ok {
				newBeamsPos[beamPosKey+1] = newBeamsPos[beamPosKey+1] + numBeamsInPos
			} else {
				newBeamsPos[beamPosKey+1] = numBeamsInPos
			}
		} else {
			if _, ok := newBeamsPos[beamPosKey]; ok {
				newBeamsPos[beamPosKey] = newBeamsPos[beamPosKey] + numBeamsInPos
			} else {
				newBeamsPos[beamPosKey] = numBeamsInPos
			}
		}
	}
	return newBeamsPos, numSplit
}

func applyAll(currentBeams map[int]int, rows [][]string) (map[int]int, int) {
	totalSplit := 0
	var numSplit int
	for _, row := range rows {
		currentBeams, numSplit = apply(currentBeams, row)
		totalSplit += numSplit
	}
	return currentBeams, totalSplit
}

func getStart(s string) map[int]int {
	ss := strings.Split(s, "")
	for pos, value := range ss {
		if value == "S" {
			return map[int]int{pos: 1}
		}
	}
	panic("could not find start")
}

func sumMap(m map[int]int) int {
	total := 0
	for _, value := range m {
		total += value
	}
	return total
}
