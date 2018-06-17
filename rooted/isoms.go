package rooted

import (
	"fmt"

	"github.com/deckarep/golang-set"
	"github.com/holdenmatt/caplib/cells"
	"github.com/holdenmatt/util"
)

// Isoms represents all CIsoms and QIsoms (see cells/isoms.go) for a Cells with a fixed root.
// We restrict to only 1) CIsoms that minimize the root, and 2) QIsoms that preserve
// counts (precomputed for each depth of the search tree).
type Isoms struct {
	CIsoms   *cells.CellPerms // CIsoms that minimize the root.
	QIsoms   []util.Perms     // QIsoms that preserve counts, for each depth.
	QBases   [][][]int        // For each qIsom, pre-compute the basis that gets mapped to the std basis.
	QNzBases [][][]int        // For each qIsom, pre-compute the basis that gets mapped to the nonzero basis.
}

// newIsoms creates a new Isoms for a Rooted cell space.
func newIsoms(rooted Rooted) *Isoms {
	c := rooted.Cells
	cIsoms := cIsomsMinimizingRoot(rooted)

	var qIsoms []util.Perms
	countsPrefix := make([]int, c.Len())
	for _, cell := range c.ProjCells.Indices {
		invCell := c.QSpace.Inv[cell]
		count := c.Counts[cell]
		countsPrefix[cell] = count
		countsPrefix[invCell] = count

		qIsoms = append(qIsoms, qIsomsFixingCounts(c, countsPrefix))
	}

	qBases := make([][][]int, len(qIsoms))
	qNzBases := make([][][]int, len(qIsoms))
	for depth, qIs := range qIsoms {
		qBases[depth] = make([][]int, len(qIs.Perms))
		qNzBases[depth] = make([][]int, len(qIs.Perms))
		for i, qIsom := range qIs.Perms {
			qInv := util.InversePerm(qIsom)
			qBases[depth][i] = util.GetIndices(qInv, c.QSpace.StdBasis)
			qNzBases[depth][i] = util.GetIndices(qInv, c.NonzeroBasis)
		}
	}

	return &Isoms{cIsoms, qIsoms, qBases, qNzBases}
}

// cIsomsMinimizingRoot returns the subset of CIsoms that minimize the root.
func cIsomsMinimizingRoot(rooted Rooted) *cells.CellPerms {
	var perms [][]int
	var image1, image2, minImage cells.Bits32

	c := rooted.Cells
	root := rooted.Root

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

// qIsomsFixingCounts returns the subset of QIsoms that preserve the given counts,
// up to uniqueness on the nonzero counts.
func qIsomsFixingCounts(c cells.Cells, counts []int) util.Perms {
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
