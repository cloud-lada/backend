// Package location provides the full application stack for location data. Including transport & querying.
package location

type (
	// The Location type describes a point on the globe as a latitude and longitude.
	Location struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}
)
