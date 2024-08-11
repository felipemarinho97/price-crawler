package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/felipemarinho97/price-crawler/data-bucket/databucket"
)

// Handler is the struct that holds the methods for the handler
type Handler struct {
	dataBucket databucket.DataBucket
}

// NewHandler creates a new handler
func NewHandler(dataBucket databucket.DataBucket) *Handler {
	return &Handler{dataBucket: dataBucket}
}

// HandleGetDatapoint handles the get request for the datapoint
func (h *Handler) HandleGetDatapoint(w http.ResponseWriter, r *http.Request) {
	// parse the query parameters
	query := r.URL.Query()
	name := query.Get("name")
	start, err := time.Parse(time.RFC3339, query.Get("start"))
	if err != nil {
		msg := fmt.Sprintf("failed to parse start time: %v", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	end, err := time.Parse(time.RFC3339, query.Get("end"))
	if err != nil {
		msg := fmt.Sprintf("failed to parse end time: %v", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// get the data points
	dataPoints, err := h.dataBucket.GetDataPoints(name, start, end)
	if err != nil {
		fmt.Printf("Error while getting data points: %v\n", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// write the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dataPoints)
}

// HandlePostDatapoint handles the post request for the datapoint
func (h *Handler) HandlePostDatapoint(w http.ResponseWriter, r *http.Request) {
	// decode the request body
	var dataPoint databucket.DataPoint
	err := json.NewDecoder(r.Body).Decode(&dataPoint)

	if err != nil {
		msg := fmt.Sprintf("failed to decode request body: %v", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// insert the data point
	err = h.dataBucket.AddDataPoint(dataPoint)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// write the response
	w.WriteHeader(http.StatusCreated)
}

// HandleListDatapointNames handles the list request for the datapoint names
func (h *Handler) HandleListDatapointNames(w http.ResponseWriter, r *http.Request) {
	// get the data point names
	names, err := h.dataBucket.ListDataPointNames()
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// write the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(names)
}
