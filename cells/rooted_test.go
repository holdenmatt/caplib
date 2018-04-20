package cells

import (
	"fmt"

	"github.com/holdenmatt/caplib/space"
)

func ExampleRooted_CellCaps() {
	cells := New(space.New(4), []int{4, 2, 2, 2, 2, 2, 2, 2, 2})
	rooted := Rooted{
		Cells: cells,
		Root:  78,
	}

	var elim Bits32
	var out []Bits32
	fmt.Println(rooted.CellCaps(1, elim, out))
	fmt.Println(rooted.CellCaps(2, elim, out))

	// Output:
	// [1 2 4 8 16 32 64 128 256]
	// [17 33 129 257 10 34 66 258 12 20 68 132 136 264 80 272 96 160]
}

func ExampleCells_MinRoots() {
	c := New(space.New(4), []int{4, 2, 2, 2, 2, 2, 2, 2, 2})

	empty := Bits32(0)
	caps := c.addPairs([]Bits32{empty}, 2)
	fmt.Println(caps)

	roots := c.MinRoots()
	fmt.Println(len(roots))
	fmt.Println(roots[0].Root)

	// Output:
	// [78 278 166 344 232 432]
	// 1
	// 78
}
