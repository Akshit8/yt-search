package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

const developerKey = "AIzaSyCx8fTFyaMdxYxAQllzH630Lx4MwhjqEOQ"

func main() {

	service, err := youtube.NewService(context.Background(), option.WithAPIKey(developerKey))
	if err != nil {
		log.Fatalf("Error creating new Youtube client: %v", err)
	}

	call := service.Search.List([]string{"id", "snippet"}).
		Q("blockchain").
		MaxResults(5).
		Type("video")

	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error fetching search response: %v", err)
	}

	for _, item := range response.Items {
		fmt.Printf("Video id: %s \n Title: %s \n Description: %s \n Publishing Datetime: %s \n Thumbnail: %s \n \n----------\n",
			item.Id.VideoId,
			item.Snippet.Title,
			item.Snippet.Description,
			item.Snippet.PublishedAt,
			item.Snippet.Thumbnails.Default.Url,
		)
	}

	call2 := service.Videos.List([]string{"id", "snippet"}).Id("hYip_Vuv8J0")
	video, err := call2.Do()
	if err != nil {
		log.Fatalf("Error fetching search response: %v", err)
	}

	fmt.Printf("Video id: %s \n Title: %s \n Description: %s \n ",
		video.Items[0].Id,
		video.Items[0].Snippet.Title,
		video.Items[0].Snippet.Description,
	)
}
