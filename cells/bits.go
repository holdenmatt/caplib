package cells

import (
	"github.com/cznic/mathutil"
	"github.com/holdenmatt/caplib/space"
)

// cmp stores precomputed byte comparisons such that cmp[i][j] = Cmp(i, j).
var cmp [256][256]int

func init() {
	for i := 0; i < 255; i++ {
		for j := 0; j < 255; j++ {
			cmp[i][j] = byteCmp(byte(i), byte(j))
		}
	}
}

// byteCmp compares two bytes lexicographically as bit vectors.
func byteCmp(a, b byte) int {
	for k := uint(0); k < 8; k++ {
		inA := (a & (1 << k)) != 0
		inB := (b & (1 << k)) != 0
		if inA != inB {
			if inA {
				return -1
			}
			return 1
		}
	}
	return 0
}

//
//--- Bits32 ---//
//

// Bits32 implements a simple 32-bit bitvector using a single uint32.
// Assuming CDim <= 3, this is big enough to represent any point set in a single cell.
type Bits32 uint32

// NewBits32 creates a new Bits32 with the given indices set.
func NewBits32(indices []int) Bits32 {
	var bits Bits32
	for _, index := range indices {
		bits = bits.Set(index)
	}
	return bits
}

// check checks that 0 <= i < 32.
func (b Bits32) check(i int) {
	if i < 0 || i >= 32 {
		panic("index out of bounds")
	}
}

// Set bit i to 1.
func (b Bits32) Set(i int) Bits32 {
	// b.check((i))
	return b | (1 << uint(i))
}

// Clear bit i to 0.
func (b Bits32) Clear(i int) Bits32 {
	// b.check((i))
	return b &^ (1 << uint(i))
}

// Test whether bit i is set.
func (b Bits32) Test(i int) bool {
	// b.check((i))
	return (b & (1 << uint(i))) != 0
}

// Empty returns true iff b represents the empty set (all bits are empty).
func (b Bits32) Empty() bool {
	return b == 0
}

// Cmp compares two Bits32s lexicographically as bit vectors.
func (b Bits32) Cmp(other Bits32) int {
	c := cmp[byte(b)][byte(other)]
	if c != 0 {
		return c
	}
	c = cmp[byte(b>>8)][byte(other>>8)]
	if c != 0 {
		return c
	}
	c = cmp[byte(b>>16)][byte(other>>16)]
	if c != 0 {
		return c
	}
	c = cmp[byte(b>>24)][byte(other>>24)]
	if c != 0 {
		return c
	}
	return 0
}

// Less returns true iff b < other as bit vectors.
func (b Bits32) Less(other Bits32) bool {
	return b.Cmp(other) < 0
}

// Apply applies a permutation to bits, and returns the resulting image.
func (b Bits32) Apply(perm []int) Bits32 {
	var image Bits32
	for i := 0; i < 32; i++ {
		if b.Test(i) {
			image = image.Set(perm[i])
		}
	}
	return image
}

// Intersection intersects bits with a set of indices.
func (b Bits32) Intersection(indices []int) Bits32 {
	var inter Bits32
	for _, i := range indices {
		if i < 32 && b.Test(i) {
			inter = inter.Set(i)
		}
	}
	return inter
}

// Inv inverts bits in a cSpace.
// TODO: Can we use a table independent of the space?
func (b Bits32) Inv(cSpace *space.Space) Bits32 {
	var inv Bits32
	for i := 0; i < 32; i++ {
		if b.Test(i) {
			inv = inv.Set(cSpace.Inv[i])
		}
	}
	return inv
}

// IsPreserved returns true iff a (cell) permutation preserves bits.
func (b Bits32) IsPreserved(perm []int) bool {
	for i := 0; i < 32; i++ {
		if b.Test(i) && !b.Test(perm[i]) {
			// perm maps i outside of bits.
			return false
		}
	}

	// perm maps bits into bits (and is a bijection), hence preserves bits as a set.
	return true
}

