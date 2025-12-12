package day09

import (
	"AoC_2025/src/utils"
	"slices"
	"sort"
)

func Day9_part2(input []string) int {
	coords := utils.Map(input, parse)
	descendingRecatanglesToCheck := rectangleSort(getRectangles(coords))
	segments := parseSegements(coords)
	safeChunks := loop(segments, []rectangle{})
	r := getLargestPossibleRectangle3(descendingRecatanglesToCheck, safeChunks)
	return r.area
}

func loop(segs map[segment]int, accRec []rectangle) []rectangle {
	topS := findTopHorizontal(segs)
	left, right := findTopBars(topS, segs)
	newsegs, removedInsideChunk := removeChunk(left, topS, right, segs)
	accRec = append(accRec, removedInsideChunk)
	if len(newsegs) == 0 {
		return accRec
	}
	return loop(newsegs, accRec)
}

// one horizontal or vertical border segement
type segment struct {
	start coord
	end   coord
	lineX int // if a vertical line x is always the same, else -1
	lineY int // if a horizontal line this is the y coord, if a vertical line (ie y changes) make -1
}

func getLargestPossibleRectangle3(toCheckrs []rectangle, safeRS []rectangle) rectangle {
	for _, r := range toCheckrs[0:] {
		if isRectanglePossibleFast(r, safeRS) {
			return r
		}
	}
	panic("none possible")
}

// needs to check segments, not coords
func isRectanglePossibleFast(r rectangle, safeRS []rectangle) bool {
	c1 := coord{r.a.x, r.a.y}
	c2 := coord{r.a.x, r.b.y}
	c3 := coord{r.b.x, r.b.y}
	c4 := coord{r.b.x, r.a.y}
	segsToCheck := []segment{makeSegment(c1, c2), makeSegment(c2, c3), makeSegment(c3, c4), makeSegment(c4, c1)}
	for _, seg := range segsToCheck {
		if !isSegInside([]segment{seg}, safeRS) {
			return false
		}
	}
	return true
}

func isSegInside(fragmentsToCheck []segment, safeRS []rectangle) bool {
	for _, r := range safeRS {
		newFrags := []segment{}
		for _, f := range fragmentsToCheck {
			new := r.containsSeg(f)
			newFrags = slices.Concat(newFrags, new)
		}
		fragmentsToCheck = newFrags
		if len(fragmentsToCheck) == 0 {
			return true
		}
	}
	return false // run out of safe rectabgles but still have somefragements to fit
}

// return list of 0, 1 or 2 segements, which are fragments of the orgianl segment that were not contained
// NB needs seg to be ordered
func (r rectangle) containsSeg(s segment) []segment {
	y1 := r.a.y
	y2 := r.b.y
	x1 := r.a.x
	x2 := r.b.x
	if s.isHorizontal() {
		// y is either in or out
		test_y := s.getHorizontalLinePos()
		if !(test_y >= y1 && test_y <= y2) {
			return []segment{s}
		}

		// only care about the x values in the rectagle for making chunks
		oneDranges := compareRanges(s.start.x, s.end.x, x1, x2)
		segs := []segment{}
		for _, o := range oneDranges {
			segs = append(segs, makeSegment(coord{o.a, test_y}, coord{o.b, test_y}))
		}
		return segs
	} else {
		// x is either in or out
		test_x := s.getVerticalLinePos()
		if !(test_x >= x1 && test_x <= x2) {
			return []segment{s}
		}

		// only care about the y values in the rectagle for making chunks
		oneDranges := compareRanges(s.start.y, s.end.y, y1, y2)
		segs := []segment{}
		for _, o := range oneDranges {
			segs = append(segs, makeSegment(coord{test_x, o.a}, coord{test_x, o.b}))
		}
		return segs
	}
}

type oneDRange struct {
	a int
	b int
}

