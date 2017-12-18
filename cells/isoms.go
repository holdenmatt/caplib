package cells

import (
	"log"

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
