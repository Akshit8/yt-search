package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/Akshit8/yt-search/elasticsearch"
	"github.com/Akshit8/yt-search/rest"
	"github.com/Akshit8/yt-search/yt"
	es "github.com/elastic/go-elasticsearch/v7"
	"github.com/gorilla/mux"
)

// one of the query from the list will be randomly selected and passed for search
var ytQueries = []string{"blockchain", "tesla", "dogecoin", "eth2.0", "elon musk", "maldives"}

// sets random seed on start of app.
func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// startPolling calls passed function every 30s in a seperate go routine
func startPolling(
	ctx context.Context,
	c chan error,
	f func(context.Context, string, chan error),
) {
	for {
		index := RandomInt(0, int64(len(ytQueries)-1))
		f(ctx, ytQueries[index], c)
		time.Sleep(20 * time.Minute)
	}
}

// newElasticSearch creates es cient and returns if connection to es is successful.
func newElasticSearch() (*es.Client, error) {
	client, err := es.NewDefaultClient()
	if err != nil {
		return nil, fmt.Errorf("elasticsearch.Open %w", err)
	}

	res, err := client.Info()
	if err != nil {
		return nil, fmt.Errorf("es.Info %w", err)
	}

	defer func() {
		err = res.Body.Close()
	}()

	return client, nil
}

func main() {
	var address string

	flag.StringVar(&address, "address", ":8000", "HTTP Server Address")
	flag.Parse()

	client, err := newElasticSearch()
	if err != nil {
		log.Fatalln("Couldn't create elasticsearch client ", err)
	}

	vs := elasticsearch.NewVideoSearch(client)

	api, err := yt.NewYoutubeAPI(os.Getenv("YOUTUBE_API_KEY"), vs)
	if err != nil {
		log.Fatalln("Error creating youtube api ", err)
	}

	errs := make(chan error, 1)

	go startPolling(context.Background(), errs, api.Search)

	r := mux.NewRouter()

	rest.NewVideoHandler(vs).Register(r)

	srv := &http.Server{
		Handler:           r,
		Addr:              address,
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       1 * time.Second,
	}

	log.Println("Starting server", address)

	log.Fatal(srv.ListenAndServe())
}
