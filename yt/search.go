package yt

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Akshit8/yt-search/elasticsearch"
	"github.com/Akshit8/yt-search/entity"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// YoutubeAPI defines all actions for youtube api.
type YoutubeAPI struct {
	service *youtube.Service
	vs      *elasticsearch.VideoSearch
}

// NewYoutubeAPI creates new instance of YoutubeApi
func NewYoutubeAPI(apiKey string, vs *elasticsearch.VideoSearch) (*YoutubeAPI, error) {
	service, err := youtube.NewService(
		context.Background(),
		option.WithAPIKey(apiKey),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating new service: %v", err)
	}

	return &YoutubeAPI{service: service, vs: vs}, nil
}

// parseTime converts string time object
func (y *YoutubeAPI) parseTime(timeStr string) time.Time {
	layout := "2006-01-02T15:04:05Z"

	t, err := time.Parse(layout, timeStr)

	if err != nil {
		return time.Time{}
	}

	return t
}

// Search is a worker function that extracts videos from Youtube Date API and
// indexes them to elastic search.
func (y *YoutubeAPI) Search(ctx context.Context, query string, c chan error) {
	call := y.service.Search.List([]string{"id", "snippet"}).
		Q(query).
		MaxResults(60).
		Type("video").
		Order("date").
		PublishedAfter(time.Now().Add(-15 * 60 * time.Hour).Format(time.RFC3339))

	response, err := call.Do()
	if err != nil {
		log.Println("Error fetching search response ", err)
		c <- err
	}

	for _, item := range response.Items {
		video := entity.Video{
			ID:          item.Id.VideoId,
			Title:       item.Snippet.Title,
			Description: item.Snippet.Description,
			PublishedAt: y.parseTime(item.Snippet.PublishedAt),
			Thumbnail:   item.Snippet.Thumbnails.Default.Url,
		}

		err := y.vs.Index(ctx, video)
		if err != nil {
			log.Println("Error indexing video ", err)
			c <- err
		}
	}
}
