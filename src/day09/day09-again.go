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
	// descendingSafeChunks := rectangleSort(safeChunks)
	r := getLargestPossibleRectangle3(descendingRecatanglesToCheck, safeChunks)
	return r.area
}

func loop(segs map[segment]int, accRec []rectangle) []rectangle {
	segs = orderSegMap(segs)
	topS := findTopHorizontal(segs)
	left, right := findTopBars(topS, segs)
	intersecting := findIntersecting(topS, left, right, segs)
	newsegs, removedInsideChunk := removeChunk(left, topS, right, segs, intersecting)
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

func orderSegMap(m map[segment]int) map[segment]int {
	newm := map[segment]int{}
	for k, _ := range m {
		newm[k.order()] = 1
	}
	return newm
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
		if !isSegInside([]segment{seg.order()}, safeRS) {
			return false
		}
	}
	return true
}

func isSegInside(fragmentsToCheck []segment, safeRS []rectangle) bool {
	for _, r := range safeRS {
		newFrags := []segment{}
		for _, f := range fragmentsToCheck {
			_, new := r.containsSeg(f)
			newFrags = slices.Concat(newFrags, new)
		}
		fragmentsToCheck = newFrags
		if len(fragmentsToCheck) == 0 {
			return true
		}
	}
	return false // run out of safe rectabgles but still have somefragements to fit
}

// return false if no part of segemnt is contained in rectangle, or
// return true with list of 0, 1 or 2 segements, which are fragments of the orgianl segment that were not contained
// NB seg is interanlly x ordered. make sure returned segs are interally x ordered
func (r rectangle) containsSeg(s segment) (bool, []segment) {
	if s.isHorizontal() {
		// y is either in or out
		y := s.getHorizontalLinePos()
		y1 := r.a.y
		y2 := r.b.y
		if y1 > y2 {
			y1, y2 = y2, y1
		}
		if !(y >= y1 && y <= y2) {
			return false, []segment{s}
		}
		// only care about the x values in the rectagle for ma,ing chunks
		x1 := r.a.x
		x2 := r.b.x
		if x1 > x2 {
			panic("rectabgke xs not ordered")
		}
		sx1 := s.start.x
		sx2 := s.end.x
		if sx1 > sx2 {
			panic("segment xs not ordered")
		}
		anyOverlap, oneDranges := compareRanges(sx1, sx2, x1, x2)
		segs := []segment{}
		for _, o := range oneDranges {
			segs = append(segs, makeSegment(coord{o.a, y}, coord{o.b, y}))
		}
		return anyOverlap, segs
	} else {
		// x is either in or out
		x := s.lineX
		x1 := r.a.x
		x2 := r.b.x
		if x1 > x2 {
			x1, x2 = x2, x1
		}
		if !(x >= x1 && x <= x2) {
			return false, []segment{s}
		}
		// only care about the y values in the rectagle
		y1 := r.a.y
		y2 := r.b.y
		if y1 > y2 {
			panic("rectabgke ys not ordered")
		}
		sy1 := s.start.y
		sy2 := s.end.y
		if sy1 > sy2 {
			sy1, sy2 = sy2, sy1
		}
		anyOverlap, oneDranges := compareRanges(sy1, sy2, y1, y2)
		segs := []segment{}
		for _, o := range oneDranges {
			segs = append(segs, makeSegment(coord{x, o.a}, coord{x, o.b}))
		}
		return anyOverlap, segs
	}
}

type oneDRange struct {
	a int
	b int
}

