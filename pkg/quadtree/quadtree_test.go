package quadtree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_QuadtreeAddLocation(t *testing.T) {
	t.Run("Adding a location to a quad tree", func(t *testing.T) {
		boundary := NewAABB(0, 0, 10, 10)
		qt := NewQuadtree(4, boundary)

		location := NewLocation(Point{2, 3}, []byte("somewhere"))
		qt.AddLocation(location)

		assert.Len(t, qt.locations, 1)
	})

	t.Run("Adding a location to the quadtree when the capacity is full subdivides it", func(t *testing.T) {
		boundary := NewAABB(0, 0, 10, 10)
		qt := NewQuadtree(1, boundary)

		var err error

		location := NewLocation(Point{3, 3}, []byte("north west"))
		qt.AddLocation(location)

		northEastLoc := NewLocation(Point{6, 3}, []byte("noth east"))
		err = qt.AddLocation(northEastLoc)
		assert.NoError(t, err, "no error adding the north east point")

		southWestLoc := NewLocation(Point{2, 6}, []byte("south west"))
		err = qt.AddLocation(southWestLoc)
		assert.NoError(t, err, "no error adding the south west point")

		southEastLoc := NewLocation(Point{6, 6}, []byte("south east"))
		err = qt.AddLocation(southEastLoc)
		assert.NoError(t, err, "no error adding the south east point")

		assert.Len(t, qt.locations, 0)
		assert.Len(t, qt.children, 4)
		assert.Len(t, qt.children[0].locations, 1, "north west count")
		assert.Len(t, qt.children[1].locations, 1, "north east count")
		assert.Len(t, qt.children[2].locations, 1, "south west count")
		assert.Len(t, qt.children[3].locations, 1, "south east count")

		t.Run("Adding keeps splitting leaf nodes", func(t *testing.T) {
			northWest := qt.children[0]

			qt.AddLocation(NewLocation(Point{0, 1}, []byte{}))

			assert.Len(t, northWest.locations, 0)
			assert.Len(t, northWest.children[0].locations, 1)
			assert.Len(t, northWest.children[3].locations, 1)
		})
	})
}

func Test_QuadtreeContains(t *testing.T) {
	boundary := NewAABB(0, 0, 10, 10)
	qt := NewQuadtree(4, boundary)

	assert.True(t, qt.ContainsLocation(NewLocation(Point{2, 3}, []byte{})))
	assert.True(t, qt.ContainsLocation(NewLocation(Point{0, 0}, []byte{})))
	assert.False(t, qt.ContainsLocation(NewLocation(Point{10, 11}, []byte{})))
	assert.False(t, qt.ContainsLocation(NewLocation(Point{12, 0}, []byte{})))
}

func Test_QuadtreeSubdivide(t *testing.T) {
	boundary := NewAABB(0, 0, 10, 10)
	qt := NewQuadtree(4, boundary)

	qt.subdivide()

	assert.Len(t, qt.children, 4)
	assert.Equal(t, NewAABB(0, 0, 5, 5), qt.children[0].boundary, "north west")
	assert.Equal(t, NewAABB(5, 0, 5, 5), qt.children[1].boundary, "north east")
	assert.Equal(t, NewAABB(0, 5, 5, 5), qt.children[2].boundary, "south west")
	assert.Equal(t, NewAABB(5, 5, 5, 5), qt.children[3].boundary, "south east")
}
