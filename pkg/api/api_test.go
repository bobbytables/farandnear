package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bobbytables/farandnear/pkg/quadtree"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_SearchingLocations(t *testing.T) {
	location := quadtree.NewLocation(quadtree.Point{X: 80, Y: 160}, []byte("hello world"))
	qt := quadtree.NewQuadtree(32, quadtree.WorldAABB)
	require.NoError(t, qt.AddLocation(location))

	s := &Server{index: qt}

	ts := httptest.NewServer(s.Mux())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/v1/poi?lat=80&long=160")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "status should be OK")

	var r LocationResp
	err = json.NewDecoder(resp.Body).Decode(&r)
	require.NoError(t, err)

	assert.Len(t, r.Locations, 1)
}

func Test_AddingLocations(t *testing.T) {
	qt := quadtree.NewQuadtree(32, quadtree.NewAABB(0, 0, 90, 180))
	s := &Server{index: qt}
	ts := httptest.NewServer(s.Mux())
	defer ts.Close()

	req := AddLocationReq{
		Latitude:  51.1750868,
		Longitude: 115.5767292,
		Name:      "Banff Starbucks",
	}

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(req)
	require.NoError(t, err)

	resp, err := http.Post(ts.URL+"/api/v1/locations", "application/json", buf)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	assert.Len(t, qt.Locations, 1)
}
