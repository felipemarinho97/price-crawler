package client

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type ParserClient struct {
	URL    string
	client *http.Client
}

type SearchRequest struct {
	URL string `json:"url"`
}

type SearchResponse struct {
	Results []SearchResult `json:"results"`
}

type SearchResult struct {
	Name      string  `json:"name"`
	CashPrice float64 `json:"cashPrice"`
	Link      string  `json:"link"`
}

func NewParserClient(url string) *ParserClient {
	return &ParserClient{
		URL:    url,
		client: &http.Client{},
	}
}

func (c *ParserClient) Search(request SearchRequest) (SearchResponse, error) {
	// POST /search
	requestJSON, err := json.Marshal(request)
	if err != nil {
		return SearchResponse{}, err
	}
	req, err := http.NewRequest("POST", c.URL+"/search", bytes.NewBuffer(requestJSON))
	if err != nil {
		return SearchResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return SearchResponse{}, err
	}
	defer resp.Body.Close()

	var response SearchResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return SearchResponse{}, err
	}

	return response, nil
}
