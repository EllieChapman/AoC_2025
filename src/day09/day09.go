package day09

import (
	"AoC_2025/src/utils"
	"sort"
	"strings"
)

func Day9_part1(input []string) int {
	coords := utils.Map(input, parse)
	sizesDescending := rectangleSort(getRectangles(coords))
	return sizesDescending[0].area
}

type coord struct {
	x int
	y int
}

type rectangle struct {
	a    coord
	b    coord
	area int
}

func parse(s string) coord {
	ss := utils.Map(strings.Split(s, ","), utils.Atoi)
	return coord{ss[0], ss[1]}
}

func calcRectangleSize(a coord, b coord) rectangle {
	xlen := 1 + utils.AbsDiffInt(a.x, b.x)
	ylen := 1 + utils.AbsDiffInt(a.y, b.y)
	return rectangle{a, b, xlen * ylen}
}

func getRectangles(cs []coord) []rectangle {
	sizes := []rectangle{}
	for i := 0; i < len(cs); i++ {
		for j := i + 1; j < len(cs); j++ {
			sizes = append(sizes, calcRectangleSize(cs[i], cs[j]))
		}
	}
	return sizes
}

func rectangleSort(rs []rectangle) []rectangle {
	sort.Slice(rs, func(i, j int) bool {
		return rs[i].area > rs[j].area
	})
	return rs
}
