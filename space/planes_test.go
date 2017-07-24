package space

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestPlanes(t *testing.T) {
	assert := assert.New(t)

	space := New(2)
	planes := space.Planes

	assert.Equal(planes.perp[0][1], true)
	assert.Equal(planes.perp[2][2], false)
	assert.Equal(planes.perp[1][3], true)

	pts := NewPoints(space, []int{3, 4, 5, 6, 7, 8})

	assert.Equal(planes.planeCount(pts, 1), 2)
	assert.Equal(planes.planeCount(pts, 3), 0)
	assert.Panics(func() {
		planes.planeCount(pts, 0)
	})

	assert.Equal(planes.PlaneCountsString(pts), "[0 2] => [1 3]")
}
