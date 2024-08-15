package databucket

import (
	"time"
)

type DataPoint struct {
	// The ID of the data point.
	ID int `json:"id"`
	// The name of the data point.
	Name string `json:"name"`
	// The value of the data point.
	Value float64 `json:"value"`
	// The timestamp of the data point.
	Timestamp time.Time `json:"timestamp"`
}

type DataPointName struct {
	// The name of the data point.
	Name string `json:"name" db:"name"`
	// The average value of the data point.
	AvgValue float64 `json:"avgValue" db:"avg"`
	// The minimum value of the data point.
	MinValue float64 `json:"minValue" db:"min"`
	// The maximum value of the data point.
	MaxValue float64 `json:"maxValue" db:"max"`
	// The last value of the data point.
	LastValue float64 `json:"lastValue" db:"last"`
}

type DataBucket interface {
	// Add a data point to the bucket.
	AddDataPoint(dataPoint DataPoint) error
	// Get all data points in the bucket within a time range and with a specific name.
	GetDataPoints(name string, start time.Time, end time.Time) ([]DataPoint, error)
	// List all data point names in the bucket.
	ListDataPointNames() ([]DataPointName, error)
}
