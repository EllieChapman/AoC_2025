package day08

import (
	"AoC_2025/src/utils"
	"slices"
	"sort"
	"strings"
)

func Day8_part1(input []string) int {
	coords := utils.Map(input, parseCoord)
	allPossiblelinksOrdered := mySort(findAllLinks(coords))
	n := 10
	if len(input) != 20 {
		n = 1000
	}
	circuits, _ := createCircuits(allPossiblelinksOrdered[0:n])
	lens := utils.Map(circuits, func(m map[coord]string) int { return len(m) })
	sort.Ints(lens)
	lens = utils.Reverse(lens)[0:3]
	return utils.Mul(lens)
}

func Day8_part2(input []string) int {
	coords := utils.Map(input, parseCoord)
	allPossiblelinksOrdered := mySort(findAllLinks(coords))
	_, lastLinkAdded := createCircuits(allPossiblelinksOrdered)
	return lastLinkAdded.a.x * lastLinkAdded.b.x
}

type coord struct {
	x int
	y int
	z int
}

type link struct {
	a        coord
	b        coord
	distance int
}

func createCircuits(orderedLinks []link) ([]map[coord]string, link) {
	// everythingLinked := map[coord]string{}
	circuits := []map[coord]string{}
	lastAdded := link{}
	var added bool
	for _, l := range orderedLinks {
		circuits, added = addLink(l, circuits)
		if added {
			lastAdded = l
		}
	}

	return circuits, lastAdded
}

func addLink(l link, cs []map[coord]string) ([]map[coord]string, bool) {
	// find, if any cicruit containing l.a
	// find, if any cicruit containing l.b
	aPos := -1
	bPos := -1
	for pos, circuit := range cs {
		_, aok := circuit[l.a]
		_, bok := circuit[l.b]
		if aok {
			aPos = pos
		}
		if bok {
			bPos = pos
		}
	}
	// options:
	if aPos == -1 && bPos == -1 {
		// 1. neitehr a nor b in any cicuit -> create a new cicruit with a and b, update lastAdded
		return append(cs, map[coord]string{l.a: "", l.b: ""}), true
	}
	if aPos == -1 || bPos == -1 {
		// 2. only one of a or b is already in a circuit
		if aPos != -1 {
			// 2a. only a in a cicuit -> add b to that cicuit updat lastAdded
			newM := cs[aPos]
			newM[l.b] = ""
			newCs := slices.Concat(cs[0:aPos], []map[coord]string{newM}, cs[aPos+1:])
			return newCs, true
		}
		if bPos != -1 {
			// 2b. only b in a cicuit -> add a to that cicuit updat lastAdded
			newM := cs[bPos]
			newM[l.a] = ""
			newCs := slices.Concat(cs[0:bPos], []map[coord]string{newM}, cs[bPos+1:])
			return newCs, true
		}
	}
	if aPos != bPos {
		// 3. both in diff cicuit -> merge cicuits, updat lastAdded
		return mergeCircuits(aPos, bPos, cs), true
	}
	// 4. both must be in same cicruti -> nothing to do, don't update lastAdded
	return cs, false
}

func mergeCircuits(a int, b int, cs []map[coord]string) []map[coord]string {
	combined := merge(cs[a], cs[b])
	if a < b {
		cs = slices.Concat(cs[0:a], cs[a+1:b], cs[b+1:])
	} else {
		cs = slices.Concat(cs[0:b], cs[b+1:a], cs[a+1:])
	}
	return append(cs, combined)
}

func merge(a map[coord]string, b map[coord]string) map[coord]string {
	for k, v := range a {
		b[k] = v
	}
	return b
}

// the actual distance is square root of returned value. But to compare distances can just use squared value as is an int
func calcDistanceSquared(a coord, b coord) int {
	xdiff := (a.x - b.x) * (a.x - b.x)
	ydiff := (a.y - b.y) * (a.y - b.y)
	zdiff := (a.z - b.z) * (a.z - b.z)
	return xdiff + ydiff + zdiff
}

func parseCoord(s string) coord {
	is := utils.Map(strings.Split(s, ","), utils.Atoi)
	return coord{is[0], is[1], is[2]}
}

func findAllLinks(cs []coord) []link {
	ls := []link{}
	for i := 0; i < len(cs); i++ {
		for j := i + 1; j < len(cs); j++ {
			ls = append(ls, makeLink(cs[i], cs[j]))
		}
	}
	return ls
}

func makeLink(a coord, b coord) link {
	return link{a, b, calcDistanceSquared(a, b)}
}

func mySort(ls []link) []link {
	sort.Slice(ls, func(i, j int) bool {
		return ls[i].distance < ls[j].distance
	})
	return ls
}
