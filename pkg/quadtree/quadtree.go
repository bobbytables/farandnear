package quadtree

import (
	"encoding/json"
	"errors"

	"github.com/mmcloughlin/geohash"
)

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
	Capacity int
	Boundary AABB

	Children  []*Quadtree `json:,omitempty`
	Locations []*Location `json:,omitempty`
}

func Decode(b []byte) (*Quadtree, error) {
	q := new(Quadtree)

	if err := json.Unmarshal(b, q); err != nil {
		return nil, err
	}

	return q, nil
}

// NewQuadtree constructs a quadtree with a capacity and boundary
func NewQuadtree(cap int, boundary AABB) *Quadtree {
	return &Quadtree{
		Capacity: cap,
		Boundary: boundary,
	}
}

// AddLocation adds a location to the quadtree or relevant nodes
func (q *Quadtree) AddLocation(l *Location) error {
	if !q.Boundary.Contains(l.Point) {
		return ErrOutOfBounds
	}

	// if we're attempting to add to a subdivided quadtree we need to skip
	// to the portion where we add to the leaf nodes instead of this node
	if len(q.Children) > 0 {
		goto ADDTOCHILD
	}

	// our capacity is still lower than locations inside of this so we can just
	// append
	if len(q.Locations) < q.Capacity {
		q.Locations = append(q.Locations, l)
		return nil
	}

	// if we hit here,  the locations are at capacity and we need
	// to subdivide (if we havent already)
	if len(q.Children) == 0 {
		q.subdivide()
	}

ADDTOCHILD:
	for _, qt := range q.Children {
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

// CanFitLocation checks if the location fits inside of this Quadtrees boundary
func (q *Quadtree) CanFitLocation(l *Location) bool {
	return q.Boundary.Contains(l.Point)
}

// Geohash returns a geohash encoding of the center of the quadtrees boundary
func (q *Quadtree) Geohash(precision uint) string {
	b := q.Boundary
	halfWidth, halfHeight := b.Width()/2, b.Height()/2

	return geohash.EncodeWithPrecision(halfWidth, halfHeight, precision)
}

// FindLocations returns all locations within the quadtree in the given boundary
func (q *Quadtree) FindLocations(boundary AABB) []*Location {
	var locations []*Location

	for _, l := range q.Locations {
		if boundary.Contains(l.Point) {
			locations = append(locations, l)
		}
	}

	for _, child := range q.Children {
		locations = append(locations, child.FindLocations(boundary)...)
	}

	return locations
}

func (q *Quadtree) subdivide() {
	b := q.Boundary
	halfWidth, halfHeight := b.Width()/2, b.Height()/2
	northwest := NewQuadtree(q.Capacity, NewAABB(b.Min.X, b.Min.Y, halfWidth, halfHeight))
	northeast := NewQuadtree(q.Capacity, NewAABB(b.Min.X+halfWidth, b.Min.Y, halfWidth, halfHeight))
	southwest := NewQuadtree(q.Capacity, NewAABB(b.Min.X, b.Min.Y+halfHeight, halfWidth, halfHeight))
	southeast := NewQuadtree(q.Capacity, NewAABB(b.Min.X+halfWidth, b.Min.Y+halfHeight, halfWidth, halfHeight))

	q.Children = []*Quadtree{northwest, northeast, southwest, southeast}

	// since we've subdivided, we need to rebalance all of our locations to our
	// new quadrants. Since we've modified this quadtree to have children,
	// we can just call each location with AddLocation on this instance
	// and it'll be divided naturally
	for _, l := range q.Locations {
		q.AddLocation(l)
	}

	// We've distributed the locations so now we can wipe out this nodes locations slice
	q.Locations = nil
}

func (q *Quadtree) Encode() ([]byte, error) {
	return json.Marshal(q)
}
