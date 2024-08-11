package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/felipemarinho97/price-crawler-requester/cookies"
	"github.com/felipemarinho97/price-crawler-requester/flaresolverr"
)

type indexer struct {
	fs *flaresolverr.FlareSolverr
	uc *cookies.UserCookie
}

func NewIndexer(fs *flaresolverr.FlareSolverr, uc *cookies.UserCookie) *indexer {
	return &indexer{fs: fs, uc: uc}
}

type Request struct {
	URL    string `json:"url"`
	Data   []byte `json:"data"`
	Method string `json:"method"`
}

func GetDocument(fs *flaresolverr.FlareSolverr, uc *cookies.UserCookie, url string) (io.ReadCloser, error) {
	var body io.ReadCloser
	body, err := fs.Get(url)
	if err != nil {
		err = fmt.Errorf("failed to get body: %w", err)
		fmt.Printf("Error: %v\n", err)

		// try request with cookies
		client := http.Client{}

		fmt.Println("Trying request with cookies")
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		addUserDetails(req, uc, url)

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to do request: %w", err)
		}

		body = resp.Body
	}

	return body, nil
}

func PostDocument(fs *flaresolverr.FlareSolverr, uc *cookies.UserCookie, url string, data []byte) (io.ReadCloser, error) {
	var body io.ReadCloser
	// try request with cookies
	client := http.Client{}

	fmt.Println("Trying request with cookies")
	req, err := http.NewRequest("POST", url, io.NopCloser(bytes.NewReader(data)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	addUserDetails(req, uc, url)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}

	body = resp.Body

	return body, nil
}

func addUserDetails(req *http.Request, uc *cookies.UserCookie, url string) {
	// add headers
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/png,image/svg+xml,*/*;q=0.8")
	req.Header.Add("Accept-Language", "en,pt-BR;q=0.8,pt;q=0.5,en-US;q=0.3")
	//req.Header.Add("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Add("Referer", gerReferer(url))
	req.Header.Add("DNT", "1")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("Sec-Fetch-Dest", "document")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("Sec-Fetch-User", "?1")
	req.Header.Add("Priority", "u=0, i")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("TE", "trailers")
	uc.AddCookies(req)
}

func gerReferer(target string) string {
	targetURL, err := url.Parse(target)
	if err != nil {
		return target
	}
	return targetURL.Scheme + "://" + targetURL.Host
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
		doc, err = GetDocument(i.fs, i.uc, req.URL)
		if err != nil {
			err = fmt.Errorf("failed to get document: %w", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if req.Method == "POST" {
		// get the document from the URL
		doc, err = PostDocument(i.fs, i.uc, req.URL, req.Data)
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
