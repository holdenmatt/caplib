# caplib

A Go library defining basic data structures used to search for large caps in ternary
affine spaces.

To keep code organized, this library defines the data structures, while capsearch
implements the search algorithm using caplib.

# Summary

    modn.go - Defines the the ring ℤ/nℤ and its d-dimensional coordinate space.
    space.go - Defines the vector space (ℤ/3ℤ)^d.
