package rooted

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/holdenmatt/caplib/cells"
	"github.com/holdenmatt/caplib/space"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func ExampleRooted_CellCaps() {
	c := cells.New(space.New(4), []int{4, 2, 2, 2, 2, 2, 2, 2, 2})
	rooted := Rooted{
		Cells: c,
		Root:  78,
	}

	var elim cells.Bits32
	var out []cells.Bits32
	fmt.Println(rooted.CellCaps(1, elim, out))
	fmt.Println(rooted.CellCaps(2, elim, out))

	// Output:
	// [1 2 4 8 16 32 64 128 256]
	// [17 33 129 257 10 34 66 258 12 20 68 132 136 264 80 272 96 160]
}

func ExampleMinRoots() {
	c := cells.New(space.New(4), []int{4, 2, 2, 2, 2, 2, 2, 2, 2})

	empty := cells.Bits32(0)
	caps := addPairs(c, []cells.Bits32{empty}, 2)
	fmt.Println(caps)

	roots := MinRoots(c)
	fmt.Println(len(roots))
	fmt.Println(roots[0].Root)

	// Output:
	// [78 278 166 344 232 432]
	// 1
	// 78
}
