package main

import (
	"net/http"
	"os"

	"github.com/felipemarinho97/price-crawler/data-bucket/databucket"
	"github.com/felipemarinho97/price-crawler/data-bucket/handler"
)

func main() {
	// create a mux router
	mux := http.NewServeMux()

	// create the databucket service
	var postgresURL string = os.Getenv("POSTGRES_URL")
	if postgresURL == "" {
		postgresURL = "postgres://postgres:password@localhost:5434/postgres?sslmode=disable"
	}
	db, err := databucket.NewPostgresDataBucket(postgresURL)
	if err != nil {
		panic(err)
	}

	// create the handler for the /datapoint endpoint
	h := handler.NewHandler(db)

	// register the handler for the /datapoint endpoint
	mux.HandleFunc("/datapoints", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.HandleGetDatapoint(w, r)
		case http.MethodPost:
			h.HandlePostDatapoint(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/datapoints/name", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.HandleListDatapointNames(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// start the server and wait for signals to stop
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
