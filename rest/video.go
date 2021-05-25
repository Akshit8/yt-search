package rest

import (
	"net/http"

	"github.com/gorilla/mux"
)

// VideoHandler defines all handlers for resource Video.
type VideoHandler struct {
}

// NewVideoHandler creates new instance of VideoHandler.
func NewVideoHandler() *VideoHandler {
	return &VideoHandler{}
}

// Register connects the handlers to the router.
func (v *VideoHandler) Register(r *mux.Router) {
	r.HandleFunc("/videos/{id}", v.get).Methods(http.MethodGet)
	r.HandleFunc("/search/videos", v.search).Methods(http.MethodGet)
}

func (v *VideoHandler) get(w http.ResponseWriter, r *http.Request) {
	renderResponse(w, "get working", 201)
}

func (v *VideoHandler) search(w http.ResponseWriter, r *http.Request) {
	renderResponse(w, "search working", 200)
}
