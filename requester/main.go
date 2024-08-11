package main

import (
	"net/http"
	"os"
	"time"

	"github.com/felipemarinho97/price-crawler-requester/cache"
	"github.com/felipemarinho97/price-crawler-requester/cookies"
	"github.com/felipemarinho97/price-crawler-requester/flaresolverr"
	"github.com/felipemarinho97/price-crawler-requester/handlers"
	"github.com/felipemarinho97/price-crawler-requester/proxy"
)

func main() {
	// create a mux router
	indexerMux := http.NewServeMux()

	// create flare solver
	flaresolverrURL := os.Getenv("FLARESOLVERR_URL")
	if flaresolverrURL == "" {
		flaresolverrURL = "http://localhost:8191"
	}
	fs := flaresolverr.NewFlareSolverr(flaresolverrURL, 60000)

	// create a user cookie
	var cookiesFile string = os.Getenv("COOKIES_FILE")
	if cookiesFile == "" {
		cookiesFile = "./cookies.txt"
	}
	uc, err := cookies.NewCookie(cookiesFile)
	if err != nil {
		panic(err)
	}

	// create a proxy
	var proxyURL string = os.Getenv("PROXY_URL")
	if proxyURL == "" {
		proxyURL = "http://localhost:8080"
	}
	px := proxy.NewProxy(proxyURL)

	// create a cache
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}
	c := cache.NewRedisCache(redisURL, 30*time.Minute)

	// register the handler for the /indexer endpoint
	indexer := handlers.NewIndexer(c, fs, uc, px)
	indexerMux.HandleFunc("/request", indexer.HandleFlareSolvarr)

	// start the server and wait for signals to stop
	server := &http.Server{
		Addr:    ":9080",
		Handler: indexerMux,
	}

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
