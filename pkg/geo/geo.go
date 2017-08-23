package geo

import (
	"math"

	"github.com/bobbytables/farandnear/pkg/quadtree"
)

// Semi-axes of WGS-84 geoidal reference
const (
	WGS84A = float64(6378137.0)
	WGS84B = float64(6356752.3)
)

// BoundingBoxFromCoords creates a bounding box from a center point of longitude
// and latitude in degrees.
// Basically all of this  was converted from this stackoverflow answer:
// https://stackoverflow.com/a/238558
func BoundingBoxFromCoords(latD, longD float64, halfSideKm float64) quadtree.AABB {
	lat := DegreesToRadians(latD)
	lon := DegreesToRadians(longD)
	halfSide := 1000 * halfSideKm

	// Radius of Earth at given latitude
	radius := WGS84EarthRadius(lat)
	// Radius of the parallel at given latitude
	pradius := radius * math.Cos(lat)

	latMin := RadiansToDegrees(lat - halfSide/radius)
	latMax := RadiansToDegrees(lat + halfSide/radius)
	lonMin := RadiansToDegrees(lon - halfSide/pradius)
	lonMax := RadiansToDegrees(lon + halfSide/pradius)

	return quadtree.NewAABB(latMin, lonMin, latMax, lonMax)
}

// DegreesToRadians converts degrees to radians
func DegreesToRadians(deg float64) float64 {
	return (math.Pi * deg) / 180
}

// RadiansToDegrees converts radians to degrees
func RadiansToDegrees(rad float64) float64 {
	return (rad * 180) / math.Pi
}

// WGS84EarthRadius returns Earth's radius at a given latitude,
// according to the WGS-84 ellipsoid [m]
func WGS84EarthRadius(lat float64) float64 {
	// http://en.wikipedia.org/wiki/Earth_radius
	An := WGS84A * WGS84A * math.Cos(lat)
	Bn := WGS84B * WGS84B * math.Sin(lat)
	Ad := WGS84A * math.Cos(lat)
	Bd := WGS84B * math.Sin(lat)

	return math.Sqrt((An*An + Bn*Bn) / (Ad*Ad + Bd*Bd))
}
