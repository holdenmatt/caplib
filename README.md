# caplib

A Go library defining basic data structures used to search for large caps in ternary
affine spaces.

To keep code organized, this library defines the basic data structures, while capsearch
implements the search algorithm using this library.


# Summary


## Package space

Define a ternary vector space, subsets of points, and isomorphism classes.

    modn.go:
        The ring ℤ/nℤ and its d-dimensional coordinate space.

    space.go:
        Space - The vector space (ℤ/3ℤ)^d.
        Points - A vector of sorted points in a Space (and ASCII-encoding).

    isoms.go:
        Define several isomorphism classes in a space: Translations,
        CoordPerms, CoordReflections, LinearIsoms.


## Package cells

Overlay a "cell" structure on a ternary space. This is one of the primary ways we tame the highly exponential nature of the cap search problem, since it allows us to exploit
symmetries and constrain a search to caps with defined cell counts.

    cells.go:
        Define Cells (a partition of a Space into cosets with counts) and
        ProjCells (the projective version).

    isoms.go:
        Define several isomorphism classes that preserve the cell structure:
        CIsoms, QIsoms, Shears.

    bits.go:
        Define Bits32 (representing a set of points in a single cell), and
        BitsVec (representing a set of points in the entire space).
