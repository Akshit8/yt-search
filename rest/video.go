package rest

import (
	"encoding/json"
	"net/http"

	"github.com/Akshit8/yt-search/entity"
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
	r.HandleFunc("/videos", v.get).Methods(http.MethodGet)
	r.HandleFunc("/search/videos", v.search).Methods(http.MethodGet)
}

func (v *VideoHandler) get(w http.ResponseWriter, r *http.Request) {
	renderResponse(w, "get working", 201)
}

// SearchVideosRequest defines the request used for searching videos.
type SearchVideosRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
}

// SearchVideosResponse defines the response returned back after searching for any video.
type SearchVideosResponse struct {
	Videos []entity.Video `json:"videos"`
}

func (v *VideoHandler) search(w http.ResponseWriter, r *http.Request) {
	var req SearchVideosRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}

	defer r.Body.Close()
}
