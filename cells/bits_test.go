package cells

import (
	"fmt"
	"testing"

	"github.com/holdenmatt/caplib/space"
	"github.com/stretchr/testify/assert"
)

func ExampleBits32() {
	bits := NewBits32([]int{0, 1})
	fmt.Println(bits)

	bits = bits.Set(2)
	fmt.Println(bits)
	fmt.Println(bits.Test(0))

	bits = bits.Clear(0)
	fmt.Println(bits)
	fmt.Println(bits.Test(0))

	fmt.Println(bits.Empty())
	bits = bits.Clear(1).Clear(2)
	fmt.Println(bits.Empty())

	// Output:
	// 3
	// 7
	// true
	// 6
	// false
	// false
	// true
}

func ExampleBits32_Cmp() {
	bits1 := NewBits32([]int{0, 1})
	bits2 := NewBits32([]int{1, 2})
	fmt.Println(bits1.Less(bits2))
	fmt.Println(bits2.Less(bits1))
	fmt.Println(bits1.Cmp(bits1))
	fmt.Println(bits1.Cmp(bits2))

	// Output:
	// true
	// false
	// 0
	// -1
}

func ExampleBits32_Apply() {
	bits := NewBits32([]int{0, 1})
	image := bits.Apply([]int{1, 2, 0})
	fmt.Println(image)

	// Output:
	// 6
}

func ExampleBits32_Intersection() {
	bits := NewBits32([]int{0, 1, 2, 3})
	even := bits.Intersection([]int{0, 2, 4})
	odd := bits.Intersection([]int{1, 3, 5})
	fmt.Println(even)
	fmt.Println(odd)

	// Output:
	// 5
	// 10
}

func ExampleBits32_Inv() {
	cSpace := space.New(2)
	fmt.Println(NewBits32([]int{1}).Inv(cSpace))
	fmt.Println(NewBits32([]int{4}).Inv(cSpace))
	fmt.Println(NewBits32([]int{}).Inv(cSpace))

	// Output:
	// 4
	// 256
	// 0
}

func ExampleBits32_Maximum() {
	fmt.Println(NewBits32([]int{0, 1, 2, 3}).Maximum())
	fmt.Println(NewBits32([]int{31}).Maximum())
	fmt.Println(NewBits32([]int{0}).Maximum())
	fmt.Println(NewBits32([]int{}).Maximum())

	// Output:
	// 3
	// 31
	// 0
	// -1
}

func ExampleBits32_PopCount() {
	fmt.Println(NewBits32([]int{}).PopCount())
	fmt.Println(NewBits32([]int{1}).PopCount())
	fmt.Println(NewBits32([]int{0, 1, 2, 3}).PopCount())

	// Output:
	// 0
	// 1
	// 4
}

func TestEliminated(t *testing.T) {
	assert := assert.New(t)

	cells := New(space.New(4), []int{4, 2, 2, 2, 2, 2, 2, 2, 2})

	bits1 := NewBits32([]int{1, 2})
	bits2 := NewBits32([]int{3, 6})
	expected := NewBits32([]int{4, 5, 7, 8})
	assert.Equal(expected, cells.Eliminated(bits1, bits2))
}

func ExampleBitsVec() {
	cells := New(space.New(4), []int{4, 2, 2, 2, 2, 2, 2, 2, 2})
	vec := cells.NewBitsVec()
	vec[0] = NewBits32([]int{3, 6})
	fmt.Println(vec)

	transpose := []int{0, 3, 6, 1, 4, 7, 2, 5, 8}
	out := cells.NewBitsVec()
	vec.Apply(transpose, out)
	fmt.Println(out)
	fmt.Println(out.Cmp(vec))
	fmt.Println(out.EliminatedInCell(cells, 0))

	cells.Translations.MinImages(vec, out)
	fmt.Println(out)

	vec[1] = 1
	fmt.Println(vec)
	vec.PermuteValues(transpose, out)
	fmt.Println(out)

	vec.Clear()
	fmt.Println(vec)

	// Output:
	// [72 0 0 0 0 0 0 0 0]
	// [6 0 0 0 0 0 0 0 0]
	// -1
	// 7
	// [9 0 0 0 0 0 0 0 0]
	// [72 1 0 0 0 0 0 0 0]
	// [72 0 0 1 0 0 0 0 0]
	// [0 0 0 0 0 0 0 0 0]
}
