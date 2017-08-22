package quadtree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AABBIntersections(t *testing.T) {
	t.Run("2 bounding boxes that intersect", func(t *testing.T) {
		b1 := NewAABB(2, 2, 3, 3)
		b2 := NewAABB(0, 0, 3, 3)

		assert.True(t, b2.Intersects(b1))
		assert.True(t, b1.Intersects(b2))
	})

	t.Run("2 bounding boxes that do not intersect", func(t *testing.T) {
		b1 := NewAABB(2, 2, 3, 3)
		b2 := NewAABB(0, 0, 1, 1)

		assert.False(t, b2.Intersects(b1))
		assert.False(t, b1.Intersects(b2))
	})
}
