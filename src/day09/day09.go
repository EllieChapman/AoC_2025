package day09

import (
	"AoC_2025/src/utils"
	"slices"
	"sort"
	"strings"
)

func Day9_part1(input []string) int {
	coords := utils.Map(input, parse)
	sizesDescending := rectangleSort(getRectangles(coords))
	return sizesDescending[0].area
}

func Day9_part2(input []string) int {
	coords := utils.Map(input, parse)
	sizesDescending := rectangleSort(getRectangles(coords))
	// fmt.Println(len(sizesDescending)) // 28 pairwise checks, good (8*7)/2
	redGreenBorderCoords := getStartingRedGreenCoords(coords)
	// fmt.Println(len(redGreenBorderCoords)) // 30, good
	r := getLargestPossibleRectangle(sizesDescending, redGreenBorderCoords)
	return r.area
}

// func Day9_part2try3(input []string) int {
// 	coords := utils.Map(input, parse)
// 	sizesDescending := rectangleSort(getRectangles(coords))
// 	// fmt.Println(len(sizesDescending)) // 28 pairwise checks, good (8*7)/2
// 	redGreenBorderCoords := getStartingRedGreenCoords(coords)
// 	// fmt.Println(len(redGreenBorderCoords)) // 30, good
// 	r := getLargestPossibleRectangle3(sizesDescending, redGreenBorderCoords)
// 	return r.area
// }

// func Day9_part2try2(input []string) int {
// 	coords := utils.Map(input, parse)
// 	sizesDescending := rectangleSort(getRectangles(coords))
// 	fmt.Println(len(sizesDescending)) // 28 pairwise checks, good (8*7)/2
// 	redGreenBorderCoords := getStartingRedGreenCoords(coords)
// 	fmt.Println(len(redGreenBorderCoords)) // 30, good
// 	// create map of just outside coords
// 	startCoord := coord{} // todo
// 	outlineBorderCoords := createOutlineBorderCoords(startCoord, map[coord]int{}, redGreenBorderCoords)
// 	r := getLargestPossibleRectangle2(sizesDescending, redGreenBorderCoords, outlineBorderCoords)
// 	return r.area
// }

type coord struct {
	x int
	y int
}

type rectangle struct {
	a    coord
	b    coord
	area int
}

// func getLargestPossibleRectangle3(rs []rectangle, redGreenBorderCoords map[coord]int) rectangle {
// 	for _, r := range rs {
// 		// once fix updating red green, also fix outside updating rather than starting empty for eahc rectangle
// 		if isRectanglePossible3(r, redGreenBorderCoords) {
// 			return r
// 		}
// 	}
// 	panic("none possible")
// }

// func isRectanglePossible3(r rectangle, redGreenBorderCoords map[coord]int) bool {
// 	c1 := coord{r.a.x, r.a.y}
// 	c2 := coord{r.a.x, r.b.y}
// 	c3 := coord{r.b.x, r.b.y}
// 	c4 := coord{r.b.x, r.a.y}
// 	greenLineCoords := slices.Concat(getLineCoords(c1, c2), getLineCoords(c2, c3), getLineCoords(c3, c4), getLineCoords(c4, c1))
// 	for _, c := range greenLineCoords {
// 		_, partOfRedGreenBorder := redGreenBorderCoords[c]
// 		if !partOfRedGreenBorder {
// 			// if not part of border, check if is outside, rectangle fails if a single coord is
// 			if isOuter3(c, redGreenBorderCoords) {
// 				return false
// 			}
// 		}
// 	}

// 	// move up until hit redGreenBorderCoords. Space under is startCoord

// 	// recurse and try
// 	// try move up
// 	// if not, try move right
// 	// if not, try move down
// 	// if not, try move left

// 	// never move onto the redGreenBorderCoords
// 	// never move where you have already been - except for on to startCoord
// 	// if cant' move need to back track in recursive call

// 	// eventually will either:
// 	// reach the edge x/y <0 or >1000000
// 	// reach where you started
// }

// =====================================================

// func createOutlineBorderCoords(currentlyCheckingOutline coord, knownOutline map[coord]int, redGreenBorderCoords map[coord]int) map[coord]int {

// }

// func getLargestPossibleRectangle2(rs []rectangle, redGreenBorderCoords map[coord]int, outlineBorderCoords map[coord]int) rectangle {
// 	for _, r := range rs {
// 		if isRectanglePossible2(r, redGreenBorderCoords, outlineBorderCoords) {
// 			return r
// 		}
// 	}
// 	panic("none possible")
// }

// func isRectanglePossible2(r rectangle, redGreenBorderCoords map[coord]int, outlineBorderCoords map[coord]int) bool {
// 	c1 := coord{r.a.x, r.a.y}
// 	c2 := coord{r.a.x, r.b.y}
// 	c3 := coord{r.b.x, r.b.y}
// 	c4 := coord{r.b.x, r.a.y}
// 	greenLineCoords := slices.Concat(getLineCoords(c1, c2), getLineCoords(c2, c3), getLineCoords(c3, c4), getLineCoords(c4, c1))
// 	for _, c := range greenLineCoords {
// 		_, partOfOutlineBorderCoords := outlineBorderCoords[c]
// 		if partOfOutlineBorderCoords {
// 			return false
// 		}
// 		_, partOfRedGreenBorder := redGreenBorderCoords[c]
// 		if !partOfRedGreenBorder {
// 			// not in iner or outer border
// 			if isOuter(c, redGreenBorderCoords, outlineBorderCoords) {
// 				return false
// 			}
// 		}
// 	}
// 	return true
// }

