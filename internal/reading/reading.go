// Package reading provides the full application stack for sensor readings. Including transport, persistence and
// event handling.
package reading

import (
	"fmt"
	"time"
)

type (
	// The Reading type describes a single sensor reading as stored in the database.
	Reading struct {
		Sensor    string    `json:"sensor"`
		Value     float64   `json:"value"`
		Timestamp time.Time `json:"timestamp"`
	}
)

// String returns a string representation of the reading.
func (r Reading) String() string {
	return fmt.Sprint(r.Sensor, " ", r.Value, r.Timestamp)
}

// Valid returns true if the Reading is deemed to be in a valid state. This means a sensor name, a zero or
// positive value and a non-zero timestamp.
func (r Reading) Valid() bool {
	return r.Sensor != "" && r.Value >= 0 && !r.Timestamp.IsZero()
}
