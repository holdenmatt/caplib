package space

import (
	"fmt"
	"sort"
	"strings"

	"github.com/holdenmatt/util"
)

// Space defines a d-dimensional vector space over ℤ/3ℤ.
type Space struct {
	D    int              // Dimension
	Vecs *CoordinatesModN // The space of coordinate d-vectors

	Pts  []int    // Instead of dealing with vecs, we work with their indices ("points")
	Inv  []int    // Map each pt p -> -p
	Sum  [][]int  // Map (p, q) -> p + q
	Elim [][]int  // Map (p, q) -> the eliminated pt (the one that creates a line)
	Perp [][]bool // Map (p, q) -> true iff p and q are orthogonal.

	StdBasis   []int // Standard basis indices
	Directions []int // Indices of unique directions
}

// Cache a unique Space for each dimension.
var spaceCache = make(map[int]*Space)

// New creates a new ternary Space of dimension d.
func New(d int) *Space {
	if _, ok := spaceCache[d]; !ok {
		vecs := NewCoordinatesModN(3, d)
		l := len(vecs.Vecs)

		// Compute pts/inverses.
		Pts := make([]int, l)
		Inv := make([]int, l)
		for i := range vecs.Vecs {
			Pts[i] = i
			Inv[i] = vecs.Inv(i)
		}

		// Compute sum/elim/dot.
		Sum := make([][]int, l)
		Elim := make([][]int, l)
		Perp := make([][]bool, l)
		for i := 0; i < l; i++ {
			Sum[i] = make([]int, l)
			Elim[i] = make([]int, l)
			Perp[i] = make([]bool, l)
			for j := 0; j < l; j++ {
				sum := vecs.Sum(i, j)
				Sum[i][j] = sum
				Elim[i][j] = vecs.Inv(sum)
				Perp[i][j] = vecs.Dot(i, j) == 0
			}
		}

		StdBasis := vecs.StdBasis()
		Directions := vecs.Directions()

		spaceCache[d] = &Space{d, vecs, Pts, Inv, Sum, Elim, Perp, StdBasis, Directions}
	}

	return spaceCache[d]
}

// String returns the default string representation of a Space.
func (s Space) String() string {
	return fmt.Sprintf("Space[d = %d]", s.D)
}

// Size returns the number of pts (or vectors) in a Space.
func (s Space) Size() int {
	return len(s.Pts)
}

// Span returns the span of vectors with given indices,
// in the order of coefficient vectors.
func (s Space) Span(indices []int) []int {
	if len(indices) == 0 {
		return []int{ORIGIN}
	}

	last := indices[len(indices)-1]
	head := indices[:len(indices)-1]

	lastInv := s.Inv[last]
	headSpan := s.Span(head)

	if util.Contains(headSpan, last) {
		return headSpan
	}

	// Append headSpan, headSpan + last, headSpan + lastInv
	span := make([]int, 0, 3*len(headSpan))
	span = append(span, headSpan...)
	for _, p := range headSpan {
		span = append(span, s.Sum[p][last])
	}
	for _, p := range headSpan {
		span = append(span, s.Sum[p][lastInv])
	}
	return span
}

// LinearCombo computes the linear combination of pts with the given coeffs.
func (s Space) LinearCombo(pts []int, coeffs []int) int {
	if len(pts) != len(coeffs) {
		panic("length mismatch")
	}

	if len(pts) == 0 {
		return ORIGIN
	}

	res := s.LinearCombo(pts[1:], coeffs[1:])

	pt0 := pts[0]
	c0 := coeffs[0]
	if c0 == 1 {
		res = s.Sum[res][pt0]
	} else if c0 == 2 {
		res = s.Sum[res][pt0]
		res = s.Sum[res][pt0]
	} else if c0 != 0 {
		panic("coeffs must be 0, 1, or 2")
	}

	return res
}

//
//--- Points ---//
//

