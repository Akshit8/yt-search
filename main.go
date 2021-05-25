package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Akshit8/yt-search/rest"
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

func main() {
	var env string

	flag.StringVar(&env, "env", "", "Environment Variables filename")
	flag.Parse()

	if err := load(env); err != nil {
		log.Fatalln("Couldn't load configuration", err)
	}

	r := mux.NewRouter()

	rest.NewVideoHandler().Register(r)

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
