package cells

import (
	"fmt"
	"log"

	"github.com/deckarep/golang-set"
	"github.com/holdenmatt/util"
)

// CellPerms represents a class of permutations in CSpace.
// In order to quickly compute images of Bits32s, we precompute images for the 4 bytes.
type CellPerms struct {
	Perms        [][]int
	byteImage    [][4][256]Bits32  // Image of each byte in a Bits32, for each perm.
	zeroPreimage []int             // Preimage of 0 under each perm.
	minImage     map[Bits32]Bits32 // Map each Bits32 to its smallest image.
}

// NewCellPerms creates a new CellPerms for the given perms.
func (c Cells) NewCellPerms(perms [][]int) *CellPerms {
	cellSize := len(c.CSpace.Pts)

	// Create a mask to check that all set bits are < cellSize.
	var mask Bits32
	for i := 0; i < cellSize; i++ {
		mask |= (1 << uint(i))
	}

	var byteImage [][4][256]Bits32
	var zeroPreimage []int

	for _, perm := range perms {

		var image [4][256]Bits32
		for shift := uint(0); shift < 4; shift++ {
			for k := uint32(0); k < 256; k++ {
				preimage := Bits32(k << (8 * shift))
				if preimage&mask == preimage {
					image[shift][k] = preimage.Apply(perm)
				}
			}
		}
		byteImage = append(byteImage, image)

		zeroPreimage = append(zeroPreimage, perm[perm[0]])
	}

	minImage := make(map[Bits32]Bits32)
	return &CellPerms{perms, byteImage, zeroPreimage, minImage}
}

// Apply applies the kth cell perm to a Bits32.
func (p *CellPerms) Apply(k int, b Bits32) Bits32 {
	im := p.byteImage
	return im[k][0][byte(b)] |
		im[k][1][byte(b>>8)] |
		im[k][2][byte(b>>16)] |
		im[k][3][byte(b>>24)]
}

// ApplyVec applies the kth cell perm to every cell in a BitsVec, using an out vector.
func (p *CellPerms) ApplyVec(k int, vec BitsVec, out BitsVec) {
	if len(vec) != len(out) {
		panic("length mismatch")
	}

	for i, bits := range vec {
		out[i] = p.Apply(k, bits)
	}
}

// MinImageIndex applies a class of perms to a Bits32, and returns the smallest
// perm index for which the image is minimal.
func (p *CellPerms) MinImageIndex(b Bits32) int {
	index := 0
	minImage := b
	var image Bits32

	for k, preimage := range p.zeroPreimage {
		// Only consider perms whose image includes 0.
		if b.Test(preimage) {
			image = p.Apply(k, b)
			if image.Less(minImage) {
				index = k
				minImage = image
			}
		}
	}
	return index
}

// MinImage returns the smallest image of a Bits32 under a class of CellPerms.
func (p *CellPerms) MinImage(b Bits32) Bits32 {
	image, ok := p.minImage[b]
	if ok {
		return image
	}

	index := p.MinImageIndex(b)
	image = p.Apply(index, b)

	p.minImage[b] = image
	return image
}

// MinImages computes the MinImage for each bits in a BitsVec, using an out vector.
func (p *CellPerms) MinImages(vec BitsVec, out BitsVec) {
	if len(vec) != len(out) {
		panic("length mismatch")
	}

	for i, bits := range vec {
		out[i] = p.MinImage(bits)
	}
}

// We now define several classes of cell-preserving isomorphisms.
//
// Our goal is to represent the full class of linear isomorphisms that preserve cells (as a set).
// To that end, consider a point p in a cell-partitioned space, and write it as p = (x, y),
// where x represents coordinates of the cell (in q_space) and y represents the coordinates
// within the cell (in c_space).
//
// Every linear isomorphism M can then be represented as an invertible block matrix
//    |A B|
//    |C D|,
// partitioned by the dimensions of q_space/c_space.
//
// Assume M preserves the set of cells.
//
// Claim: B = 0
// ============
// Proof: M maps the point (x, y) into the cell with coordinates Ax + By.
// The points in a single cell are obtained by fixing x and varying y.
// These must all map to a single cell with fixed coordinates Ax + By.
// Taking x = 0, we must have By = 0 for all y, hence B = 0.
// --//--
//
// Now, we can decompose an arbitrary cell-preserving linear isomorphism
// into 3 (invertible) parts:
//    |A 0| |I 0| |I  0|
//    |0 I|.|0 D|.|C' I|,
// where C' = D^{-1} C.
//
// We gives names to these 3 classes (which commute, up to reordering their elements):
//
// 1. QIsoms
// Represented by an arbitrary invertible matrix A in GL(QSpace), which permutes
// cells around rigidly, while keeping points fixed within each cell.
//
// 2. CIsoms
// Represented by an arbitrary invertible matrix D in GL(CSpace), which applies
// a single linear isomorphism to all cells in parallel.
//
// 3. Shears
// Represented by an arbitrary matrix C of size CDim x QDim, which maps the
// point (x, y) -> (x, Cx + y), aka a "vertical shear".
//
// Each column C_i represents a translate by C_i within cells as we move in the
// direction of the i'th (QSpace) coordinate.
//
// Notice that each column shear C_i is independent of the others, and we can obtain
// the full set of shears by composing the shears for each coordinate direction.
//

