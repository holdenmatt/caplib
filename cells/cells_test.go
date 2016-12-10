package cells

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/holdenmatt/caplib/space"
	"github.com/holdenmatt/util"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func ExampleCells() {
	cells := New(space.New(4), []int{4, 2, 2, 2, 2, 2, 2, 2, 2})

	fmt.Println("D:", cells.Space.D)
	fmt.Println("CDim:", cells.CSpace.D)
	fmt.Println("QDim:", cells.QSpace.D)
	fmt.Println("Cells:", cells.Cells)
	fmt.Println("Counts:", cells.Counts)

	// Output:
	// D: 4
	// CDim: 2
	// QDim: 2
	// Cells: [[0 1 2 3 4 5 6 7 8] [9 10 11 12 13 14 15 16 17] [18 19 20 21 22 23 24 25 26] [27 28 29 30 31 32 33 34 35] [36 37 38 39 40 41 42 43 44] [45 46 47 48 49 50 51 52 53] [54 55 56 57 58 59 60 61 62] [63 64 65 66 67 68 69 70 71] [72 73 74 75 76 77 78 79 80]]
	// Counts: [4 2 2 2 2 2 2 2 2]
}

func TestCells(t *testing.T) {
	s := space.New(4)
	assert.Panics(t, func() { New(s, []int{1}) })
	assert.Panics(t, func() { New(s, []int{1, 2, 3}) })
	assert.Panics(t, func() { New(s, util.Range(81)) })
}

func ExampleProjCells() {
	cells := New(space.New(4), []int{4, 2, 2, 2, 2, 2, 2, 2, 2})

	proj := cells.ProjCells
	fmt.Println(proj.Indices)
	fmt.Println(proj.Counts)
	fmt.Println(proj.Sizes)

	// Output:
	// [0 1 3 4 5]
	// [4 2 2 2 2]
	// [4 8 12 16 20]
}
