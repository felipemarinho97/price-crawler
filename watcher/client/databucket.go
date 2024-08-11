package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type DataBucketClient struct {
	URL    string
	client *http.Client
}

type DataPoint struct {
	Name      string    `json:"name"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

func NewDataBucketClient(url string) *DataBucketClient {
	return &DataBucketClient{
		URL:    url,
		client: &http.Client{},
	}
}

func (c *DataBucketClient) AddDataPoint(dataPoint DataPoint) error {
	// POST /datapoints
	dataJSON, err := json.Marshal(dataPoint)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.URL+"/datapoints", bytes.NewBuffer(dataJSON))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
