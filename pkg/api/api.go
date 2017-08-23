package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/bobbytables/farandnear/pkg/geo"
	"github.com/bobbytables/farandnear/pkg/quadtree"
)

// Server contains all of the logic necessary for searching for locations
type Server struct {
	index *quadtree.Quadtree
}

// NewServer constructs a new server with the given quadtree
func NewServer(index *quadtree.Quadtree) *Server {
	return &Server{index: index}
}

// Mux returns a http.Handler
func (s *Server) Mux() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/poi", decorateRequest(s.SearchLocations))
	r.HandleFunc("/api/v1/locations", decorateRequest(s.AddLocations)).Methods("POST")

	return r
}

// SearchLocations looks for locations within the given storage
func (s *Server) SearchLocations(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	latS, longS := q.Get("lat"), q.Get("long")
	lat, err := strconv.ParseFloat(latS, 64)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("invalid latitude passed"))
		return
	}

	long, err := strconv.ParseFloat(longS, 64)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("invalid longitude passed"))
		return
	}

	aabb := geo.BoundingBoxFromCoords(lat, long, 10)

	var locations []Location
	for _, loc := range s.index.FindLocations(aabb) {
		newLoc := Location{Latitude: loc.Point.X, Longitude: loc.Point.Y, Name: string(loc.Data)}
		locations = append(locations, newLoc)
	}

	resp := LocationResp{Locations: locations}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - internal server error"))
	}
}

func (s *Server) AddLocations(w http.ResponseWriter, req *http.Request) {
	var add AddLocationReq
	if err := json.NewDecoder(req.Body).Decode(&add); err != nil {
		logrus.WithError(err).Error("could not decode request body")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not decode request body"))
		return
	}

	l := quadtree.NewLocation(quadtree.Point{X: add.Latitude, Y: add.Longitude}, []byte(add.Name))
	logrus.WithFields(logrus.Fields{"latitude": add.Latitude, "longitude": add.Longitude}).Info("added point to index")

	if err := s.index.AddLocation(l); err != nil {
		logrus.WithError(err).Error("could not add location")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not add location"))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func decorateRequest(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		logrus.WithField("path", req.URL.Path).WithField("method", req.Method).Info("request received")

		handler(w, req)

		logrus.WithField("path", req.URL.Path).WithField("duration", time.Since(start)).Info("request finished")
	}
}
