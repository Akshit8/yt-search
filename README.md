```go
fmt.Printf("Video id: %s \n Title: %s \n Description: %s \n Publishing Datetime: %s \n Thumbnail: %s \n \n----------\n",
			item.Id.VideoId,
			item.Snippet.Title,
			item.Snippet.Description,
			item.Snippet.PublishedAt,
			item.Snippet.Thumbnails.Default.Url,
		)
```