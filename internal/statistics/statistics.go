// Package statistics provides all components required to serve statistical data via an HTTP API. It includes
// the transport & persistence layers.
package statistics

type (
	// The Statistics type contains fields describing the current state of the Lada.
	Statistics struct {
		Speed             float64 `json:"speed"`
		Fuel              float64 `json:"fuel"`
		EngineTemperature float64 `json:"engineTemperature"`
		Revolutions       float64 `json:"revolutions"`
	}
)
