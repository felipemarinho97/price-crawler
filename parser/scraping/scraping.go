package scraping

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/felipemarinho97/price-crawler/requester"
)

type WebsiteSelectorSpec struct {
	Selector    string `json:"selector"`
	PostProcess func(*requester.Requester, string) string
}

func GetDocument(rq *requester.Requester, url string) (*goquery.Document, error) {
	res, err := rq.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get body: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create document: %w", err)
	}

	return doc, nil
}
