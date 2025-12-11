package day09

import (
	"AoC_2025/src/utils"
	"fmt"
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
	// for pos, r := range toCheckrs {
	startPos := 0
	if len(toCheckrs) != 28 {
		startPos = 49000
	}
	for pos, r := range toCheckrs[startPos:] { // TODO remove startPos hack once isRectanglePossible is faster. Start just before where it is.
		if pos%100 == 0 {
			fmt.Println(len(toCheckrs) - pos)
		}
		if isRectanglePossible(r, safeRS) {
			return r
		}
	}
	panic("none possible")
}

// TODO - this is the slow part. Generating safeRS is fast, checking test rectangles is not. Need to check segments, not each point separately.
// ?? is it true that for each of 4 test rec segement, each segemnt must be is entriely cotained in a single safe chunk? no
func isRectanglePossible(r rectangle, safeRS []rectangle) bool {
	c1 := coord{r.a.x, r.a.y}
	c2 := coord{r.a.x, r.b.y}
	c3 := coord{r.b.x, r.b.y}
	c4 := coord{r.b.x, r.a.y}
	greenLineCoords := slices.Concat(getLineCoords(c1, c2), getLineCoords(c2, c3), getLineCoords(c3, c4), getLineCoords(c4, c1))
	for _, c := range greenLineCoords {
		coordOk := checkIfCoordIsOk(c, safeRS)
		if !coordOk {
			return false
		}
	}
	return true
}

// needs to check segments, not coords
// func isRectanglePossibleFast(r rectangle, safeRS []rectangle) bool {
// 	c1 := coord{r.a.x, r.a.y}
// 	c2 := coord{r.a.x, r.b.y}
// 	c3 := coord{r.b.x, r.b.y}
// 	c4 := coord{r.b.x, r.a.y}
// 	segsToCheck := []segement{makeSegment(c1, c2), makeSegment(c2, c3), makeSegment(c3, c4), makeSegment(c4, c1)}

// 	return true
// }

func checkIfCoordIsOk(c coord, safeRS []rectangle) bool {
	for _, r := range safeRS {
		if r.contains(c) {
			return true
		}
	}
	return false
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

func (r rectangle) contains(c coord) bool {
	leftLimit := r.a.x
	rightLimit := r.b.x
	if leftLimit > rightLimit {
		leftLimit, rightLimit = rightLimit, leftLimit
	}
	upLimit := r.a.y
	downLimit := r.b.y
	if upLimit > downLimit {
		upLimit, downLimit = downLimit, upLimit
	}
	if c.x >= leftLimit && c.x <= rightLimit && c.y >= upLimit && c.y <= downLimit {
		return true
	}
	return false
}