// return segs not included in parant
func compareRanges(test1, test2, parent1, parent2 int) []oneDRange {
	if test2 < parent1 || test1 > parent2 {
		// non included, so fragemnt is entire orginal segment
		return []oneDRange{{test1, test2}}
		// return false, []oneDRange{} // change for 2 mthdos
	}
	if test1 >= parent1 && test2 <= parent2 {
		// all included, so no fragmenst to return
		return []oneDRange{}
	}
	if test1 >= parent1 && test1 <= parent2 && test2 > parent2 {
		// will generate one fragement, beyond parent 2
		return []oneDRange{{parent2 + 1, test2}}
	}
	if test2 >= parent1 && test2 <= parent2 && test1 < parent1 {
		// will generate one fragement, before parent 1
		return []oneDRange{{test1, parent1 - 1}}
	}
	if test1 < parent1 && test2 > parent2 {
		// two frgament, on eitehr sid eof th eparents
		return []oneDRange{{parent2 + 1, test2}, {test1, parent1 - 1}}
	}
	panic("another case didnt consider")
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
func removeChunk(left, top, right segment, segments map[segment]int) (map[segment]int, rectangle) {
	intersecting, newTopY := findIntersecting(top, left, right, segments)
	newSegments, rectangleBeingRemoved := createNewSegements(left, top, right, newTopY, intersecting)

	delete(segments, left)
	delete(segments, top)
	delete(segments, right)

	for _, i := range intersecting {
		delete(segments, i)
	}
	for _, s := range newSegments {
		segments[s] = 1
	}
	segments = mergeHorizontals(segments)
	return segments, rectangleBeingRemoved
}

func createNewSegements(left, top, right segment, newTopY int, intersecting []segment) ([]segment, rectangle) {
	newSegments := []segment{}
	rx := right.getVerticalLinePos()
	lx := left.getVerticalLinePos()

	newTopSegments := findNextTopFast(lx, rx, newTopY, intersectSort(intersecting), []segment{})

	newSegments = slices.Concat(newSegments, newTopSegments)
	rectangle := calcRectangleSize(top.start, coord{rx, newTopY})

	leftRightSegs := makeLeftRigthSegments(left, right, newTopY)
	newSegments = slices.Concat(newSegments, leftRightSegs)

	return newSegments, rectangle
}

func intersectSort(ls []segment) []segment {
	sort.Slice(ls, func(i, j int) bool {
		return ls[i].start.x < ls[j].end.x
	})
	return ls
}

func findNextTopFast(lx, rx, ytop int, orderedIntersects []segment, accTops []segment) []segment {
	if len(orderedIntersects) == 0 {
		if lx < rx {
			return append(accTops, makeSegment(coord{lx, ytop}, coord{rx, ytop}))
		} else {
			return accTops
		}
	}
	if lx > rx {
		panic("lx bigger")
	}
	nextIntersectX := orderedIntersects[0].start.x
	endOfNextIntersect := orderedIntersects[0].end.x

	if lx == nextIntersectX {
		// recurse, did not make any
		return findNextTopFast(endOfNextIntersect, rx, ytop, orderedIntersects[1:], accTops)
	}
	if lx+1 != nextIntersectX {
		// recurse, found a chunk before next intersect
		newS := makeSegment(coord{lx, ytop}, coord{nextIntersectX, ytop})
		return findNextTopFast(endOfNextIntersect, rx, ytop, orderedIntersects[1:], append(accTops, newS))
	}
	// end of previosu intersect immedietaly abuts next one, no intersect
	return findNextTopFast(endOfNextIntersect, rx, ytop, orderedIntersects[1:], accTops)
}

// return segements (interalluy ordered) if one end is within, or is exact wodth of lx and rx
// needs to walk down from starting top y, between l and r, find first y that has intersect or pass min or leeft and rigth downness
func findIntersecting(top, left, right segment, segments map[segment]int) ([]segment, int) {
	res := []segment{}
	lx := left.start.x
	rx := right.end.x

	ytop := slices.Min([]int{left.end.y, right.end.y})

	for y := top.start.y; y <= slices.Min([]int{left.end.y, right.end.y}); y++ {
		if len(res) > 0 {
			return res, ytop
		}
		for k, _ := range segments {
			if k.isHorizontal() {
				if k.getHorizontalLinePos() == y {
					if k != top {
						if k.start.x > lx && k.start.x < rx || k.end.x > lx && k.end.x < rx || k.start.x == lx && k.end.x == rx || k.start.x == rx && k.end.x == lx {
							res = append(res, k)
							ytop = y
						}
					}
				}
			}
		}
	}
	return res, ytop
}

func makeLeftRigthSegments(left, right segment, newTopY int) []segment {
	// if l/r bar origianlly descends past the newY, then need to create the section of line starting at newY and going down to whre bar origianlly ended
	segments := []segment{}
	if left.end.y > newTopY {
		segments = append(segments, makeSegment(coord{left.start.x, newTopY}, left.end))
	}
	if right.end.y > newTopY {
		segments = append(segments, makeSegment(coord{right.start.x, newTopY}, right.end))
	}
	return segments
}

func mergeHorizontals(segs map[segment]int) map[segment]int {
	for k1, _ := range segs {
		for k2, _ := range segs {
			if k1 != k2 && k1.isHorizontal() && k2.isHorizontal() && k1.getHorizontalLinePos() == k2.getHorizontalLinePos() {
				// not same and both horizontal
				if k1.start == k2.end {
					// merge
					delete(segs, k1)
					delete(segs, k2)
					newS := makeSegment(k2.start, k1.end)
					segs[newS] = 1
				}
				if k1.end == k2.start {
					delete(segs, k1)
					delete(segs, k2)
					newS := makeSegment(k1.start, k2.end)
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
		newS := makeSegment(cs[i], cs[j])
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
		if k.isHorizontal() {
			if highest.lineY == -1 {
				highest = k
			}
			if k.getHorizontalLinePos() < highest.getHorizontalLinePos() {
				highest = k
			}
		}
	}
	return highest
}

func findTopBars(top segment, segments map[segment]int) (segment, segment) {
	var left segment
	var right segment
	for k, _ := range segments {
		if k.isVertical() {
			if k.start == top.start {
				left = k
			}
			if k.start == top.end {
				right = k
			}
		}
	}
	if left.end.x == 0 || right.end.x == 0 {
		panic("bad, l/r is not filled in")
	}
	return left, right
}

// ======== helpers

// only for horizontal lines, returns the common y pos
func (s segment) getHorizontalLinePos() int {
	if s.lineY == -1 {
		panic("trying to use yline from a vertical coord")
	}
	return s.lineY
}

// only for vertical lines, returns the common x pos
func (s segment) getVerticalLinePos() int {
	if s.lineX == -1 {
		panic("trying to use xline from a vertical coord")
	}
	return s.lineX
}

func (s segment) isVertical() bool {
	return s.lineX != -1
}

func (s segment) isHorizontal() bool {
	return s.lineY != -1
}

// If horizontal, left to right. If vertical, up to down.
func (s segment) order() segment {
	if s.isHorizontal() && s.start.x > s.end.x || s.isVertical() && s.start.y > s.end.y {
		return s.reverse()
	}
	return s
}

func (s segment) reverse() segment {
	return makeSegment(s.end, s.start)
}

func makeSegment(c1 coord, c2 coord) segment {
	if c1.x == c2.x {
		// vertical seg
		return segment{c1, c2, c1.x, -1}.order()
	}
	if c1.y == c2.y {
		// horizotnal seg
		return segment{c1, c2, -1, c1.y}.order()
	}
	panic("tryign to create segment from 2 coords which share neitehr x nor y")
}
