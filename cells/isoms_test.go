package cells

import (
	"fmt"

	"github.com/holdenmatt/caplib/space"
)

func ExampleCellPerms_MinImage() {
	cells := New(space.New(4), []int{4, 2, 2, 2, 2, 2, 2, 2, 2})
	translations := cells.Translations
	bits := NewBits32([]int{4, 8})
	fmt.Println(translations.MinImageIndex(bits))
	fmt.Println(translations.MinImage(bits))

	bits = NewBits32([]int{0, 4})
	fmt.Println(bits)
	fmt.Println(translations.MinImageIndex(bits))
	fmt.Println(translations.MinImage(bits))

	// Output:
	// 8
	// 17
	// 17
	// 0
	// 17
}

func ExampleCells_GetCIsoms() {
	cells := New(space.New(4), []int{4, 2, 2, 2, 2, 2, 2, 2, 2})
	cIsoms := cells.GetCIsoms()
	fmt.Println(cIsoms.Len())
	fmt.Println(cIsoms.Perms1.Perms)
	fmt.Println(cIsoms.Perms2.Perms)

	// Output:
	// 48
	// [[0 1 2 3 4 5 6 7 8] [0 2 1 3 5 4 6 8 7] [0 1 2 6 7 8 3 4 5] [0 2 1 6 8 7 3 5 4] [0 3 6 1 4 7 2 5 8] [0 3 6 2 5 8 1 4 7] [0 6 3 1 7 4 2 8 5] [0 6 3 2 8 5 1 7 4]]
	// [[0 1 2 3 4 5 6 7 8] [0 1 2 4 5 3 8 6 7] [0 1 2 5 3 4 7 8 6] [0 3 6 4 7 1 8 2 5] [0 3 6 5 8 2 7 1 4] [0 4 8 5 6 1 7 2 3]]
}

func ExampleCells_GetQIsoms() {
	cells := New(space.New(4), []int{4, 2, 2, 2, 2, 2, 2, 2, 2})
	qIsoms := cells.GetQIsoms()
	fmt.Println(qIsoms.Len())

	// Output:
	// 48
}

func ExampleCells_CIsomsMinimizingRoot() {
	cells := New(space.New(4), []int{4, 2, 2, 2, 2, 2, 2, 2, 2})
	root := NewBits32([]int{1, 2, 3, 6})
	cIsoms := cells.CIsomsMinimizingRoot(root)
	fmt.Println(cIsoms.Perms)

	// Output:
	// [[0 1 2 3 4 5 6 7 8] [0 2 1 3 5 4 6 8 7] [0 1 2 6 7 8 3 4 5] [0 2 1 6 8 7 3 5 4] [0 3 6 1 4 7 2 5 8] [0 3 6 2 5 8 1 4 7] [0 6 3 1 7 4 2 8 5] [0 6 3 2 8 5 1 7 4]]
}

func ExampleCells_QIsomsFixingCounts() {
	cells := New(space.New(4), []int{4, 2, 2, 2, 2, 2, 2, 2, 2})
	qIsoms := cells.QIsomsFixingCounts([]int{4, 2, 2, 2, 0, 0, 2, 0, 0})
	fmt.Println(qIsoms.Perms)

	// Output:
	// [[0 1 2 3 4 5 6 7 8] [0 1 2 6 7 8 3 4 5] [0 2 1 3 5 4 6 8 7] [0 2 1 6 8 7 3 5 4] [0 3 6 1 4 7 2 5 8] [0 3 6 2 5 8 1 4 7] [0 6 3 1 7 4 2 8 5] [0 6 3 2 8 5 1 7 4]]
}