// return segs not included in parant, and if there was any overlap
func compareRanges(test1, test2, parent1, parent2 int) (bool, []oneDRange) {
	if test2 < parent1 || test1 > parent2 {
		// non included
		return false, []oneDRange{{test1, test2}}
		// return false, []oneDRange{} // change for 2 mthdos
	}
	if test1 >= parent1 && test2 <= parent2 {
		// all included
		return true, []oneDRange{}
	}
	if test1 >= parent1 && test1 <= parent2 && test2 > parent2 {
		// will generate one fragement, beyond parent 2
		return true, []oneDRange{{parent2 + 1, test2}}
	}
	if test2 >= parent1 && test2 <= parent2 && test1 < parent1 {
		// will generate one fragement, before parent 1
		return true, []oneDRange{{test1, parent1 - 1}}
	}
	if test1 < parent1 && test2 > parent2 {
		// two frgament, on eitehr sid eof th eparents
		return true, []oneDRange{{parent2 + 1, test2}, {test1, parent1 - 1}}
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
func removeChunk(left, top, right segment, segments map[segment]int, intersecting []segment) (map[segment]int, rectangle) {
	newSegments, rectangleBeingRemoved := createNewSegements(left, top, right, intersecting)

	delete(segments, left.order())
	delete(segments, top.order())
	delete(segments, right.order())
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

func createNewSegements(left, top, right segment, intersecting []segment) ([]segment, rectangle) {
	newSegments := []segment{}
	newTopY := getNewY(left, right, intersecting)
	rx := right.start.x
	lx := left.start.x
	if lx > rx {
		lx, rx = rx, lx
	}

	newTopSegments := makeTopSegments(newTopY, lx, rx, intersecting)

	newSegments = slices.Concat(newSegments, newTopSegments)
	rectangle := calcRectangleSize(top.start, coord{rx, newTopY})

	leftRightSegs := makeLeftRigthSegments(left, right, newTopY)
	newSegments = slices.Concat(newSegments, leftRightSegs)

	return newSegments, rectangle
}

func makeTopSegments(newTopY, lx, rx int, intersecting []segment) []segment {
	intersects := order(intersecting)

	if lx > rx {
		lx, rx = rx, lx
	}
	return findNextTopFast(lx, rx, newTopY, intersects, []segment{})
}

func order(intersecting []segment) []segment {
	new := []segment{}
	for _, i := range intersecting {
		internallyOrdered := i.order()
		new = append(new, internallyOrdered)
	}
	return intersectSort(new)
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
func findIntersecting(top, left, right segment, segments map[segment]int) []segment {
	res := []segment{}
	lx := left.start.x
	rx := right.end.x
	if lx > rx {
		lx, rx = rx, lx
	}
	for y := top.start.y; y <= slices.Min([]int{left.start.y, right.end.y}); y++ {
		if len(res) > 0 {
			return res
		}
		for k, _ := range segments {
			if k.isHorizontal() {
				if k.getHorizontalLinePos() == y {
					if k != top {
						if k.start.x > lx && k.start.x < rx || k.end.x > lx && k.end.x < rx || k.start.x == lx && k.end.x == rx || k.start.x == rx && k.end.x == lx {
							res = append(res, k.order())
						}
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
		segments = append(segments, makeSegment(left.start, coord{left.start.x, newTopY}))
	}
	if right.end.y > newTopY {
		segments = append(segments, makeSegment(coord{right.start.x, newTopY}, right.end))
	}
	return segments
}

func getNewY(l segment, r segment, intersecting []segment) int {
	if len(intersecting) > 0 {
		return intersecting[0].getHorizontalLinePos()
	}
	if l.start.y > r.end.y {
		return r.end.y
	}
	return l.start.y
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
			if k.end == top.start {
				left = k
			}
			if k.start == top.start {
				left = k.order()
			}
			if k.start == top.end {
				right = k
			}
			if k.end == top.end {
				right = k.reverse()
			}
		}
	}
	if left.end.x == 0 || right.end.x == 0 {
		panic("bad, l/r is not filled in")
	}
	return left, right
}

// ======== helpers

func (s segment) getHorizontalLinePos() int {
	if s.lineY == -1 {
		panic("trying to use yline from a vertical coord")
	}
	return s.lineY
}

func (s segment) isVertical() bool {
	return s.lineX != -1
}

func (s segment) isHorizontal() bool {
	return s.lineY != -1
}

// Todo, amke x reverse and y reverse, and make order work for both horoizontal and verticla segs, insetad of doi gnoithing for vertical lie now
func (s segment) order() segment {
	if s.start.x < s.end.x {
		return s
	}
	return s.reverse()
}

func (s segment) reverse() segment {
	return makeSegment(s.end, s.start)
}

func makeSegment(c1 coord, c2 coord) segment {
	if c1.x == c2.x {
		// vertical seg
		return segment{c1, c2, c1.x, -1}
	}
	if c1.y == c2.y {
		// horizotnal seg
		return segment{c1, c2, -1, c1.y}.order()
	}
	panic("tryign to create segment from 2 coords which share neitehr x nor y")
}
