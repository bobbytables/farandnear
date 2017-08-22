package quadtree

import "errors"

var (
	// ErrOutOfBounds is returned when attempting to add a location to a quadtree
	// that is outside of its boundary
	ErrOutOfBounds = errors.New("quadtree: location out of bounds")

	// ErrCouldNotAdd is returned when an unknown error has occurred adding a
	// location to a quadtree
	ErrCouldNotAdd = errors.New("quadtree: could not add location")
)

// Location contains a point and arbitrary information in the Data field
type Location struct {
	Point Point
	Data  []byte
}

// NewLocation constructs a location
func NewLocation(p Point, data []byte) *Location {
	return &Location{p, data}
}

// Point represents a place within a quadtree boundary
type Point struct {
	X, Y float64
}

// Quadtree stores points within a bounding box
type Quadtree struct {
	capacity int
	boundary AABB

	children  []*Quadtree
	locations []*Location
}

// NewQuadtree constructs a quadtree with a capacity and boundary
func NewQuadtree(cap int, boundary AABB) *Quadtree {
	return &Quadtree{
		capacity: cap,
		boundary: boundary,
	}
}

// AddLocation adds a location to the quadtree or relevant nodes
func (q *Quadtree) AddLocation(l *Location) error {
	if !q.boundary.Contains(l.Point) {
		return ErrOutOfBounds
	}

	// if we're attempting to add to a subdivided quadtree we need to skip
	// to the portion where we add to the leaf nodes instead of this node
	if len(q.children) > 0 {
		goto ADDTOCHILD
	}

	// our capacity is still lower than locations inside of this so we can just
	// append
	if len(q.locations) < q.capacity {
		q.locations = append(q.locations, l)
		return nil
	}

	// if we hit here,  the locations are at capacity and we need
	// to subdivide (if we havent already)
	if len(q.children) == 0 {
		q.subdivide()
	}

ADDTOCHILD:
	for _, qt := range q.children {
		// fmt.Printf("%#v\n", qt.boundary)
		err := qt.AddLocation(l)
		switch err {
		case ErrOutOfBounds:
			continue
		case nil:
			return nil
		default:
			return err
		}
	}

	return ErrCouldNotAdd
}

// ContainsLocation checks if the location fits inside of this Quadtrees boundary
func (q *Quadtree) ContainsLocation(l *Location) bool {
	return q.boundary.Contains(l.Point)
}

func (q *Quadtree) subdivide() {
	b := q.boundary
	halfWidth, halfHeight := b.Width()/2, b.Height()/2
	northwest := NewQuadtree(q.capacity, NewAABB(b.Min.X, b.Min.Y, halfWidth, halfHeight))
	northeast := NewQuadtree(q.capacity, NewAABB(b.Min.X+halfWidth, b.Min.Y, halfWidth, halfHeight))
	southwest := NewQuadtree(q.capacity, NewAABB(b.Min.X, b.Min.Y+halfHeight, halfWidth, halfHeight))
	southeast := NewQuadtree(q.capacity, NewAABB(b.Min.X+halfWidth, b.Min.Y+halfHeight, halfWidth, halfHeight))

	q.children = []*Quadtree{northwest, northeast, southwest, southeast}

	// since we've subdivided, we need to rebalance all of our locations to our
	// new quadrants. Since we've modified this quadtree to have children,
	// we can just call each location with AddLocation on this instance
	// and it'll be divided naturally
	for _, l := range q.locations {
		q.AddLocation(l)
	}

	// We've distributed the locations so now we can wipe out this nodes locations slice
	q.locations = nil
}
