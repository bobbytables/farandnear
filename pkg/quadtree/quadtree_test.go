package quadtree

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_QuadtreeAddLocation(t *testing.T) {
	t.Run("Adding a location to a quad tree", func(t *testing.T) {
		boundary := NewAABB(0, 0, 10, 10)
		qt := NewQuadtree(4, boundary)

		location := NewLocation(Point{2, 3}, []byte("somewhere"))
		qt.AddLocation(location)

		assert.Len(t, qt.Locations, 1)
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

		assert.Len(t, qt.Locations, 0)
		assert.Len(t, qt.Children, 4)
		assert.Len(t, qt.Children[0].Locations, 1, "north west count")
		assert.Len(t, qt.Children[1].Locations, 1, "north east count")
		assert.Len(t, qt.Children[2].Locations, 1, "south west count")
		assert.Len(t, qt.Children[3].Locations, 1, "south east count")

		t.Run("Adding keeps splitting leaf nodes", func(t *testing.T) {
			northWest := qt.Children[0]

			qt.AddLocation(NewLocation(Point{0, 1}, []byte{}))

			assert.Len(t, northWest.Locations, 0)
			assert.Len(t, northWest.Children[0].Locations, 1)
			assert.Len(t, northWest.Children[3].Locations, 1)
		})
	})
}

func Test_QuadtreeContains(t *testing.T) {
	boundary := NewAABB(0, 0, 10, 10)
	qt := NewQuadtree(4, boundary)

	assert.True(t, qt.CanFitLocation(NewLocation(Point{2, 3}, []byte{})))
	assert.True(t, qt.CanFitLocation(NewLocation(Point{0, 0}, []byte{})))
	assert.False(t, qt.CanFitLocation(NewLocation(Point{10, 11}, []byte{})))
	assert.False(t, qt.CanFitLocation(NewLocation(Point{12, 0}, []byte{})))
}

func Test_QuadtreeSubdivide(t *testing.T) {
	boundary := NewAABB(0, 0, 10, 10)
	qt := NewQuadtree(4, boundary)

	qt.subdivide()

	assert.Len(t, qt.Children, 4)
	assert.Equal(t, NewAABB(0, 0, 5, 5), qt.Children[0].Boundary, "north west")
	assert.Equal(t, NewAABB(5, 0, 5, 5), qt.Children[1].Boundary, "north east")
	assert.Equal(t, NewAABB(0, 5, 5, 5), qt.Children[2].Boundary, "south west")
	assert.Equal(t, NewAABB(5, 5, 5, 5), qt.Children[3].Boundary, "south east")
}

func Test_QuadtreeSearch(t *testing.T) {
	boundary := NewAABB(0, 0, 10, 10)
	qt := NewQuadtree(2, boundary)

	location := NewLocation(Point{3, 4}, []byte("waldo"))
	qt.AddLocation(location)
	farAwayLoc := NewLocation(Point{9, 9}, []byte("bobby"))
	qt.AddLocation(farAwayLoc)

	locations := qt.FindLocations(NewAABB(0, 0, 5, 5))

	require.Len(t, locations, 1)
	assert.Equal(t, locations[0].Data, []byte("waldo"))

	t.Run("Includes locations from a subdivided tree", func(t *testing.T) {
		farAwayLoc := NewLocation(Point{5, 5}, []byte("tables"))
		qt.AddLocation(farAwayLoc)

		locations := qt.FindLocations(NewAABB(0, 0, 5, 5))

		require.Len(t, locations, 2)
		assert.Equal(t, locations[0].Data, []byte("waldo"))
		assert.Equal(t, locations[1].Data, []byte("tables"))
	})
}

func Test_QuadtreeGeohash(t *testing.T) {
	boundary := NewAABB(0, 0, 10, 10)
	qt := NewQuadtree(4, boundary)

	assert.Equal(t, "s0gs3y", qt.Geohash(6))
}

func Test_QuadtreeEncoding(t *testing.T) {
	boundary := NewAABB(0, 0, 10, 10)
	qt := NewQuadtree(4, boundary)
	location := NewLocation(Point{3, 4}, []byte("waldo"))

	qt.AddLocation(location)

	b, err := qt.Encode()
	require.NoError(t, err)

	expected := `{"Capacity":4,"Boundary":{"Min":{"X":0,"Y":0},"Max":{"X":10,"Y":10}},"Children":null,"Locations":[{"Point":{"X":3,"Y":4},"Data":"d2FsZG8="}]}`
	assert.Equal(t, expected, string(b))
}

func Test_QuadtreeDecoding(t *testing.T) {
	encoded := `{"Capacity":4,"Boundary":{"Min":{"X":0,"Y":0},"Max":{"X":10,"Y":10}},"Children":null,"Locations":[{"Point":{"X":3,"Y":4},"Data":"d2FsZG8="}]}`

	qt, err := Decode([]byte(encoded))
	require.NoError(t, err)

	assert.Len(t, qt.Locations, 1)
	locations := qt.FindLocations(NewAABB(0, 0, 5, 5))
	assert.Len(t, locations, 1)
	assert.Equal(t, []byte("waldo"), locations[0].Data, "data matches after decoding")
}
