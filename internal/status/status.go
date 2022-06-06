// Package status provides all components required to serve data ingestion status via an HTTP API. It includes
// the transport & persistence layers.
package status

import "time"

type (
	// The Status type contains fields describing the current status of the data within the backend.
	Status struct {
		LastIngestTimestamp time.Time `json:"lastIngestTimestamp"`
	}
)
