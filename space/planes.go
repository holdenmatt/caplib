package space

import (
	"fmt"
	"sort"
)

// Planes represents all planes in a Space through the origin, and provides
// methods to compute plane intersection counts of Points.
type Planes struct {
	space *Space
	perp  [][]bool // Map (p, q) -> true iff p and q are orthogonal.
}

// NewPlanes creates a new Planes.
func NewPlanes(s *Space) *Planes {
	size := s.Size()
	perp := make([][]bool, size)
	for _, p := range s.Pts {
		perp[p] = make([]bool, size)
		for _, q := range s.Pts {
			perp[p][q] = s.Vecs.Dot(p, q) == 0
		}
	}

	return &Planes{s, perp}
}

// PlaneCountsString returns plane counts as a string "[keys] => [values]"
// sorted by key.
func (p Planes) PlaneCountsString(pts *Points) string {
	counts := p.planeCounts(pts)

	// Sort keys.
	var keys []int
	for k := range counts {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	// Sort values by key.
	values := make([]int, 0, len(counts))
	for _, k := range keys {
		values = append(values, counts[k])
	}

	// Print "[keys] => [values]"
	return fmt.Sprintf("%v => %v", keys, values)
}

// CountsAtMost returns true iff all plane intersections of a Points
// are at most a given max.
func (p Planes) CountsAtMost(pts *Points, max int) bool {
	for _, normal := range p.space.Directions {
		count := p.planeCount(pts, normal)

		if count > max {
			return false
		}
	}
	return true
}

// planeCount returns the count of pts in the plane orthogonal to a given normal pt.
func (p Planes) planeCount(pts *Points, normal int) int {
	if normal == ORIGIN {
		panic("normal must be nonzero")
	}

	count := 0
	perp := p.perp[normal]
	for _, p := range pts.Pts {
		if perp[p] {
			count++
		}
	}
	return count
}

// planeCounts maps each plane count to the # of planes through the origin
// with that count.
//
// This is invariant under linear isomorphisms, so if two Points have differing
// PlaneCounts they cannot be isomorphic.
func (p Planes) planeCounts(pts *Points) map[int]int {
	counts := make(map[int]int)
	for _, normal := range pts.Space.Directions {
		count := p.planeCount(pts, normal)
		if _, exists := counts[count]; exists {
			counts[count]++
		} else {
			counts[count] = 1
		}
	}
	return counts
}
