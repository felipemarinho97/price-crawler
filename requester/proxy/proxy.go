package proxy

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
)

// Proxy is a struct that represents a proxy server.
type Proxy struct {
	// The URL of the proxy server.
	URL    url.URL
	client http.Client
}

// NewProxy creates a new Proxy with the given URL.
func NewProxy(u string) *Proxy {
	_u, err := url.Parse(u)
	if err != nil {
		panic(err)
	}
	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(_u),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return &Proxy{URL: *_u, client: client}
}

// Get performs a GET request to the given URL using the proxy server.
func (p *Proxy) Get(url string) (io.ReadCloser, error) {
	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Perform the request
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
