package space

import (
	"fmt"
	"sort"

	"github.com/willf/bitset"
)

// Planes represents all planes in a Space through the origin, and provides
// methods to compute plane intersection counts.
type Planes struct {
	space      *Space
	planes     []bitset.BitSet // For each direction, store the orthogonal plane directions.
	directions bitset.BitSet   // Reusable buffer for directions.
}

// NewPlanes creates a new Planes.
func NewPlanes(s *Space) *Planes {
	var planes []bitset.BitSet

	l := uint(len(s.Directions))

	for _, p := range s.Directions {
		plane := bitset.New(l)
		for i, q := range s.Directions {
			if s.Vecs.Dot(p, q) == 0 {
				plane.Set(uint(i))
			}
		}
		planes = append(planes, *plane)
	}

	directions := bitset.New(l)

	return &Planes{s, planes, *directions}
}

func (p Planes) planeCount(directions *bitset.BitSet, index int) int {
	plane := p.planes[index]
	return 2 * int(plane.IntersectionCardinality(directions))
}

// CountExceeds returns true iff a given directions BitSet has a plane exceeding max.
func (p Planes) CountExceeds(directions *bitset.BitSet, max int) bool {
	for i := range p.planes {
		count := p.planeCount(directions, i)
		if count > max {
			return true
		}
	}
	return false
}

// PlaneCountsString returns plane counts as a string "[keys] => [values]" sorted by key.
func (p Planes) PlaneCountsString(directions *bitset.BitSet) string {
	// Map each plane count to the # of planes through the origin with that count.
	counts := make(map[int]int)
	for i := range p.planes {
		count := p.planeCount(directions, i)
		if _, exists := counts[count]; exists {
			counts[count]++
		} else {
			counts[count] = 1
		}
	}

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
