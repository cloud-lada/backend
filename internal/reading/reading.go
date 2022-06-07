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
		Sensor    SensorType `json:"sensor"`
		Value     float64    `json:"value"`
		Timestamp time.Time  `json:"timestamp"`
	}

	// The SensorType type describes the kind of sensor that the value relates to.
	SensorType string
)

// Constants for sensor types.
const (
	SensorTypeSpeed             = SensorType("speed")
	SensorTypeFuel              = SensorType("fuel")
	SensorTypeRevolution        = SensorType("revolution")
	SensorTypeEngineTemperature = SensorType("engine_temperature")
	SensorTypeLocationLatitude  = SensorType("location_latitude")
	SensorTypeLocationLongitude = SensorType("location_longitude")
)

var validSensorTypes = map[SensorType]struct{}{
	SensorTypeSpeed:             {},
	SensorTypeFuel:              {},
	SensorTypeRevolution:        {},
	SensorTypeEngineTemperature: {},
	SensorTypeLocationLatitude:  {},
	SensorTypeLocationLongitude: {},
}

// String returns a string representation of the reading.
func (r Reading) String() string {
	return fmt.Sprint(r.Sensor, " ", r.Value, r.Timestamp)
}

// Valid returns true if the Reading is deemed to be in a valid state. This means a valid sensor type and a non-zero
// timestamp.
func (r Reading) Valid() bool {
	return r.Sensor.Valid() && !r.Timestamp.IsZero()
}

// Valid returns true if the SensorType is one of the valid types of sensor.
func (st SensorType) Valid() bool {
	_, ok := validSensorTypes[st]
	return ok
}
