package main

import (
	"net/http"
	"os"

	"github.com/felipemarinho97/price-crawler/handlers"
	"github.com/felipemarinho97/price-crawler/requester"
)

func main() {
	// create a mux router
	indexerMux := http.NewServeMux()

	// create a requester service
	requesterURL := os.Getenv("REQUESTER_URL")
	if requesterURL == "" {
		requesterURL = "http://localhost:9080"
	}
	rs := requester.NewRequester(requesterURL)

	// register the handler for the /indexer endpoint
	indexer := handlers.NewIndexer(rs)
	indexerMux.HandleFunc("/product", indexer.HandleProduct)
	indexerMux.HandleFunc("/search", indexer.HandleSearch)

	// start the server and wait for signals to stop
	server := &http.Server{
		Addr:    ":8080",
		Handler: indexerMux,
	}

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
