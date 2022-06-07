// Package statistics provides all components required to serve statistical data via an HTTP API. It includes
// the transport & persistence layers.
package statistics

import (
	"time"

	"github.com/cloud-lada/backend/internal/reading"
)

type (
	// The Statistics type contains fields describing the current state of the Lada.
	Statistics struct {
		Speed             float64 `json:"speed"`
		Fuel              float64 `json:"fuel"`
		EngineTemperature float64 `json:"engineTemperature"`
		Revolutions       float64 `json:"revolutions"`
	}

	// The Statistic type contains fields describing a bucketed sensor value at a given time.
	Statistic struct {
		Sensor    reading.SensorType `json:"sensor"`
		Value     float64            `json:"value"`
		Timestamp time.Time          `json:"timestamp"`
	}
)
