package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Akshit8/yt-search/elasticsearch"
	"github.com/Akshit8/yt-search/rest"
	"github.com/Akshit8/yt-search/yt"
	es "github.com/elastic/go-elasticsearch/v7"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// load read the env filename and load it into ENV for this process.
func load(filename string) error {
	if err := godotenv.Load(filename); err != nil {
		return fmt.Errorf("loading env var file: %w", err)
	}

	return nil
}

// startPolling calls passed function every 30s in a seperate go routine
func startPolling(ctx context.Context, f func(context.Context)) {
	for {
		// f(ctx)
		time.Sleep(30 * time.Minute)
	}
}

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
	var env string

	flag.StringVar(&env, "env", "", "Environment Variables filename")
	flag.Parse()

	if err := load(env); err != nil {
		log.Fatalln("Couldn't load configuration ", err)
	}

	client, err := newElasticSearch()
	if err != nil {
		log.Fatalln("Couldn't create elasticsearch client ", err)
	}

	vs := elasticsearch.NewVideoSearch(client)

	api, err := yt.NewYoutubeAPI(os.Getenv("YOUTUBE_API_KEY"), vs)
	if err != nil {
		log.Fatalln("Error creating youtube api ", err)
	}

	go startPolling(context.Background(), api.Search)

	r := mux.NewRouter()

	rest.NewVideoHandler(vs).Register(r)

	address := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))

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
