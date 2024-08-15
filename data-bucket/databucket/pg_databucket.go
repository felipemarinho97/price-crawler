package databucket

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresDataBucket struct {
	// The database connection string.
	db *sqlx.DB
}

const createTableQuery = `
CREATE TABLE IF NOT EXISTS data_points (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	value DOUBLE PRECISION NOT NULL,
	timestamp TIMESTAMP NOT NULL
)`

func NewPostgresDataBucket(connectionString string) (*PostgresDataBucket, error) {
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	return &PostgresDataBucket{
		db: db,
	}, nil
}

func (p *PostgresDataBucket) AddDataPoint(dataPoint DataPoint) error {
	_, err := p.db.Exec("INSERT INTO data_points (name, value, timestamp) VALUES ($1, $2, $3)", dataPoint.Name, dataPoint.Value, dataPoint.Timestamp)
	return err
}

func (p *PostgresDataBucket) GetDataPoints(name string, start time.Time, end time.Time) ([]DataPoint, error) {
	var dataPoints []DataPoint
	err := p.db.Select(&dataPoints, `
	SELECT * FROM data_points 
		WHERE name = $1 AND timestamp >= $2 
		AND timestamp <= $3 
		ORDER BY timestamp ASC`, name, start, end)
	return dataPoints, err
}

func (p *PostgresDataBucket) ListDataPointNames() ([]DataPointName, error) {
	var names []DataPointName
	err := p.db.Select(&names, `
	SELECT name, 
		AVG(value),
		MIN(value),
		MAX(value),
		(SELECT value FROM data_points WHERE name = dp.name ORDER BY timestamp DESC LIMIT 1) AS last
	FROM data_points dp
	WHERE timestamp >= NOW() - INTERVAL '90 day'
	GROUP BY name 
	ORDER BY last ASC`)

	return names, err
}
