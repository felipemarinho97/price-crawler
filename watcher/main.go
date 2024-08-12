package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/felipemarinho97/price-crawler/watcher/client"
	"github.com/felipemarinho97/price-crawler/watcher/handler"
)

func main() {
	wactchedSearchs := strings.Split(os.Getenv("WATCHED_SEARCHES"), ",")
	if len(wactchedSearchs) == 0 {
		panic("no watched searchs provided")
	}

	dataBucketURL := os.Getenv("DATA_BUCKET_URL")
	if dataBucketURL == "" {
		dataBucketURL = "http://localhost:9081"
	}
	parserURL := os.Getenv("PARSER_URL")
	if parserURL == "" {
		parserURL = "http://localhost:8080"
	}

	// create clients
	db := client.NewDataBucketClient(dataBucketURL)
	parser := client.NewParserClient(parserURL)

	// create the handler
	h := handler.NewHandler(db, parser)

	// create the server
	mux := http.NewServeMux()
	mux.HandleFunc("/update-prices", h.HandleUpdatePrices)

	server := &http.Server{
		Addr:    ":9082",
		Handler: mux,
	}

	// start the server and wait for signals to stop inside a goroutine
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// wait for signals to stop
	select {}

}
