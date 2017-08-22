package quadtree

// AABB represents an axis aligned bounding box
type AABB struct {
	Min, Max Point
}

// NewAABB constructs a new axis aligned bounding box
func NewAABB(x, y, width, height float64) AABB {
	return AABB{
		Min: Point{x, y},
		Max: Point{x + width, y + height},
	}
}

// Intersects checks to see if 2 axis aligned bounding boxes
// intersect eachother
func (a AABB) Intersects(b AABB) bool {
	return (a.Min.X <= b.Max.X && a.Max.X >= b.Min.X) &&
		(a.Min.Y <= b.Max.Y && a.Max.Y >= b.Min.Y)
}

func (a AABB) Width() float64 {
	return a.Max.X - a.Min.X
}

func (a AABB) Height() float64 {
	return a.Max.Y - a.Min.Y
}

func (a AABB) Contains(p Point) bool {
	return (p.X >= a.Min.X && p.X <= a.Max.X && p.Y >= a.Min.Y && p.Y <= a.Max.Y)
}
