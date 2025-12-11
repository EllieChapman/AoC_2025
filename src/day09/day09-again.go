package day09

import (
	"AoC_2025/src/utils"
	"fmt"
	"slices"
	"sort"
)

func Day9_part2try4(input []string) int {
	coords := utils.Map(input, parse)
	descendingRecatanglesToCheck := rectangleSort(getRectangles(coords))
	fmt.Println(len(descendingRecatanglesToCheck)) // 28 pairwise checks, good (8*7)/2
	segments := parseSegements(coords)
	rs := loop(segments, []rectangle{})
	fmt.Println("[{{7 1} {11 3} 0} {{2 3} {11 5} 0} {{9 5} {11 7} 0}]")
	fmt.Println(rs)
	r := getLargestPossibleRectangle3(descendingRecatanglesToCheck, rs)
	return r.area
}

func loop(segs map[segment]int, accRec []rectangle) []rectangle {
	// segs = utils.Map(segs, func(s segment) segment { return s.order() }) // to do
	// fmt.Println("before order len:", len(segs))
	segs = orderSegMap(segs)
	// fmt.Println("after order len:", len(segs))
	topS := findTopHorizontal(segs)
	left, right := findTopBars(topS, segs)
	intersecting := findIntersecting(topS, left, right, segs)
	// fmt.Println("intersecting", intersecting)
	// fmt.Println(left, topS, right)
	// fmt.Println("looping, before call removeChunk, len(egs):", len(segs), segs)
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
	for _, r := range toCheckrs {
		// once fix updating red green, also fix outside updating rather than starting empty for eahc rectangle
		if isRectanglePossibleNew(r, safeRS) {
			return r
		}
	}
	panic("none possible")
}

func isRectanglePossibleNew(r rectangle, safeRS []rectangle) bool {
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
	rectangle := rectangle{top.start, coord{rx, newTopY}, 0}

	leftRightSegs := makeLeftRigthSegments(left, right, newTopY)
	newSegments = slices.Concat(newSegments, leftRightSegs)

	return newSegments, rectangle
}

func makeTopSegments(newTopY, lx, rx int, intersecting []segment) []segment {
	// EBC todo need to deal with whiole intersect. Intersect includes end coords, but we dont pick up unless one end is wholy within range being considered
	// intersectingCoords := getIntersectingXs(intersecting)
	intersects := order(intersecting)

	if lx > rx {
		lx, rx = rx, lx
	}
	// return findNextTop(lx, rx, newTopY, intersectingCoords, []segment{})
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

// EBC todo broken
// need interscets to be ordered in list, and start and end for each to be ordered
func findNextTopFast(lx, rx, ytop int, orderedIntersects []segment, accTops []segment) []segment {
	// fmt.Println("findNextTopFast", lx, rx, ytop)

	if len(orderedIntersects) == 0 {
		if lx < rx {
			return append(accTops, segment{coord{lx, ytop}, coord{rx, ytop}, -1, ytop})
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
		newS := segment{coord{lx, ytop}, coord{nextIntersectX, ytop}, -1, ytop}
		return findNextTopFast(endOfNextIntersect, rx, ytop, orderedIntersects[1:], append(accTops, newS))
	}
	// send of previosu intersect immedietaly abuts next one, no intersect
	return findNextTopFast(endOfNextIntersect, rx, ytop, orderedIntersects[1:], accTops)
}

func (s segment) order() segment {
	if s.start.x < s.end.x {
		return s
	}
	return segment{s.end, s.start, s.lineX, s.lineY}
}

func (s segment) reverse() segment {
	return segment{s.end, s.start, s.lineX, s.lineY}
}

// // EBC todo This is just too slow for main, going one by one. Could jump by more if has ordered intersect ranges instead
// func findNextTop(lx, rx, ytop int, intersectXs map[int]bool, accTops []segment) []segment {
// 	fmt.Println("findNextTop", lx, rx, ytop)
// 	if !(lx < rx-1) {
// 		return accTops
// 	}
// 	tmp := segment{coord{lx, ytop}, coord{lx, ytop}, -1, ytop}
// 	for ii := lx + 1; ii < rx; ii++ {
// 		_, isIntersecting := intersectXs[ii]
// 		if isIntersecting {
// 			if tmp.start.x == tmp.end.x {
// 				return findNextTop(lx+1, rx, ytop, intersectXs, accTops)
// 			} else {
// 				return findNextTop(lx+1, rx, ytop, intersectXs, append(accTops, segment{coord{lx, ytop}, coord{ii, ytop}, -1, ytop}))
// 			}
// 		} else {
// 			tmp.end.x = ii
// 		}
// 	}
// 	return findNextTop(rx, rx, ytop, intersectXs, append(accTops, segment{coord{lx, ytop}, coord{rx, ytop}, -1, ytop}))
// }

// func getIntersectingXs(intersecting []segment) map[int]bool {
// 	xintersects := map[int]bool{}
// 	for _, i := range intersecting {
// 		s := i.start.x
// 		e := i.end.x
// 		if i.start.x > i.end.x {
// 			s = i.end.x
// 			e = i.start.x
// 		}
// 		for ii := s; ii <= e; ii++ {
// 			xintersects[ii] = true
// 		}
// 	}
// 	return xintersects
// }

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
		for k, _ := range segments {
			if k.isHorizontal() {
				if k.getHorizontalLinePos() == y {
					if k != top {
						if k.start.x > lx && k.start.x < rx || k.end.x > lx && k.end.x < rx || k.start.x == lx && k.end.x == rx || k.start.x == rx && k.end.x == lx {
							res = append(res, k.order())
							// res = append(res, segment{k.end, k.start, k.lineX, k.lineY})
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
		segments = append(segments, segment{left.start, coord{left.start.x, newTopY}, left.start.x, -1})
	}
	if right.end.y > newTopY {
		segments = append(segments, segment{coord{right.start.x, newTopY}, right.end, right.start.x, -1})
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
					newS := segment{k2.start, k1.end, -1, k1.getHorizontalLinePos()}
					segs[newS] = 1
				}
				if k1.end == k2.start {
					delete(segs, k1)
					delete(segs, k2)
					newS := segment{k1.start, k2.end, -1, k2.getHorizontalLinePos()}
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

// ebc directionality might become a problem. ie might have two bars that connect via their ends. but code above needs left and right with start and end as expected
// if k.start = top.start for exmaple, need to reverse k before set as left. prints for now
func findTopBars(top segment, segments map[segment]int) (segment, segment) {
	var left segment
	var right segment
	for k, _ := range segments {
		if k.isVertical() {
			if k.end == top.start {
				left = k
			}
			if k.start == top.start {
				// fmt.Println("need to implement 1") // Todo almost 100%
				// panic("here1")
				left = k.order()
			}
			if k.start == top.end {
				right = k
			}
			if k.end == top.end {
				// fmt.Println("need to implement 2")
				// panic("here2")
				right = k.reverse()
			}
		}
	}
	return left, right
}

// ======== helpers

func (s segment) getVerticalPos() int {
	if s.lineX == -1 {
		panic("trying to use xline from a horoizontal coord")
	}
	return s.lineX
}

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
