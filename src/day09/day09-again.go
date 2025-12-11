package day09

import (
	"AoC_2025/src/utils"
	"fmt"
	"slices"
)

func Day9_part2try4(input []string) int {
	coords := utils.Map(input, parse)
	// descendingRecatanglesToCheck := rectangleSort(getRectangles(coords))
	// fmt.Println(len(descendingRecatanglesToCheck)) // 28 pairwise checks, good (8*7)/2
	segments := parseSegements(coords)
	// Round 1
	fmt.Println("round 1")
	topS := findTopHorizontal(segments)
	left, right := findTopBars(topS, segments)
	fmt.Println("top", topS) // should be {{7 1} {11 1} -1 1} {{7 3} {7 1} 7 -1} {{11 1} {11 7} 11 -1}
	fmt.Println("left", left)
	fmt.Println("right", right)
	intersecting := findIntersecting(topS, left, right, segments)
	fmt.Println("intersecting", intersecting)
	newsegs, removedInsideChunk := removeChunk(left, topS, right, segments, intersecting)
	fmt.Println(newsegs) // should be map[{{2 3} {11 3} -1 3}:1 {{2 5} {2 3} 2 -1}:1 {{9 5} {2 5} -1 5}:1 {{9 7} {9 5} 9 -1}:1 {{11 3} {11 7} 11 -1}:1 {{11 7} {9 7} -1 7}:1]
	fmt.Println(removedInsideChunk)
	// Round 2
	fmt.Println("round 2")
	topS = findTopHorizontal(newsegs)
	left, right = findTopBars(topS, newsegs)
	fmt.Println("top", topS)
	fmt.Println("left", left)
	fmt.Println("right", right)
	intersecting = findIntersecting(topS, left, right, newsegs)
	fmt.Println("intersecting", intersecting)
	newsegs, removedInsideChunk = removeChunk(left, topS, right, newsegs, intersecting)
	fmt.Println(newsegs) // should eb map[{{9 5} {11 5} -1 5}:1 {{9 7} {9 5} 9 -1}:1 {{11 5} {11 7} 11 -1}:1 {{11 7} {9 7} -1 7}:1]
	fmt.Println(removedInsideChunk)
	// Round 3
	fmt.Println("round 3")
	topS = findTopHorizontal(newsegs)
	left, right = findTopBars(topS, newsegs)
	fmt.Println("top", topS)
	fmt.Println("left", left)
	fmt.Println("right", right)
	intersecting = findIntersecting(topS, left, right, newsegs)
	fmt.Println("intersecting", intersecting)
	newsegs, removedInsideChunk = removeChunk(left, topS, right, newsegs, intersecting)
	fmt.Println(newsegs)
	fmt.Println(removedInsideChunk)
	return 0
}

// one horizontal or vertical border segement
type segment struct {
	start coord
	end   coord
	lineX int // if a vertical line x is always the same, else -1
	lineY int // if a horizontal line this is the y coord, if a vertical line (ie y changes) make -1
}

// Easy mode (intersecting empty):
// This is the case where the chunk can be removed down until one of the left or right bar segements ends.
// Ie, there is no other line segemetn which is "inside" the chunk being removed
//
//	a         _
//	b        | |_  ->    _ _
//	c       _|         _|
//	d
//
// Harder mode (intersecting used)::
// This is the case where the chunk can be removed down only until it hits another line segment contained in the chunk.
// As the intruding chunk can be assumed to be the "outide", only replace the line segement where there was no line previously
// This will also clean up in the case where it hits and collapses its own border
//
//	a      _ _ _                                     _ _ _ _ _
//	b     |  _  |   ->    _   _                     |         |
//	c    _| | | |_      _| | | |_        or         |_ _ _    |    ->    _ _
//	d       | |            | |                            |   |         |   |
//	e
func removeChunk(left, top, right segment, segments map[segment]int, intersecting []segment) (map[segment]int, rectangle) {
	newSegments, rectangleBeingRemoved := createNewSegements(left, top, right, intersecting)

	delete(segments, left)
	delete(segments, top)
	delete(segments, right)
	for _, s := range newSegments {
		segments[s] = 1
	}
	segments = mergeHorizontals(segments)
	return segments, rectangleBeingRemoved
}

func createNewSegements(left, top, right segment, intersecting []segment) ([]segment, rectangle) {
	newSegments := []segment{}
	newTopY := getNewY(left, right, intersecting)
	rx := right.start.x
	lx := left.start.x
	// fmt.Println(newTopY, lx, rx)
	// fmt.Println(top)
	// fmt.Println(left)
	// fmt.Println(right)

	newTopSegments := makeTopSegments(newTopY, lx, rx, intersecting)
	newSegments = slices.Concat(newSegments, newTopSegments)
	rectangle := rectangle{top.start, coord{rx, newTopY}, 0}

	leftRightSegs := makeLeftRigthSegments(left, right, newTopY)
	newSegments = slices.Concat(newSegments, leftRightSegs)

	return newSegments, rectangle
}