// Maximum returns the index of the largest set bit, or -1 if empty.
func (b Bits32) Maximum() int {
	for i := 31; i >= 0; i-- {
		if b.Test(i) {
			return i
		}
	}
	return -1
}

// PopCount returns population count of b (number of bits set).
func (b Bits32) PopCount() int {
	return mathutil.PopCountUint32(uint32(b))
}

// Eliminated returns the pts eliminated by a & b.
func (c Cells) Eliminated(a Bits32, b Bits32) Bits32 {
	cellSize := len(c.CSpace.Pts)
	var bits Bits32

	for i := 0; i < cellSize; i++ {
		if a.Test(i) {
			elim := c.CSpace.Elim[i]
			for j := 0; j < cellSize; j++ {
				if b.Test(j) {
					bits = bits.Set(elim[j])
				}
			}
		}
	}
	return bits
}

//
//--- BitsVec ---//
//

// A BitsVec represents a point set in the entire space as a slice of Bits32,
// one for each cell.
type BitsVec []Bits32

// NewBitsVec creates a new BitsVec for a Cells.
func (c Cells) NewBitsVec() BitsVec {
	return make([]Bits32, len(c.Cells))
}

// Hash a BitsVec. This is the modular hashing algorithm used by Java's hashCode(),
// which is very simple and fast.
func (vec BitsVec) Hash() uint32 {
	var hash uint32 = 1
	for _, bits := range vec {
		hash = 31*hash + uint32(bits)
	}
	return hash
}

// Apply applies a (cell) perm to each cell in a BitsVec, using an out vector.
func (vec BitsVec) Apply(perm []int, out BitsVec) {
	if len(vec) != len(out) {
		panic("length mismatch")
	}

	for i, bits := range vec {
		out[i] = bits.Apply(perm)
	}
}

// Clear clears a BitsVec.
func (vec BitsVec) Clear() {
	for i := range vec {
		vec[i] = 0
	}
}

// Cmp compares two BitsVecs lexicographically.
func (vec BitsVec) Cmp(other BitsVec) int {
	if len(vec) != len(other) {
		panic("length mismatch")
	}

	for i, bits := range vec {
		if bits != other[i] {
			return bits.Cmp(other[i])
		}
	}
	return 0
}

// EliminatedInCell returns the pts eliminated by a BitsVec in a single cell.
func (vec BitsVec) EliminatedInCell(cells Cells, cell int) Bits32 {
	elimByCell := cells.QSpace.Elim[cell]
	var bits Bits32

	for _, p := range cells.QSpace.Pts {
		third := elimByCell[p]
		if p <= third {
			bitsP := vec[p]
			bitsThird := vec[third]
			if !bitsP.Empty() && !bitsThird.Empty() {
				bits |= cells.Eliminated(bitsP, bitsThird)
			}
		}
	}
	return bits
}

// IsPreserved returns true iff a (cell) permutation preserves all cells in a BitsVec.
func (vec BitsVec) IsPreserved(perm []int) bool {
	for _, bits := range vec {
		if !bits.IsPreserved(perm) {
			return false
		}
	}
	return true
}

// PermuteValues permutes cells in a BitsVec, using an out vector.
func (vec BitsVec) PermuteValues(qIsom []int, out BitsVec) {
	if len(vec) != len(out) {
		panic("length mismatch")
	}

	for i, bits := range vec {
		out[qIsom[i]] = bits
	}
}

// ToPoints converts a BitsVec to a Points.
func (vec BitsVec) ToPoints(c Cells) *space.Points {
	if len(vec) != len(c.Cells) {
		panic("length mismatch")
	}

	var pts []int
	for k, bits := range vec {
		offset := c.Cells[k][0]
		for i := 0; i < 32; i++ {
			if bits.Test(i) {
				pts = append(pts, offset+i)
			}
		}
	}
	return space.NewPoints(c.Space, pts)
}
