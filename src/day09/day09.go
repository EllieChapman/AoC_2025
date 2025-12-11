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

func getLineCoords(a coord, b coord) []coord {
	if a.x != b.x && a.y != b.y {
		panic("line not horizontal or vertical")
	}
	if a.x == b.x {
		if a.y < b.y {
			return makeVertical(a.x, a.y, b.y)
		} else {
			return makeVertical(a.x, b.y, a.y)
		}
	} else {
		if a.x < b.x {
			return makeHorizontal(a.x, b.x, a.y)
		} else {
			return makeHorizontal(b.x, a.x, a.y)
		}
	}

}

func makeHorizontal(x1, x2, y int) []coord {
	cs := []coord{}
	for x := x1 + 1; x < x2; x++ {
		cs = append(cs, coord{x, y})
	}
	return cs
}

func makeVertical(x, y1, y2 int) []coord {
	cs := []coord{}
	for y := y1 + 1; y < y2; y++ {
		cs = append(cs, coord{x, y})
	}
	return cs
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
