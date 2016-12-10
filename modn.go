package caplib

import "github.com/holdenmatt/util"

//
// Implement the ring ℤ/nℤ and its d-dimensional coordinate space.
//

// ORIGIN is the index of the zero vector in a coordinate space.
const ORIGIN = 0

// IntegersModN represents the ring ℤ/nℤ of integers mod n.
type IntegersModN struct {
	n      int   // The modulus
	values []int // The residues 0,...,n-1
}

// NewIntegersModN creates a new IntegersModN with modulus n.
func NewIntegersModN(n int) *IntegersModN {
	values := util.Range(n)
	return &IntegersModN{n, values}
}

// CoordinatesModN represents the d-dimensional coordinate space (ℤ/nℤ)^d.
// This is a (ℤ/nℤ)-module in general, and a vector space when n is prime.
type CoordinatesModN struct {
	n    int     // The modulus
	d    int     // The dimension
	Vecs [][]int // The d-vectors of coordinates mod n

	buffer []int // Reusable buffer of size d.
}

// NewCoordinatesModN creates a new CoordinatesModN of dimension d and modulus n.
func NewCoordinatesModN(n int, d int) *CoordinatesModN {
	coords := NewIntegersModN(n)
	var vecs [][]int
	for vec := range util.Product(coords.values, d) {
		vecs = append(vecs, vec)
	}

	buffer := make([]int, d)

	return &CoordinatesModN{n, d, vecs, buffer}
}

// VecToIndex computes the index of a coordinate vector.
func (c CoordinatesModN) VecToIndex(vec []int) int {
	if len(vec) != c.d {
		panic("VecToIndex: vec must have length d")
	}

	index := 0
	for _, val := range vec {
		if val < 0 || val >= c.n {
			panic("VecToIndex: value out of bounds")
		}
		index = c.n*index + val
	}
	return index
}

// Inv computes the additive inverse of vector i, and returns its index.
func (c CoordinatesModN) Inv(i int) int {
	inv := c.buffer
	for k := 0; k < c.d; k++ {
		inv[k] = (c.n - c.Vecs[i][k]) % c.n
	}
	return c.VecToIndex(inv)
}

// Sum adds vectors i and j, and returns the index.
func (c CoordinatesModN) Sum(i int, j int) int {
	sum := c.buffer
	for k := 0; k < c.d; k++ {
		sum[k] = (c.Vecs[i][k] + c.Vecs[j][k]) % c.n
	}
	return c.VecToIndex(sum)
}

// Dot computes the dot product of vectors i and j (mod n).
func (c CoordinatesModN) Dot(i int, j int) int {
	dot := 0
	for k := 0; k < c.d; k++ {
		dot += c.Vecs[i][k] * c.Vecs[j][k]
	}
	return dot % c.n
}

// Coords returns the coordinates [0,1,...,d-1] for a space.
func (c CoordinatesModN) Coords() []int {
	return util.Range(c.d)
}

// StdBasis returns the indices of the standard basis for the coordinate space.
func (c CoordinatesModN) StdBasis() []int {
	basis := make([]int, c.d)
	for i := range c.Coords() {
		basis[i] = util.Pow(c.n, i)
	}
	return basis
}

// Directions returns indices for each unique "direction" in the coordinate space,
// i.e. vectors whose first nonzero value is a 1.
func (c CoordinatesModN) Directions() []int {
	var dirs []int
	for i, vec := range c.Vecs {
		for _, value := range vec {
			if value == 1 {
				dirs = append(dirs, i)
				break
			} else if value != 0 {
				break
			}
		}
	}
	return dirs
}

// Span returns the span of vectors with given indices,
// in the order of coefficient vectors.
func (c CoordinatesModN) Span(indices []int) []int {
	if len(indices) == 0 {
		return []int{ORIGIN}
	}

	last := indices[len(indices)-1]
	head := indices[:len(indices)-1]

	lastInv := c.Inv(last)
	headSpan := c.Span(head)

	if util.Contains(headSpan, last) {
		return headSpan
	}

	// Append headSpan, headSpan + last, headSpan + lastInv
	span := make([]int, 0, 3*len(headSpan))
	span = append(span, headSpan...)
	for _, p := range headSpan {
		span = append(span, c.Sum(p, last))
	}
	for _, p := range headSpan {
		span = append(span, c.Sum(p, lastInv))
	}
	return span
}