// func isOuter(c coord, redGreenBorderCoords map[coord]int, outlineBorderCoords map[coord]int) bool {
// 	// keep chekcing one up (x = x, y = y-1)
// 	// if hit outer or  y = 0, return true
// 	// if hit inner border, rewturn false
// 	for yy := c.y; yy > 0; yy-- {
// 		newC := coord{c.x, yy}
// 		_, hitOuter := outlineBorderCoords[newC]
// 		_, hitInner := redGreenBorderCoords[newC]
// 		if hitOuter {
// 			return true
// 		}
// 		if hitInner {
// 			return false
// 		}
// 	}
// 	return true // hit outer bounds, y=0
// }

// ==============================================================

func getLargestPossibleRectangle(rs []rectangle, redGreenBorderCoords map[coord]int) rectangle {
	outsideCoords := map[coord]int{}
	var isPossible bool
	for _, r := range rs {
		// once fix updating red green, also fix outside updating rather than starting empty for eahc rectangle
		isPossible, redGreenBorderCoords, outsideCoords = isRectanglePossible(r, redGreenBorderCoords, outsideCoords)
		if isPossible {
			return r
		}
	}
	panic("none possible")
}

func isRectanglePossible(r rectangle, redGreenBorderCoords map[coord]int, outsideCoords map[coord]int) (bool, map[coord]int, map[coord]int) {
	c1 := coord{r.a.x, r.a.y}
	c2 := coord{r.a.x, r.b.y}
	c3 := coord{r.b.x, r.b.y}
	c4 := coord{r.b.x, r.a.y}
	greenLineCoords := slices.Concat(getLineCoords(c1, c2), getLineCoords(c2, c3), getLineCoords(c3, c4), getLineCoords(c4, c1))
	var isInside bool
	for _, c := range greenLineCoords {
		_, partOfBorder := redGreenBorderCoords[c]
		isInside, redGreenBorderCoords, outsideCoords = isCoordInside(redGreenBorderCoords, outsideCoords, map[coord]int{}, map[coord]int{c: 1})
		if !partOfBorder && !isInside {
			return false, redGreenBorderCoords, outsideCoords
		}
	}
	// no coord is outside the overall border
	return true, redGreenBorderCoords, outsideCoords
}

func isCoordInside(redGreenBorderCoords map[coord]int, outsideCoords map[coord]int, allVisited map[coord]int, currentFront map[coord]int) (bool, map[coord]int, map[coord]int) {
	nextFront := map[coord]int{}
	for k, _ := range currentFront {
		gen8 := generateAdjacent(k)
		for _, c := range gen8 {
			_, isOutside := outsideCoords[c]
			if isOutside || isOffGrid(c) {
				return false, redGreenBorderCoords, outsideCoords // outside coord
			}
			_, visited := allVisited[c]
			_, isBorder := redGreenBorderCoords[c]
			if !visited && !isBorder {
				nextFront[c] = 1
			}
		}
	}
	if len(nextFront) == 0 {
		// fmt.Println("inside len", len(redGreenBorderCoords))
		// must be an inside coord, checked everywhere can reach without crossing border and never went off grid
		// means everything in allVisited, currentFront and nextFront(but is empty) is also an inside coord, can add to redGreenBorderCoords
		for k, _ := range allVisited {
			redGreenBorderCoords[k] = 1
		}
		for k, _ := range currentFront {
			redGreenBorderCoords[k] = 1
		}
		return true, redGreenBorderCoords, outsideCoords
	}
	for k, _ := range currentFront {
		allVisited[k] = 1
	}
	return isCoordInside(redGreenBorderCoords, outsideCoords, allVisited, nextFront)
}

func isOffGrid(c coord) bool {
	return c.x == -1 || c.y == -1
}

func generateAdjacent(c coord) []coord {
	x, y := c.x, c.y
	return []coord{
		{x + 1, y + 1},
		{x + 1, y},
		{x + 1, y - 1},
		{x, y + 1},
		{x - 1, y + 1},
		{x - 1, y},
		{x - 1, y - 1},
		{x, y - 1},
	}
}

func getStartingRedGreenCoords(redCs []coord) map[coord]int {
	m := map[coord]int{}
	for i := 0; i < len(redCs); i++ {
		m[redCs[i]] = 0
		j := i + 1
		if i == len(redCs)-1 {
			j = 0
		}
		greenLine := getLineCoords(redCs[i], redCs[j])
		for _, greenC := range greenLine {
			m[greenC] = 0
		}
	}
	return m
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

// start same, get list of red coords RedCoords
// create RedGreenCoords union/map of all red/green coords, start by putting in all red and work out green coord lines between red coords
// create UncolouredCoods, will add coords that are not red or green, starts empty. Might nto actually need
// create Unassigned, calc all coords and add those not yet in RedGreenCoords (or UncolouredCoods)
// for coord in Uassigned, work out if should be in RedGreenCoords or UncolouredCoods. Needs a coord in RedGreenCoords above, below, left and right of it.
// !!!! test nto so simple, need to take a point and try to walk from it to off edge without crossing the line made up of RedGreenCoords, dikjstras search
// (nb might be easier not to immediately add from Uassigned to RedGreenCoords, to avoid bloat when looking for coord u/d/l/r ect)

// Should end up with
// RedCoords - starting list of red coords
// RedGreenCoords - all coords which can be in a rectangle
// UncolouredCoods - coords whihc cannot be in a rectangle
// (Unassigned - empty)

// Basically repeat part 1, but instead get a sorted list of {a coord, b coord, size} instead of just the size
// Then work through, generate all coords in rectangle for a,b red coord corners, if any are in UncolouredCoods bad
// Return the size for first rectangle which has no bad coords

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
