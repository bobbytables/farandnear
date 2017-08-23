package api

// Location represents a location for a API response
type Location struct {
	Name      string
	Latitude  float64
	Longitude float64
}

// LocationResp is encoded as JSON with locations for a POI request
type LocationResp struct {
	Locations []Location
}

// AddLocationReq is the JSON body received to add locations to our index
type AddLocationReq struct {
	Latitude  float64
	Longitude float64
	Name      string
}
