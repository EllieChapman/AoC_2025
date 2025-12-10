package day09

import (
	"AoC_2025/src/utils"
	"fmt"
)

func Day9_part2try4(input []string) int {
	coords := utils.Map(input, parse)
	descendingRecatanglesToCheck := rectangleSort(getRectangles(coords))
	fmt.Println(len(descendingRecatanglesToCheck)) // 28 pairwise checks, good (8*7)/2
	segments := getSegements(coords)
	fmt.Println(len(segments)) // 8 segments, good
	topS := findTopHorizontal(segments)
	left, right := findTopBars(topS, segments)
	fmt.Println(topS, left, right)
	newsegs, r := removeChunkEasy(left, topS, right, segments)
	fmt.Println(newsegs)
	fmt.Println(r)
	return 0
}

// one horizontal or vertical border segement
type segment struct {
	start coord
	end   coord
	lineY int // if a horizontal line this is the y coord, if a vertical line (ie y changes) make -1
}

// This is the case where the chunk can be removed down until one of the left or right bar segements ends.
// Ie, there is no other line segemetn which is "inside" the chunk being removed
//
//	a         _
//	b        | |_  ->    _ _
//	c       _|         _|
//	d
func removeChunkEasy(left, top, right segment, segments map[segment]int, intersecting []segment) (map[segment]int, rectangle) {
	// create new horizointal segement
	// create new left or right bar - which ever was longer so not fully used up. If both same dont need a new bar as both fully used
	// remove origianl top, left and right from segments.
	// If exists, add new left/right sgement (this can never be merged, is just a shorter version of what already existsed)
	//
	// if can merge new horxitonal with another horizojtal segement, do it (nends ame y and connecting start or end point). create longer continuous lines
	// add new or merged horizontal
	newSegments := []segment{}
	rectangleBeingRemoved := rectangle{top.start, coord{}, 0}
	if left.start.y > right.end.y {
		// left bar reaches down further
		newTop := segment{coord{left.start.x, right.end.y}, coord{right.start.x, right.end.y}, right.end.y}
		newLeft := segment{left.start, coord{left.start.x, right.end.y}, -1}
		newSegments = append(newSegments, newTop)
		newSegments = append(newSegments, newLeft)
		rectangleBeingRemoved.b = newTop.end
	} else if left.start.y < right.end.y {
		// right bar reaches down further
		newTop := segment{coord{left.start.x, left.start.y}, coord{right.start.x, left.start.y}, left.start.y}
		newRight := segment{coord{right.start.x, left.start.y}, right.end, -1}
		newSegments = append(newSegments, newTop)
		newSegments = append(newSegments, newRight)
		rectangleBeingRemoved.b = newTop.end
	} else {
		// they are the same
		newTop := segment{coord{left.start.x, right.end.y}, coord{right.start.x, right.end.y}, right.end.y}
		newSegments = append(newSegments, newTop)
		rectangleBeingRemoved.b = newTop.end
	}
	delete(segments, left)
	delete(segments, top)
	delete(segments, right)
	for _, s := range newSegments {
		segments[s] = 1
	}
	segments = mergeHorizontals(segments)
	return segments, rectangleBeingRemoved
}

// This is the case where the chunk can be removed down only until it hits another line segment contained in the chunk.
// As the intruding chunk can be assumed to be the "outide", only replace the line segement where there was no line previously
// This will also clean up in the case where it hits and collapses its own border
//
//	a      _ _ _                                     _ _ _ _ _
//	b     |  _  |   ->    _   _                     |         |
//	c    _| | | |_      _| | | |_        or         |_ _ _    |    ->    _ _
//	d       | |            | |                            |   |         |   |
//	e
// func removeChunkHard() {
// TODO
// need to know all segments which are in the way, this can be multiple if they have the same y
//
// create new horizointal segements
// segements need to include any x coords from original top which are not present in any the interrupting segments
// should not include any x whihc is in an interruption (excpet for corners)
// Same - create new left or right bar - which ever was longer so not fully used up. If both same dont need a new bar as both fully used
// Same - remove origianl top, left and right from segments.
// Same - If exists, add new left/right sgement (this can never be merged, is just a shorter version of what already existsed)
//
// Sameish - if can merge ANY of the new horxitonal with another horizojtal segement, do it (nends ame y and connecting start or end point). create longer continuous lines
// }

func mergeHorizontals(segs map[segment]int) map[segment]int {
	for k1, _ := range segs {
		for k2, _ := range segs {
			if k1 != k2 && k1.lineY != -1 && k2.lineY != -1 && k1.lineY == k2.lineY {
				// not same and both horizontal
				delete(segs, k1)
				delete(segs, k2)
				if k1.start == k2.end {
					// merge
					newS := segment{k2.start, k1.end, k1.lineY}
					segs[newS] = 1
				}
				if k1.end == k2.start {
					newS := segment{k1.start, k2.end, k2.lineY}
					segs[newS] = 1
				}
			}
		}
	}
	return segs
}

func getSegements(cs []coord) map[segment]int {
	segments := map[segment]int{}
	for i := 0; i < len(cs); i++ {
		j := i + 1
		if i == len(cs)-1 {
			j = 0
		}
		newS := segment{cs[i], cs[j], -1}
		if cs[i].y == cs[j].y {
			newS.lineY = cs[i].y
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
		if k.end == top.start {
			left = k
		}
		if k.start == top.start {
			fmt.Println("need to implement 1")
		}
		if k.start == top.end {
			right = k
		}
		if k.end == top.end {
			fmt.Println("need to implement 2")
		}
	}
	return left, right
}
