package space

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegersModN(t *testing.T) {
	Z5 := NewIntegersModN(5)
	assert.Equal(t, Z5.values, []int{0, 1, 2, 3, 4})
}

func TestCoordinatesModN(t *testing.T) {
	assert := assert.New(t)

	n := 3
	d := 2
	c := NewCoordinatesModN(n, d)

	assert.Equal(fmt.Sprint(c.Vecs), "[[0 0] [0 1] [0 2] [1 0] [1 1] [1 2] [2 0] [2 1] [2 2]]")

	for i := range []int{0, 1, 10, 50} {
		vec := c.Vecs[i]
		index := c.VecToIndex(vec)
		assert.Equal(index, i)
	}

	assert.Panics(func() {
		c.VecToIndex([]int{1})
	})
	assert.Panics(func() {
		c.VecToIndex([]int{-1, -1})
	})

	assert.Equal(c.Inv(0), 0)
	assert.Equal(c.Inv(1), 2)
	assert.Equal(c.Inv(2), 1)

	assert.Equal(c.Sum(0, 5), 5)
	assert.Equal(c.Sum(1, 2), 0)
	assert.Equal(c.Sum(2, 4), 3)

	assert.Equal(c.Dot(0, 1), 0)
	assert.Equal(c.Dot(1, 2), 2)
	assert.Equal(c.Dot(1, 3), 0)

	assert.Equal(c.Coords(), []int{0, 1})
	assert.Equal(c.StdBasis(), []int{1, 3})
	assert.Equal(c.Directions(), []int{1, 3, 4, 5})
}
