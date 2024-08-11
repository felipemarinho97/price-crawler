package requester

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Requester struct {
	addr   string
	client *http.Client
}

func NewRequester(addr string) *Requester {
	return &Requester{
		addr:   addr,
		client: &http.Client{},
	}
}

func (r *Requester) Get(url string) (*http.Response, error) {
	body := map[string]string{"url": url}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/request", r.addr), bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return r.client.Do(req)
}

func (r *Requester) Post(url string, data interface{}) (*http.Response, error) {
	body := map[string]interface{}{"url": url, "data": data, "method": "POST"}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/request", r.addr), bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return r.client.Do(req)
}