// GetCIsoms returns the class of all "cell isoms" for a Cells,
// i.e. all linear isoms in cell space.
func (c Cells) GetCIsoms() *util.PermsProduct {
	perms := c.CSpace.LinearIsoms()
	log.Println("# of CIsoms:", perms.Len())
	return &perms
}

// GetQIsoms returns the class of "quotient isoms" for a Cells,
// i.e. isoms that rigidly permute cells while preserving the cell counts.
func (c Cells) GetQIsoms() *util.Perms {
	perms := c.QSpace.LinearIsomsFixingCounts(c.Counts)
	log.Println("# of QIsoms:", perms.Len())
	return &perms
}

// CIsomsMinimizingRoot returns the subset of CIsoms that minimize the given root bits.
func (c Cells) CIsomsMinimizingRoot(root Bits32) *CellPerms {
	var perms [][]int
	var image1, image2, minImage Bits32

	for _, perm1 := range c.CIsoms.Perms1.Perms {
		image1 = root.Apply(perm1)
		for _, perm2 := range c.CIsoms.Perms2.Perms {
			image2 = image1.Apply(perm2)

			if image2.Less(minImage) {
				minImage = image2
				perms = nil
			}

			if image2 == minImage {
				perm := util.Compose(perm1, perm2)
				perms = append(perms, perm)
			}
		}
	}
	// log.Printf("%d of %d CIsoms preserve: %v", len(perms), c.CIsoms.Len(), root)
	return c.NewCellPerms(perms)
}

// QIsomsFixingCounts returns the subset of QIsoms that preserve the given counts,
// up to uniqueness on the nonzero counts.
func (c Cells) QIsomsFixingCounts(counts []int) util.Perms {
	nonzeroIndices := util.Nonzero(counts)
	seen := mapset.NewSet()

	var perms [][]int
	for _, perm := range c.QIsoms.Perms {
		if util.PreservesValues(perm, counts) {
			// Take unique perms on the nonzero indices (using its String as a Set key)
			nonzeroPerm := util.GetIndices(perm, nonzeroIndices)
			key := fmt.Sprintf("%v", nonzeroPerm)
			if !(seen.Contains(key)) {
				seen.Add(key)
				perms = append(perms, perm)
			}
		}
	}
	// log.Printf("%d of %d QIsoms preserve counts: %v", len(perms), c.QIsoms.Len(), counts)
	return util.NewPerms(perms)
}

//
//--- Shears ---//
//

// MinShear computes the minimal shear of a BitsVec, in place.
func (c Cells) MinShear(vec BitsVec) {
	basis := c.QSpace.StdBasis
	lastIndex := c.lastNonemptyBasisIndex(vec)

	for i, basisPt := range basis {
		if i <= lastIndex {
			bits := vec[basisPt]
			minTIndex := c.Translations.MinImageIndex(bits)
			c.iShear(vec, minTIndex, i)
		}
	}
}

// iShear applies a shear to a BitsVec in place: apply the translation with a given
// index along the ith coordinate direction.
func (c Cells) iShear(vec BitsVec, transIndex int, i int) {
	translations := c.Translations

	index := (c.QSpace.D - 1) - i // Why?
	for k, qVec := range c.QSpace.Vecs.Vecs {
		value := qVec[index]
		if value == 1 {
			vec[k] = translations.Apply(transIndex, vec[k])
		} else if value == 2 {
			vec[k] = translations.Apply(transIndex, vec[k])
			vec[k] = translations.Apply(transIndex, vec[k])
		}
	}
}

// lastNonemptyBasisIndex returns the last std basis index for which the cell is nonempty.
// We expect (and assert) that if any std basis cell is empty, then all later cells
// must be empty as well.
func (c Cells) lastNonemptyBasisIndex(vec BitsVec) int {
	basis := c.QSpace.StdBasis
	for i, basisPt := range basis {
		if vec[basisPt] == 0 {
			// Check all larger cells are also empty.
			for k := basisPt + 1; k < len(vec); k++ {
				if vec[k] != 0 {
					panic("if a std basis cell is empty, all larger pts must also be empty")
				}
			}
			lastNonempty := i - 1
			return lastNonempty
		}
	}
	return len(basis) - 1
}
