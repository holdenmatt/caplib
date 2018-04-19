package cells

// Rooted represents a Cells in which we have selected points in the origin cell.
// We generate caps in a single cell compatible with the root.
type Rooted struct {
	Cells
	Root Bits32
}

// CellCaps returns all caps of given size in a single cell (CSpace) that
// a) avoid some eliminated bits, and b) are compatible with the root.
func (r *Rooted) CellCaps(size int, elim Bits32, out []Bits32) []Bits32 {
	var empty Bits32
	return r.extendBits(empty, size, elim, out)
}

// extendBits extends a given starting cap, generating all possible bit combinations for
// larger bit indices, while avoiding elim bits.
func (r *Rooted) extendBits(bits Bits32, size int, elim Bits32, out []Bits32) []Bits32 {
	if bits.PopCount() == size {
		out = append(out, bits)
		return out
	}

	// Compute pts eliminated by bits with itself, and [root, -bits].
	el := elim
	el |= r.EliminatedFast(bits, bits)
	el |= r.EliminatedFast(r.Root, bits.Inv(r.CSpace))

	// Only change larger indices.
	cellSize := len(r.CSpace.Pts)
	nextBit := bits.Maximum() + 1
	for i := nextBit; i < cellSize; i++ {
		if !el.Test(i) {
			// Include bit i.
			nextBits := bits.Set(i)
			out = r.extendBits(nextBits, size, elim, out)
		}
	}
	return out
}