func makeTopSegments(newTopY, lx, rx int, intersecting []segment) []segment {
	if len(intersecting) == 0 {
		return []segment{segment{coord{lx, newTopY}, coord{rx, newTopY}, -1, newTopY}}
	} else {
		middleIntersectingCoords := getIntersectingXs(intersecting)
		return makeTops(newTopY, lx, rx, middleIntersectingCoords)
	}
}

func makeTops(newTopY, lx, rx int, intersectXs map[int]bool) []segment {
	segments := []segment{}
	if lx > rx {
		lx, rx = rx, lx // TODO make sure this works
	}
	fmt.Println("lx", lx)
	fmt.Println("rx", rx)
	fmt.Println("newTopY", newTopY)
	s := -1
	for ii := lx; ii <= rx; ii++ {
		_, isIntersecting := intersectXs[ii]
		if s == -1 && !isIntersecting {
			s = ii
		}
		if s != -1 && isIntersecting {
			fmt.Println("hit intersect", ii)
			segments = append(segments, segment{coord{s, newTopY}, coord{ii - 1, newTopY}, -1, newTopY})
			s = -1
		}
	}
	if s != -1 {
		// if -1 at end means never added a segment
		segments = append(segments, segment{coord{lx, newTopY}, coord{rx, newTopY}, -1, newTopY})
	}
	return segments
}

func getIntersectingXs(intersecting []segment) map[int]bool {
	xintersects := map[int]bool{}
	for _, i := range intersecting {
		s := i.start.x
		e := i.end.x
		if i.start.x > i.end.x {
			s = i.end.x
			e = i.start.x
		}
		for ii := s + 1; ii < e; ii++ {
			xintersects[ii] = true
		}
	}
	fmt.Println("Possiblexintersects:", xintersects)
	return xintersects
}

// needs to walk down from starting top y, between l and r, find first y that has intersect or pass min or leeft and rigth downness
func findIntersecting(top, left, right segment, segments map[segment]int) []segment {
	res := []segment{}
	lx := left.start.x
	rx := right.end.x
	if lx > rx {
		lx, rx = rx, lx // TODO make sure this works
	}
	for y := top.start.y; y <= slices.Min([]int{left.start.y, right.end.y}); y++ {
		for k, _ := range segments {
			if k.lineY == y {
				if k != top {
					if k.start.x >= lx && k.start.x <= rx || k.end.x >= lx && k.end.x <= rx {
						res = append(res, k)
					}
				}
			}
		}
	}
	return res
}

func makeLeftRigthSegments(left, right segment, newTopY int) []segment {
	// if l/r bar origianlly descends past the newY, then need to create the section of line starting at newY and going down to whre bar origianlly ended
	segments := []segment{}
	if left.start.y > newTopY {
		segments = append(segments, segment{left.start, coord{left.start.x, newTopY}, left.start.x, -1})
	}
	if right.end.y > newTopY {
		segments = append(segments, segment{coord{right.start.x, newTopY}, right.end, right.start.x, -1})
	}
	return segments
}

func getNewY(l segment, r segment, intersecting []segment) int {
	if len(intersecting) > 0 {
		return intersecting[0].lineY
	}
	if l.start.y > r.end.y {
		return r.end.y
	}
	return l.start.y
}

func mergeHorizontals(segs map[segment]int) map[segment]int {
	for k1, _ := range segs {
		for k2, _ := range segs {
			if k1 != k2 && k1.lineY != -1 && k2.lineY != -1 && k1.lineY == k2.lineY {
				// not same and both horizontal
				delete(segs, k1)
				delete(segs, k2)
				if k1.start == k2.end {
					// merge
					newS := segment{k2.start, k1.end, -1, k1.lineY}
					segs[newS] = 1
				}
				if k1.end == k2.start {
					newS := segment{k1.start, k2.end, -1, k2.lineY}
					segs[newS] = 1
				}
			}
		}
	}
	return segs
}

func parseSegements(cs []coord) map[segment]int {
	segments := map[segment]int{}
	for i := 0; i < len(cs); i++ {
		j := i + 1
		if i == len(cs)-1 {
			j = 0
		}
		newS := segment{cs[i], cs[j], -1, -1}
		if cs[i].y == cs[j].y {
			newS.lineY = cs[i].y
		}
		if cs[i].x == cs[j].x {
			newS.lineX = cs[i].x
		}
		segments[newS] = 1
	}
	return segments
}

func findTopHorizontal(segments map[segment]int) segment {
	highest := segment{}
	highest.lineY = -1
	for k, _ := range segments {
		if k.lineY != -1 {
			if highest.lineY == -1 {
				highest = k
			}
			if k.lineY < highest.lineY {
				highest = k
			}
		}
	}
	return highest
}

// ebc directionality might become a problem. ie might have two bars that connect via their ends. but code above needs left and right with start and end as expected
// if k.start = top.start for exmaple, need to reverse k before set as left. prints for now
func findTopBars(top segment, segments map[segment]int) (segment, segment) {
	var left segment
	var right segment
	for k, _ := range segments {
		if k.lineY == -1 {
			if k.end == top.start {
				left = k
			}
			if k.start == top.start {
				fmt.Println("need to implement 1") // Todo almost 100%
			}
			if k.start == top.end {
				right = k
			}
			if k.end == top.end {
				fmt.Println("need to implement 2")
			}
		}
	}
	return left, right
}
