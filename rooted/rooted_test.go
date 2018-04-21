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
	root := cells.Bits32(78)
	rooted := New(c, root)

	var out []cells.Bits32

	vec := cells.BitsVec{78, 0, 0, 0, 0, 0, 0, 0, 0}
	fmt.Println(rooted.CellCaps(vec, 1, out))

	vec = cells.BitsVec{78, 17, 257, 0, 0, 0, 0, 0, 0}
	fmt.Println(rooted.CellCaps(vec, 3, out))

	vec = cells.BitsVec{78, 17, 257, 17, 0, 0, 257, 0, 0}
	fmt.Println(rooted.CellCaps(vec, 4, out))

	vec = cells.BitsVec{78, 17, 257, 17, 68, 0, 257, 0, 10}
	fmt.Println(rooted.CellCaps(vec, 5, out))

	// Output:
	// [17 33 129 257 10 34 66 258 12 20 68 132 136 264 80 272 96 160]
	// [17 33 129 257 10 34 66 258 12 20 68 132 136 264 80 272 96 160]
	// [10 34 66 12 68 132 136 96 160]
	// [160]
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

func Example_cIsomsMinimizingRoot() {
	c := cells.New(space.New(4), []int{4, 2, 2, 2, 2, 2, 2, 2, 2})
	root := cells.NewBits32([]int{1, 2, 3, 6})
	rooted := New(c, root)

	cIsoms := cIsomsMinimizingRoot(rooted)
	fmt.Println(cIsoms.Perms)

	// Output:
	// [[0 1 2 3 4 5 6 7 8] [0 2 1 3 5 4 6 8 7] [0 1 2 6 7 8 3 4 5] [0 2 1 6 8 7 3 5 4] [0 3 6 1 4 7 2 5 8] [0 3 6 2 5 8 1 4 7] [0 6 3 1 7 4 2 8 5] [0 6 3 2 8 5 1 7 4]]
}

func Example_qIsomsFixingCounts() {
	c := cells.New(space.New(4), []int{4, 2, 2, 2, 2, 2, 2, 2, 2})
	qIsoms := qIsomsFixingCounts(c, []int{4, 2, 2, 2, 0, 0, 2, 0, 0})
	fmt.Println(qIsoms.Perms)

	// Output:
	// [[0 1 2 3 4 5 6 7 8] [0 1 2 6 7 8 3 4 5] [0 2 1 3 5 4 6 8 7] [0 2 1 6 8 7 3 5 4] [0 3 6 1 4 7 2 5 8] [0 3 6 2 5 8 1 4 7] [0 6 3 1 7 4 2 8 5] [0 6 3 2 8 5 1 7 4]]
}
