package space

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/willf/bitset"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestPlanes(t *testing.T) {
	assert := assert.New(t)

	space := New(2)
	planes := space.Planes

	pts := []int{3, 4, 5, 6, 7, 8}
	directions := bitset.New(uint(len(space.Directions)))
	for _, p := range pts {
		dir := space.PtToDirection[p]
		if dir != -1 {
			directions.Set(uint(dir))
		}
	}

	assert.Equal(planes.PlaneCountsString(directions), "[0 2] => [1 3]")
}
