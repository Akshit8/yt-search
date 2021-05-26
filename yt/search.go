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

// YoutubeApi defines all actions for youtube api.
type YoutubeApi struct {
	service *youtube.Service
	vs      *elasticsearch.VideoSearch
}

// NewYoutubeApi creates new instance of YoutubeApi
func NewYoutubeApi(apiKey string, vs *elasticsearch.VideoSearch) (*YoutubeApi, error) {
	service, err := youtube.NewService(
		context.Background(),
		option.WithAPIKey(apiKey),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating new service: %v", err)
	}

	return &YoutubeApi{service: service, vs: vs}, nil
}

func (y *YoutubeApi) parseTime(timeStr string) time.Time {
	layout := "2006-01-02T15:04:05Z"

	t, err := time.Parse(layout, timeStr)

	if err != nil {
		return time.Time{}
	}

	return t
}

func (y *YoutubeApi) Search(ctx context.Context) {
	call := y.service.Search.List([]string{"id", "snippet"}).
		Q("blockchain|travel|tesla|bitcoin|dogecoin|").
		MaxResults(60).
		Type("video").
		Order("date").
		PublishedAfter(time.Now().Add(-15 * 60 * time.Hour).Format(time.RFC3339))

	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error fetching search response: %v", err)
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
			log.Fatalln("Error indexing video ", err)
		}
	}
}
