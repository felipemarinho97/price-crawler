package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/felipemarinho97/price-crawler-requester/cache"
	"github.com/felipemarinho97/price-crawler-requester/cookies"
	"github.com/felipemarinho97/price-crawler-requester/flaresolverr"
	"github.com/felipemarinho97/price-crawler-requester/proxy"
)

type indexer struct {
	c  cache.Cache
	fs *flaresolverr.FlareSolverr
	uc *cookies.UserCookie
	px *proxy.Proxy
}

func NewIndexer(c cache.Cache, fs *flaresolverr.FlareSolverr, uc *cookies.UserCookie, px *proxy.Proxy) *indexer {
	return &indexer{fs: fs, uc: uc, px: px, c: c}
}

type Request struct {
	URL    string `json:"url"`
	Data   []byte `json:"data"`
	Method string `json:"method"`
}

func (i *indexer) getDocument(url string) (io.ReadCloser, error) {
	var body io.ReadCloser

	// check if the document is cached
	cached, err := i.c.Get(context.Background(), url)
	if err == nil {
		fmt.Printf("Document %s is cached\n", url)
		return io.NopCloser(bytes.NewBufferString(cached)), nil
	}

	body, err = i.fs.Get(url)
	if err != nil {
		// try request with cookies
		fmt.Printf("Trying request %s with cookies\n", url)

		body, err = i.uc.Get(url)
	}

	if err != nil {
		// try request with proxy
		fmt.Printf("Trying request %s with proxy\n", url)

		body, err = i.px.Get(url)
	}

	if err != nil {
		fmt.Printf("Failed to get document: %v\n", err)
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	// cache the document
	doc := getAsStr(body)
	err = i.c.Set(context.Background(), url, doc)
	if err != nil {
		fmt.Printf("Failed to cache document: %v\n", err)
	}
	body = io.NopCloser(bytes.NewBufferString(doc))

	return body, nil
}

func getAsStr(data io.ReadCloser) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(data)
	return buf.String()
}

func (i *indexer) postDocument(url string, data []byte) (io.ReadCloser, error) {
	var body io.ReadCloser
	// try request with cookies
	client := http.Client{}

	fmt.Println("Trying request with cookies")
	req, err := http.NewRequest("POST", url, io.NopCloser(bytes.NewReader(data)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}

	body = resp.Body

	return body, nil
}

func (i *indexer) HandleFlareSolvarr(w http.ResponseWriter, r *http.Request) {
	// get the page URL from the request body
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	var doc io.ReadCloser
	if req.Method == "" || req.Method == "GET" {
		// get the document from the URL
		doc, err = i.getDocument(req.URL)
		if err != nil {
			err = fmt.Errorf("failed to get document: %w", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if req.Method == "POST" {
		// get the document from the URL
		doc, err = i.postDocument(req.URL, req.Data)
		if err != nil {
			err = fmt.Errorf("failed to get document: %w", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// return the document
	w.WriteHeader(http.StatusOK)

	_, err = io.Copy(w, doc)
	if err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
		return
	}
}
