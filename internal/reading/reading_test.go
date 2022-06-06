package reading_test

import (
	"testing"
	"time"

	"github.com/cloud-lada/backend/internal/reading"
	"github.com/stretchr/testify/assert"
)

func TestReading_Valid(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name     string
		Input    reading.Reading
		Expected bool
	}{
		{
			Name:     "It should return true for a valid reading",
			Expected: true,
			Input: reading.Reading{
				Sensor:    reading.SensorTypeSpeed,
				Value:     100,
				Timestamp: time.Now(),
			},
		},
		{
			Name: "It should return false for an invalid sensor",
			Input: reading.Reading{
				Sensor:    "invalid",
				Value:     100,
				Timestamp: time.Now(),
			},
		},
		{
			Name: "It should return false for a zero timestamp",
			Input: reading.Reading{
				Sensor:    reading.SensorTypeRevolution,
				Value:     100,
				Timestamp: time.Time{},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			assert.EqualValues(t, tc.Expected, tc.Input.Valid())
		})
	}
}
