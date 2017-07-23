package space

import (
	"log"

	"github.com/holdenmatt/util"
)

//
//--- Translations ---//
//

// Translations returns the group of all translations for a Space.
// Note: These are affine (not linear) isomorphisms.
func (s *Space) Translations() util.Perms {
	var perms [][]int
	for _, p1 := range s.Pts {
		var perm []int
		for _, p2 := range s.Pts {
			perm = append(perm, s.Sum[p1][p2])
		}
		perms = append(perms, perm)
	}

	return util.NewPerms(perms)
}

//
//--- Coord isoms ---//
//

// CoordPerms is the group generated by coordinate permutations in a Space.
func (s *Space) CoordPerms() util.Perms {
	var perms [][]int
	image := make([]int, s.D)

	coords := s.Vecs.Coords()
	for dPerm := range util.Permutations(coords, len(coords)) {
		var perm []int
		for _, vec := range s.Vecs.Vecs {
			for i, img := range dPerm {
				image[i] = vec[img]
			}
			pt := s.Vecs.VecToIndex(image)
			perm = append(perm, pt)
		}
		perms = append(perms, perm)
	}

	return util.NewPerms(perms)
}

// CoordReflections is the group generated by reflecting each coordinate in a Space.
func (s *Space) CoordReflections() util.Perms {
	var perms [][]int
	image := make([]int, s.D)

	// Iterate over all +1/-1 d-tuples (1,...,1) to (-1,...,-1).
	for signs := range util.Product([]int{1, -1}, s.D) {
		var perm []int
		for _, vec := range s.Vecs.Vecs {
			for i, sign := range signs {
				image[i] = (3 + vec[i]*sign) % 3
			}
			pt := s.Vecs.VecToIndex(image)
			perm = append(perm, pt)
		}
		perms = append(perms, perm)
	}

	return util.NewPerms(perms)
}

//
//--- Linear isoms ---//
//

// LinearIsoms returns all linear isoms of a space.
//
// Because this is a large group, we decompose it into a product of 2 classes:
// CoordPerms/CoordReflections and LinearIsomsModCoords.
// In d = 4, for example, there are ~24M linear isoms, which factors into
// classes of size 234 and 63180.
func (s *Space) LinearIsoms() util.PermsProduct {
	log.Printf("Computing linear isoms (D = %d)", s.D)
	perms1 := s.CoordPerms().Compose(s.CoordReflections())
	perms2 := s.LinearIsomsModCoords()
	perms := util.NewPermsProduct(perms1, perms2)
	log.Printf("...done (%d isoms)", perms.Len())
	return perms
}

// LinearIsomsModCoords returns a unique representative of every class of
// linear isoms of a space, modulo CoordIsoms.
func (s *Space) LinearIsomsModCoords() util.Perms {
	log.Printf("Computing linear isoms mod coords (D = %d)", s.D)

	bases := s.sorted1Bases([]int{})
	var perms [][]int
	for _, basis := range bases {
		perms = append(perms, s.basisToPerm(basis))
	}
	log.Printf("...done (%d isoms)", len(perms))
	return util.NewPerms(perms)
}

// LinearIsomsFixingCounts returns all isoms that preserve a vector of counts.
func (s *Space) LinearIsomsFixingCounts(counts []int) util.Perms {
	if len(counts) != len(s.Pts) {
		panic("length mismatch")
	}
	log.Printf("Computing linear isoms fixing counts (D = %d)", s.D)

	var perms [][]int
	bases := s.basesFixingCounts(counts, []int{})
	for _, basis := range bases {
		perms = append(perms, s.basisToPerm(basis))
	}
	log.Printf("...done (%d isoms)", len(perms))
	return util.NewPerms(perms)
}

//
//--- Bases (used to compute linear isoms) ---//
//

// basisToPerm returns the linear map taking the std basis to a given basis.
// This allows us to specify a linear isom by its image basis.
func (s *Space) basisToPerm(basis []int) []int {
	// ith pt -> ith coeff vector -> ith basis image
	images := s.Span(basis)

	if len(basis) != s.D {
		panic("basis must have dimension D")
	}
	if len(images) != len(s.Pts) {
		panic("basis is not linearly independent")
	}

	return images
}

// BasisToInvPerm returns the linear map taking a given basis to the std basis.
func (s *Space) BasisToInvPerm(basis []int) []int {
	images := s.basisToPerm(basis)
	return util.InversePerm(images)
}

// sorted1Bases returns all sorted bases consisting of only vectors with
// leading 1. This is a single representative for each basis class mod CoordIsoms.
func (s *Space) sorted1Bases(partialBasis []int) [][]int {
	if len(partialBasis) == s.D {
		return [][]int{partialBasis}
	}

	max := util.Maximum(partialBasis)
	partialSpan := s.Span(partialBasis)

	var bases [][]int
	for p := range s.Pts {
		// Only select linearly independent pts.
		if !util.Contains(partialSpan, p) {
			// Only select from larger pts, to produce sorted bases.
			if p > max {
				// Only select pts with leading 1s.
				if util.Contains(s.Directions, p) {
					nextPartialBasis := append(util.Clone(partialBasis), p)
					nextBases := s.sorted1Bases(nextPartialBasis)
					bases = append(bases, nextBases...)
				}
			}
		}
	}
	return bases
}

// basesFixingCounts returns all bases that preserve a vector of counts.
func (s *Space) basesFixingCounts(counts []int, partialBasis []int) [][]int {
	if len(counts) != len(s.Pts) {
		panic("length mismatch")
	}

	if len(partialBasis) == s.D {
		return [][]int{partialBasis}
	}

	// The partial map takes stdSpan -> partialSpan.
	partialSpan := s.Span(partialBasis)
	stdSpan := util.Range(len(partialSpan))

	// Find the pts/counts added by including the next StdBases vector.
	k := len(partialBasis)
	nextStdBasisPt := s.StdBasis[k]
	stdPts := util.GetIndices(s.Sum[nextStdBasisPt], stdSpan)
	stdCounts := util.GetIndices(counts, stdPts)

	// Each extension maps nextStdBasisPt -> p; find those that preserve counts.
	var bases [][]int
	nextPts := make([]int, len(partialSpan))
	nextCounts := make([]int, len(partialSpan))
	for p := range s.Pts {
		// Only select linearly independent pts.
		if !util.Contains(partialSpan, p) {
			// Check p itself preserves the count.
			if counts[p] == counts[nextStdBasisPt] {
				// Check all the added pts preserve the count.
				util.GetIndicesOut(s.Sum[p], partialSpan, nextPts)
				util.GetIndicesOut(counts, nextPts, nextCounts)
				if util.Equal(nextCounts, stdCounts) {
					nextPartialBasis := append(util.Clone(partialBasis), p)
					nextBases := s.basesFixingCounts(counts, nextPartialBasis)
					bases = append(bases, nextBases...)
				}
			}
		}
	}
	return bases
}
