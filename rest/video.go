package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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
	r.HandleFunc("/videos", v.get).Queries("skip", "{skip}").Queries("limit", "{limit}").Methods(http.MethodGet)
	r.HandleFunc("/search/videos", v.search).Methods(http.MethodPost)
}

// GetVideosResponse defines the response returned back after getting for any video.
type GetVideosResponse struct {
	Videos []entity.Video `json:"videos"`
}

// get returns all videos indexed in es as paginated list
func (v *VideoHandler) get(w http.ResponseWriter, r *http.Request) {
	skip := r.FormValue("skip")
	limit := r.FormValue("limit")

	if skip == "" || limit == "" {
		renderErrorResponse(r.Context(), w, "invalid request", errors.New("skip/limit missing"))
		return
	}

	skipINT, _ := strconv.Atoi(skip)   // NOTE: Safe to ignore error
	limitINT, _ := strconv.Atoi(limit) // NOTE: Safe to ignore error
	res, err := v.vs.Get(r.Context(), skipINT, limitINT)
	if err != nil {
		renderErrorResponse(r.Context(), w, "get failed", err)
		return
	}

	renderResponse(w, &GetVideosResponse{Videos: res}, http.StatusOK)
}

// SearchVideosRequest defines the request used for searching videos.
type SearchVideosRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

// SearchVideosResponse defines the response returned back after searching for any video.
type SearchVideosResponse struct {
	Videos []entity.Video `json:"videos"`
}

// search videos that matches provides search params
func (v *VideoHandler) search(w http.ResponseWriter, r *http.Request) {
	var req SearchVideosRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		renderErrorResponse(r.Context(), w, "invalid request", err)
		return
	}
	defer r.Body.Close()

	if req.Title == nil && req.Description == nil {
		renderErrorResponse(r.Context(), w, "invalid request", errors.New("no valid search param"))
		return
	}

	res, err := v.vs.Search(r.Context(), req.Title, req.Description)
	if err != nil {
		renderErrorResponse(r.Context(), w, "search failed", err)
		return
	}

	renderResponse(w, &SearchVideosResponse{Videos: res}, http.StatusOK)
}
