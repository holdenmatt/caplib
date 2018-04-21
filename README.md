# caplib

A Go library defining basic data structures used to search for large caps in ternary
affine spaces.

To keep code organized, this library defines the basic data structures, while a different
repo (capsearch) implements the search algorithm using this library.


## Package space

Define a ternary vector space, subsets of points, and isomorphism classes.

    modn.go:
        The ring ℤ/nℤ and its d-dimensional coordinate space.

    space.go:
        Space - The vector space (ℤ/3ℤ)^d.
        Points - A vector of sorted points in a Space.

    isoms.go:
        Define several isomorphism classes in a space: Translations,
        CoordPerms, CoordReflections, LinearIsoms.


## Package cells

Overlay a "cell" structure on a ternary space. This is one of the primary ways we constrain the highly exponential nature of the cap search problem, since it allows us to exploit
symmetries and confine our search to caps with defined cell counts.

    cells.go:
        Define Cells (a partition of a Space into translation cosets, with counts) and ProjCells
        (the projective version).

    isoms.go:
        Define several isomorphism classes that preserve the cell structure:
            CIsoms (isoms of the origin cell),
            QIsoms (isoms of the quotient, that permute cells), and
            Shears (cell translates in a given direction).

    bits.go:
        Define Bits32 (representing a set of points in a single cell), and
        BitsVec (representing a set of points in the entire space).
        We use byte tables and bit tricks to make operations more efficient.
