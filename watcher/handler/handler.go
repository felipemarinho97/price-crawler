package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/felipemarinho97/price-crawler/watcher/client"
)

type Handler struct {
	db     *client.DataBucketClient
	parser *client.ParserClient
}

func NewHandler(db *client.DataBucketClient, parser *client.ParserClient) *Handler {
	return &Handler{
		db:     db,
		parser: parser,
	}
}

type UpdatePriceRequest struct {
	SearchLinks []string `json:"searchLinks"`
}

// HandleUpdatePrice is the handler for the /update-prices endpoint
// It receives a search link and updates the price of the products
func (h *Handler) HandleUpdatePrices(w http.ResponseWriter, r *http.Request) {
	// decode the request body
	var req UpdatePriceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		msg := fmt.Sprintf("failed to decode request body: %v", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if len(req.SearchLinks) == 0 {
		http.Error(w, "no search links provided", http.StatusBadRequest)
		return
	}

	// update the price for each search link
	for _, searchLink := range req.SearchLinks {
		err := h.UpdatePrice(searchLink)
		if err != nil {
			msg := fmt.Sprintf("failed to update price: %v", err)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
	}

	// write the response
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UpdatePrice(searchLink string) error {
	// get the product name, cash price, and credit price from the parser
	// and update the price in the database
	currentPrices, err := h.parser.Search(client.SearchRequest{URL: searchLink})
	if err != nil {
		return fmt.Errorf("failed to get price from parser: %w", err)
	}

	pricesFiltered := []client.SearchResult{}
	for _, result := range currentPrices.Results {
		if result.CashPrice != 0 {
			pricesFiltered = append(pricesFiltered, result)
		}
	}

	var errChan = make(chan error)
	for _, result := range pricesFiltered {
		go func(sr client.SearchResult) {
			errChan <- h.db.AddDataPoint(client.DataPoint{
				Name:      sr.Name,
				Value:     sr.CashPrice,
				Timestamp: time.Now(),
			})
		}(result)
	}

	for range pricesFiltered {
		if err := <-errChan; err != nil {
			fmt.Printf("Error while adding data point: %v\n", err)
		}
	}

	return nil
}
