package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/felipemarinho97/price-crawler/requester"
	"github.com/felipemarinho97/price-crawler/scraping"
)

type indexer struct {
	rs *requester.Requester
}

func NewIndexer(rs *requester.Requester) *indexer {
	return &indexer{rs: rs}
}

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Name        string `json:"name"`
	CashPrice   string `json:"cashPrice"`
	CreditPrice string `json:"creditPrice"`
}

func (i *indexer) HandleProduct(w http.ResponseWriter, r *http.Request) {
	// get the page URL from the request body
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	// get the document from the URL
	doc, err := scraping.GetDocument(i.rs, req.URL)
	if err != nil {
		err = fmt.Errorf("failed to get document: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get prices from the document
	price, err := scraping.GetPrice(doc, i.rs)
	if err != nil {
		err = fmt.Errorf("failed to get price: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get product name from the document
	name, err := scraping.GetProductName(doc)
	if err != nil {
		err = fmt.Errorf("failed to get product name: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create the response
	resp := Response{
		Name:        name,
		CashPrice:   price.CashPrice,
		CreditPrice: price.CreditPrice,
	}

	// encode the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (i *indexer) HandleSearch(w http.ResponseWriter, r *http.Request) {
	// get the page URL from the request body
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	// get the document from the URL
	doc, err := scraping.GetDocument(i.rs, req.URL)
	if err != nil {
		err = fmt.Errorf("failed to get document: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get search results from the document
	results, err := scraping.GetSearch(doc, i.rs)
	if err != nil {
		err = fmt.Errorf("failed to get search results: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// sort results by price in ascending order
	sort.Slice(results.Results, func(i, j int) bool {
		return results.Results[i].CashPrice < results.Results[j].CashPrice
	})

	// encode the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

}
