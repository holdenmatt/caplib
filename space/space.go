package space

import (
	"fmt"

	"github.com/holdenmatt/util"
)

// Space defines a d-dimensional vector space over ℤ/3ℤ.
type Space struct {
	D    int              // Dimension
	Vecs *CoordinatesModN // The space of coordinate d-vectors

	Pts  []int   // Instead of dealing with vecs, we work with their indices ("points")
	Inv  []int   // Map each pt p -> -p
	Sum  [][]int // Map (p, q) -> p + q
	Elim [][]int // Map (p, q) -> the eliminated pt (the one that creates a line)

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
		for i := 0; i < l; i++ {
			Sum[i] = make([]int, l)
			Elim[i] = make([]int, l)
			for j := 0; j < l; j++ {
				sum := vecs.Sum(i, j)
				Sum[i][j] = sum
				Elim[i][j] = vecs.Inv(sum)
			}
		}

		StdBasis := vecs.StdBasis()
		Directions := vecs.Directions()

		spaceCache[d] = &Space{d, vecs, Pts, Inv, Sum, Elim, StdBasis, Directions}
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

// IsSymmetric returns true iff pts is closed under inverses.
func (s Space) IsSymmetric(pts []int) bool {
	for _, p := range pts {
		if !util.Contains(pts, s.Inv[p]) {
			return false
		}
	}
	return true
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

// NewPoints creates a new Points in a Space.
func (s *Space) NewPoints(pts []int) *Points {
	return NewPoints(s, pts)
}

// IsSym returns true iff a Points is symmetric across the origin.
func (pts Points) IsSym() bool {
	for _, p := range pts.Pts {
		pInv := pts.Space.Inv[p]
		if p < pInv {
			if !util.Contains(pts.Pts, pInv) {
				return false
			}
		}
	}
	return true
}

// IsCap returns true iff a Points contains no lines.
func (pts Points) IsCap() bool {
	for _, p := range pts.Pts {
		for _, q := range pts.Pts {
			if p < q {
				r := pts.Space.Elim[p][q]
				if util.Contains(pts.Pts, r) {
					fmt.Println(p, q, r)
					return false
				}
			}
		}
	}
	return true
}

// String returns the default string representation of Points.
func (pts Points) String() string {
	return fmt.Sprintf("Points[%v]", util.Join(pts.Pts, ", "))
}
