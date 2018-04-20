package cells

import "log"

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

// MinRoots finds all symmetric root caps in the origin cell that are
// minimal in their isomorphism class (effective up to CDim = 3).
func (c Cells) MinRoots() []Rooted {
	if c.CSpace.D > 3 {
		panic("MinRoots only computable for CDim <= 3")
	}

	target := c.Counts[0]
	caps := c.rootCaps(target)
	log.Printf("# of roots in cell 0: %d", len(caps))

	var res []Rooted
	for _, root := range caps {
		if c.isMinRoot(root) {
			rooted := Rooted{c, root}
			res = append(res, rooted)
		}
	}
	log.Printf("# of unique roots: %d", len(res))

	return res
}

// rootCaps returns all symmetric caps of a given size in the origin cell.
func (c Cells) rootCaps(size int) []Bits32 {
	if (size % 2) != 0 {
		panic("size must be even")
	}

	empty := []Bits32{Bits32(0)}
	nPairs := size / 2
	return c.addPairs(empty, nPairs)
}

// addPairs adds nPairs pairs to the given caps, and returns all resulting caps.
func (c Cells) addPairs(caps []Bits32, nPairs int) []Bits32 {
	if nPairs == 0 {
		return caps
	}

	prevCaps := c.addPairs(caps, nPairs-1)
	directions := c.CSpace.Directions

	var nextCaps []Bits32
	for _, cap := range prevCaps {
		// Skip any eliminated pts; only extend by larger directions.
		elim := c.EliminatedFast(cap, cap)
		maxDir := cap.Intersection(directions).Maximum()

		for _, p := range directions {
			if p > maxDir && !elim.Test(p) {
				nextCap := cap.Set(p).Set(c.CSpace.Inv[p])
				nextCaps = append(nextCaps, nextCap)
			}
		}
	}
	return nextCaps
}

// isMinRoot returns true iff root is minimal in its isomorphism class.
func (c Cells) isMinRoot(root Bits32) bool {
	var im1, im2 Bits32

	for _, perm1 := range c.CIsoms.Perms1.Perms {
		im1 = root.Apply(perm1)

		for _, perm2 := range c.CIsoms.Perms2.Perms {
			im2 = im1.Apply(perm2)

			if im2.Less(root) {
				return false
			}
		}
	}
	return true
}
