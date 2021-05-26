package rest

import (
	"encoding/json"
	"net/http"

	"github.com/Akshit8/yt-search/elasticsearch"
	"github.com/Akshit8/yt-search/entity"
	"github.com/gorilla/mux"
)

// VideoHandler defines all handlers for resource Video.
type VideoHandler struct {
	vs *elasticsearch.VideoSearch
}

// NewVideoHandler creates new instance of VideoHandler.
func NewVideoHandler(vs *elasticsearch.VideoSearch) *VideoHandler {
	return &VideoHandler{vs}
}

// Register connects the handlers to the router.
func (v *VideoHandler) Register(r *mux.Router) {
	r.HandleFunc("/videos", v.get).Methods(http.MethodGet)
	r.HandleFunc("/search/videos", v.search).Methods(http.MethodPost)
}

func (v *VideoHandler) get(w http.ResponseWriter, r *http.Request) {
	renderResponse(w, "get working", 201)
}

// SearchVideosRequest defines the request used for searching videos.
type SearchVideosRequest struct {
	Title       *string  `json:"title"`
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

	res, err := v.vs.Search(r.Context(), req.Title, req.Description)
	if err != nil {
		renderErrorResponse(r.Context(), w, "search failed", err)
		return 
	}

	renderResponse(w, &SearchVideosResponse{Videos: res}, http.StatusOK)
}
