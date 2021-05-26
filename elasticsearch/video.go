package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/Akshit8/yt-search/entity"
	es "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// VideoSearch represents repository for interacting with elastic-search for resource Video.
type VideoSearch struct {
	client *es.Client
	index  string
}

// NewVideoSearch creates new instance
func NewVideoSearch(client *es.Client) *VideoSearch {
	return &VideoSearch{
		client: client,
		index:  "videos",
	}
}

// Index creates or updates a task in an index.
func (v *VideoSearch) Index(ctx context.Context, video entity.Video) error {
	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(video)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      v.index,
		Body:       &buf,
		DocumentID: video.ID,
		Refresh:    "true",
	}

	resp, err := req.Do(ctx, v.client)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return fmt.Errorf("IndexRequest error: %d", resp.StatusCode)
	}

	io.Copy(ioutil.Discard, resp.Body)

	return nil
}

// Get return all videos indexed in elasticsearch with pagination.
func (v *VideoSearch) Get(ctx context.Context, skip, limit int) ([]entity.Video, error) {
	query := map[string]interface{}{
		"from": skip,
		"size": limit,
		"sort": map[string]interface{}{
			"published_at": map[string]interface{}{
				"order": "desc",
			},
		},
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}

	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(query)
	if err != nil {
		return nil, err
	}

	req := esapi.SearchRequest{
		Index: []string{"videos"},
		Body:  &buf,
	}

	client := v.client
	resp, err := req.Do(ctx, client)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return nil, fmt.Errorf("SearchRequest error: %d", resp.StatusCode)
	}

	var hits struct {
		Hits struct {
			Hits []struct {
				Source entity.Video `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	err = json.NewDecoder(resp.Body).Decode(&hits)
	if err != nil {
		fmt.Println("Error here", err)
		return nil, err
	}

	res := make([]entity.Video, len(hits.Hits.Hits))

	for i, hit := range hits.Hits.Hits {
		res[i].ID = hit.Source.ID
		res[i].Title = hit.Source.Title
		res[i].Description = hit.Source.Description
		res[i].PublishedAt = hit.Source.PublishedAt
		res[i].Thumbnail = hit.Source.Thumbnail
	}

	return res, nil
}

// Search returns videos matching a query.
func (v *VideoSearch) Search(ctx context.Context, title, description *string) ([]entity.Video, error) {
	should := make([]interface{}, 0, 2)

	if title != nil {
		should = append(should, map[string]interface{}{
			"match": map[string]interface{}{
				"title": *title,
			},
		})
	}

	if description != nil {
		should = append(should, map[string]interface{}{
			"match": map[string]interface{}{
				"description": *description,
			},
		})
	}

	var query map[string]interface{}

	if len(should) > 1 {
		query = map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"should": should,
				},
			},
		}
	} else {
		query = map[string]interface{}{
			"query": should[0],
		}
	}

	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(query)
	if err != nil {
		return nil, err
	}

	req := esapi.SearchRequest{
		Index: []string{"videos"},
		Body:  &buf,
	}

	client := v.client
	resp, err := req.Do(ctx, client)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return nil, fmt.Errorf("SearchRequest error: %d", resp.StatusCode)
	}

	var hits struct {
		Hits struct {
			Hits []struct {
				Source entity.Video `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	err = json.NewDecoder(resp.Body).Decode(&hits)
	if err != nil {
		fmt.Println("Error here", err)
		return nil, err
	}

	res := make([]entity.Video, len(hits.Hits.Hits))

	for i, hit := range hits.Hits.Hits {
		res[i].ID = hit.Source.ID
		res[i].Title = hit.Source.Title
		res[i].Description = hit.Source.Description
		res[i].PublishedAt = hit.Source.PublishedAt
		res[i].Thumbnail = hit.Source.Thumbnail
	}

	return res, nil
}