// Points represents a vector of sorted pts (indices) in a Space.
type Points struct {
	Space *Space // The enclosing Space
	Pts   []int  // Sorted pt indices
}

// NewPoints creates a new Points in a Space.
func NewPoints(s *Space, pts []int) *Points {
	return &Points{s, pts}
}

// String returns the default string representation of Points.
func (pts Points) String() string {
	return fmt.Sprintf("Points%v", pts.Pts)
}

// PlaneCount returns the count of pts in the plane orthogonal to a given normal pt.
func (pts Points) PlaneCount(normal int) int {
	if normal == ORIGIN {
		panic("normal must be nonzero")
	}

	count := 0
	perp := pts.Space.Perp[normal]
	for _, p := range pts.Pts {
		if perp[p] {
			count++
		}
	}
	return count
}

// PlaneCounts maps each plane count to the # of planes through the origin
// with that count.
//
// This is invariant under linear isomorphisms, so if two Points have differing
// PlaneCounts they cannot be isomorphic.
func (pts Points) PlaneCounts() map[int]int {
	counts := make(map[int]int)
	for _, normal := range pts.Space.Directions {
		count := pts.PlaneCount(normal)
		if _, exists := counts[count]; exists {
			counts[count]++
		} else {
			counts[count] = 1
		}
	}
	return counts
}

// PlaneCountsString returns PlaneCounts as a string "[keys] => [values]"
// sorted by key.
func (pts Points) PlaneCountsString() string {
	counts := pts.PlaneCounts()

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

//
//-- Encode Points as a graphical ASCII string --//
//

const pointSymbol = "#"
const emptySymbol = "-"

// Encode some Points as an ASCII string.
func (pts Points) Encode() string {
	lines := encodeLines(pts.Space.D, pts.Pts)
	return fmt.Sprintf("%v\n%v\n", pts, strings.Join(lines, "\n"))
}

// encodeLines returns a string representation of pts, as a slice of lines.
// In alternate dimensions, we move to the right and downward on the page.
func encodeLines(d int, pts []int) []string {
	if d == 0 {
		if util.Contains(pts, ORIGIN) {
			return []string{pointSymbol}
		}
		return []string{emptySymbol}
	}

	// Get string representations for the 3 projections.
	proj := projections(d, pts)
	lines0 := encodeLines(d-1, proj[0])
	lines1 := encodeLines(d-1, proj[1])
	lines2 := encodeLines(d-1, proj[2])

	var lines []string
	if d%2 == 1 {
		// In odd dimensions, combine arrays of lines (going down).
		lines = append(lines, lines0...)
		for i := 0; i < numSpacers(d); i++ {
			lines = append(lines, "")
		}
		lines = append(lines, lines1...)
		for i := 0; i < numSpacers(d); i++ {
			lines = append(lines, "")
		}
		lines = append(lines, lines2...)
	} else {
		// In even dimensions, combine corresponding lines (going rightward).
		for i := range lines0 {
			spacer := strings.Repeat(" ", numSpacers(d))
			line := fmt.Sprintf("%v%v%v%v%v", lines0[i], spacer, lines1[i], spacer, lines2[i])
			lines = append(lines, line)
		}
	}
	return lines
}

// Split pts into 3 sets by projecting along the last dimension.
func projections(d int, pts []int) [][]int {
	oneThird := util.Pow(3, d-1)
	projections := [][]int{{}, {}, {}}

	for _, p := range pts {
		index := p / oneThird
		value := p % oneThird
		projections[index] = append(projections[index], value)
	}
	return projections
}

// numSpacers returns the number of blank spacers used to separate
// encodings of larger dimension d.
// For d = 1,2,3,..., the sequence is 0, 1, 1, 4, 4, 9, 9...
func numSpacers(d int) int {
	var half int
	if d%2 == 0 {
		half = d / 2
	} else {
		half = (d - 1) / 2
	}
	return half * half
}
