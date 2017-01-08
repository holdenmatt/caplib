package space

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestSpaceCache(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(len(spaceCache), 0)
	New(1)
	assert.Equal(len(spaceCache), 1)
	New(1)
	New(1)
	assert.Equal(len(spaceCache), 1)
	New(2)
	assert.Equal(len(spaceCache), 2)
}

func TestSpace(t *testing.T) {
	assert := assert.New(t)

	space := New(2)
	assert.Equal(fmt.Sprintf("%v", space), "Space[d = 2]")

	assert.Equal(space.Pts, []int{0, 1, 2, 3, 4, 5, 6, 7, 8})
	assert.Equal(space.Inv, []int{0, 2, 1, 6, 8, 7, 3, 5, 4})

	assert.Equal(len(space.Sum), 9)
	assert.Equal(space.Sum[0], space.Pts)
	assert.Equal(space.Sum[1], []int{1, 2, 0, 4, 5, 3, 7, 8, 6})

	assert.Equal(len(space.Elim), 9)
	assert.Equal(space.Elim[0], space.Inv)
	assert.Equal(space.Elim[1], []int{2, 1, 0, 8, 7, 6, 5, 4, 3})

	assert.Equal(space.Perp[0][1], true)
	assert.Equal(space.Perp[2][2], false)
	assert.Equal(space.Perp[1][3], true)

	assert.Equal(space.StdBasis, []int{1, 3})
	assert.Equal(space.Directions, []int{1, 3, 4, 5})

	assert.Equal(space.Size(), 9)
}

func TestSpan(t *testing.T) {
	assert := assert.New(t)

	s := New(4)
	assert.Equal(s.Span([]int{0}), []int{0})
	assert.Equal(s.Span([]int{1, 2}), []int{0, 1, 2})
	assert.Equal(s.Span([]int{1, 3}), []int{0, 1, 2, 3, 4, 5, 6, 7, 8})
	assert.Equal(s.Span([]int{9, 27}), []int{0, 9, 18, 27, 36, 45, 54, 63, 72})
}

func TestPoints(t *testing.T) {
	assert := assert.New(t)

	space := New(2)
	pts := NewPoints(space, []int{3, 4, 5, 6, 7, 8})
	assert.Equal(fmt.Sprintf("%v", pts), "Points[3 4 5 6 7 8]")

	assert.Equal(pts.PlaneCount(1), 2)
	assert.Equal(pts.PlaneCount(3), 0)
	assert.Panics(func() {
		pts.PlaneCount(0)
	})
	assert.Equal(pts.PlaneCountsString(), "[0 2] => [1 3]")
}

func ExamplePoints_Encode() {
	space := New(1)
	pts := NewPoints(space, []int{})
	fmt.Println(pts.Encode())

	pts = NewPoints(space, []int{1, 2})
	fmt.Println(pts.Encode())

	space = New(2)
	pts = NewPoints(space, []int{0, 1, 2})
	fmt.Println(pts.Encode())

	space = New(3)
	pts = NewPoints(space, []int{0, 1, 2, 9, 18})
	fmt.Println(pts.Encode())

	// Output:
	// Points[]
	// -
	// -
	// -
	//
	// Points[1 2]
	// -
	// #
	// #
	//
	// Points[0 1 2]
	// # - -
	// # - -
	// # - -
	//
	// Points[0 1 2 9 18]
	// # - -
	// # - -
	// # - -
	//
	// # - -
	// - - -
	// - - -
	//
	// # - -
	// - - -
	// - - -
}
