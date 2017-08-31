package cells

import (
	"github.com/cznic/mathutil"
	"github.com/holdenmatt/caplib/space"
)

// _cmp stores precomputed byte comparisons such that _cmp[i][j] = Cmp(i, j).
var _cmp [256][256]int

// _inv stores precomputed inverses of each byte in a Bits32.
var _inv [4][256]Bits32

// _elim stores precomputed Bits32 bits eliminated by various bytes.
var _elim [16][256][256]Bits32

func init() {
	for i := 0; i < 255; i++ {
		for j := 0; j < 255; j++ {
			_cmp[i][j] = byteCmp(byte(i), byte(j))
		}
	}

	// Invert in 3-space (cells can't be any bigger).
	cSpace := space.New(3)
	for pos := 0; pos < 4; pos++ {
		for j := 0; j < 255; j++ {

			jInv := Bits32(0)
			for k := 0; k < 8; k++ {
				// Is the k'th bit of j set?
				kthSet := (j & (1 << uint(k))) != 0
				if kthSet {
					// Invert the index, adjusted for byte position.
					index := 8*pos + k
					if index < cSpace.Size() {
						jInv = jInv.Set(cSpace.Inv[index])
					}
				}
			}
			_inv[pos][j] = jInv
		}
	}

	for i := 0; i < 255; i++ {
		for j := 0; j < 255; j++ {
			i1, i2, i3, i4 := Bits32(i), Bits32(i<<8), Bits32(i<<16), Bits32(i<<24)
			j1, j2, j3, j4 := Bits32(j), Bits32(j<<8), Bits32(j<<16), Bits32(j<<24)

			_elim[0][i][j] = elim(i1, j1)
			_elim[1][i][j] = elim(i1, j2)
			_elim[2][i][j] = elim(i1, j3)
			_elim[3][i][j] = elim(i1, j4)
			_elim[4][i][j] = elim(i2, j1)
			_elim[5][i][j] = elim(i2, j2)
			_elim[6][i][j] = elim(i2, j3)
			_elim[7][i][j] = elim(i2, j4)
			_elim[8][i][j] = elim(i3, j1)
			_elim[9][i][j] = elim(i3, j2)
			_elim[10][i][j] = elim(i3, j3)
			_elim[11][i][j] = elim(i3, j4)
			_elim[12][i][j] = elim(i4, j1)
			_elim[13][i][j] = elim(i4, j2)
			_elim[14][i][j] = elim(i4, j3)
			_elim[15][i][j] = elim(i4, j4)
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

func elim(a, b Bits32) Bits32 {
	cSpace := space.New(3)
	cellSize := len(cSpace.Pts)

	var bits Bits32
	for i := 0; i < cellSize; i++ {
		if a.Test(i) {
			elim := cSpace.Elim[i]
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
	c := _cmp[byte(b)][byte(other)]
	if c != 0 {
		return c
	}
	c = _cmp[byte(b>>8)][byte(other>>8)]
	if c != 0 {
		return c
	}
	c = _cmp[byte(b>>16)][byte(other>>16)]
	if c != 0 {
		return c
	}
	c = _cmp[byte(b>>24)][byte(other>>24)]
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

// Inv inverts bits in a CSpace.
func (b Bits32) Inv(cSpace *space.Space) Bits32 {
	var inv Bits32
	for i := 0; i < 32; i++ {
		if b.Test(i) {
			inv = inv.Set(cSpace.Inv[i])
		}
	}
	return inv
}

// InvFast inverts bits in CSpace.
func (b Bits32) InvFast() Bits32 {
	return (_inv[0][byte(b)] |
		_inv[1][byte(b>>8)] |
		_inv[2][byte(b>>16)] |
		_inv[3][byte(b>>24)])
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
	return c.EliminatedFast(a, b)

	/*
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
	*/
}

// EliminatedFast returns the pts eliminated by a & b.
func (c Cells) EliminatedFast(a Bits32, b Bits32) Bits32 {
	a1, a2, a3, a4 := byte(a), byte(a>>8), byte(a>>16), byte(a>>24)
	b1, b2, b3, b4 := byte(b), byte(b>>8), byte(b>>16), byte(b>>24)

	return (_elim[0][a1][b1] |
		_elim[1][a1][b2] |
		_elim[2][a1][b3] |
		_elim[3][a1][b4] |
		_elim[4][a2][b1] |
		_elim[5][a2][b2] |
		_elim[6][a2][b3] |
		_elim[7][a2][b4] |
		_elim[8][a3][b1] |
		_elim[9][a3][b2] |
		_elim[10][a3][b3] |
		_elim[11][a3][b4] |
		_elim[12][a4][b1] |
		_elim[13][a4][b2] |
		_elim[14][a4][b3] |
		_elim[15][a4][b4])
}

//
//--- BitsVec ---//
//

// A BitsVec represents a point set in the entire space as a slice of Bits32,
// one for each cell.
type BitsVec []Bits32

// NewBitsVec creates a new BitsVec with given length.
func NewBitsVec(l int) BitsVec {
	return make([]Bits32, l)
}

// NewBitsVec creates a new BitsVec for a Cells.
func (c Cells) NewBitsVec() BitsVec {
	return NewBitsVec(len(c.Cells))
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

// CopyTo copies this vec to other.
func (vec BitsVec) CopyTo(other BitsVec) {
	if len(vec) != len(other) {
		panic("length mismatch")
	}

	for i, val := range vec {
		other[i] = val
	}
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

// Equals returns true iff vec == other.
func (vec BitsVec) Equals(other BitsVec) bool {
	return vec.Cmp(other) == 0
}

// GetIndices computes vec[indices], in an out BitsVec.
func (vec BitsVec) GetIndices(indices []int, out BitsVec) {
	if len(indices) != len(out) {
		panic("length mismatch")
	}

	for i, index := range indices {
		out[i] = vec[index]
	}
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
