// Package cells overlays a cell structure on a ternary Space.
//
// Adding cells is one of the primary ways we tame the highly exponential nature
// of the cap search problem, since it allows us to exploit symmetries and
// constrain a search to caps with defined cell counts.
//
package cells

import (
	"github.com/holdenmatt/caplib/space"
	"github.com/holdenmatt/util"
)

// Cells partitions a Space's points into equal size "cells",
// which are translation cosets of a subspace of dimension k <= d.
//
// Each cell is isomorphic to a "cell space" ("CSpace") of dimension k.
// The collection of cells (cosets) form a quotient space ("QSpace") of dimension d-k.
//
// We assign a count to each cell, so we can prune search nodes that exceed any count.
type Cells struct {
	Space    *space.Space // Space to partition into cells
	Cells    [][]int      // Partition space.pts into disjoint cells
	Counts   []int        // Assign a fixed target count to each cell
	CellSize int          // Size of each cell.

	CSpace *space.Space // "Cell space" (isomorphic to cells[0])
	QSpace *space.Space // "Quotient space" (isomorphic to space / cSpace)

	ProjCells *ProjCells // The projective subset of cells

	Translations *CellPerms         // Translations of CSpace
	CIsoms       *util.PermsProduct // CSpace isoms
	QIsoms       *util.Perms        // QSpace isoms that preserve counts

	BitsVec BitsVec // Reusable buffer
}

// New creates a new Cells in a Space with the given counts.
func New(s *space.Space, counts []int) Cells {
	nCells := len(counts)
	qDim := util.Log(nCells, 3)
	cDim := s.D - qDim

	if cDim <= 0 || qDim <= 0 {
		panic("cDim and qDim must be > 0")
	}
	if cDim > 3 {
		panic("We assume throughout that cDim <= 3")
	}

	cSpace := space.New(cDim)
	qSpace := space.New(qDim)

	cellSize := cSpace.Size()
	cells := make([][]int, nCells)
	for i := range cells {
		cells[i] = util.Range(i*cellSize, (i+1)*cellSize)
	}

	// Check counts are non-negative and symmetric.
	for i, count := range counts {
		if count < 0 {
			panic("counts must be >= 0")
		}
		invCount := counts[qSpace.Inv[i]]
		if invCount != count {
			panic("counts must be symmetric")
		}
	}

	c := Cells{s, cells, counts, cellSize, cSpace, qSpace, nil, nil, nil, nil, nil}
	c.ProjCells = NewProjCells(c)
	c.Translations = c.NewCellPerms(cSpace.Translations().Perms)
	c.CIsoms = c.GetCIsoms()
	c.QIsoms = c.GetQIsoms()
	c.BitsVec = c.NewBitsVec()
	return c
}

// MinPt returns the min point in a given cell.
func (c *Cells) MinPt(cell int) int {
	return cell * c.CellSize
}

// MaxPt returns the max point in a given cell.
func (c *Cells) MaxPt(cell int) int {
	return (cell+1)*c.CellSize - 1
}

// ProjCells represents the "projective" subset of a Cells.
//
// The cells correspond to the origin and directions in QSpace, and are
// sufficient to define a symmetric point set.
type ProjCells struct {
	Indices []int // Projective cell indices
	Counts  []int // Corresponding counts
	Sizes   []int // Sum all counts up to each cell index (inclusive).
}

// NewProjCells creates a new ProjCells for a Cells.
func NewProjCells(c Cells) *ProjCells {
	indices := append([]int{space.ORIGIN}, c.QSpace.Directions...)
	counts := util.GetIndices(c.Counts, indices)

	sizes := make([]int, len(indices))
	sizes[0] = counts[0]
	for i := 1; i < len(sizes); i++ {
		// Non-origin cells count double to account for the inverse cell.
		sizes[i] = sizes[i-1] + 2*counts[i]
	}

	return &ProjCells{indices, counts, sizes}
}
