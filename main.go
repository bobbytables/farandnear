package main

import (
	"net/http"

	"github.com/bobbytables/farandnear/pkg/api"
	"github.com/bobbytables/farandnear/pkg/quadtree"
)

func main() {
	qt := quadtree.NewQuadtree(32, quadtree.NewAABB(-90, -180, 90, 180))
	server := api.NewServer(qt)
	http.ListenAndServe(":7000", server.Mux())
}
