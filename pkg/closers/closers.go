// Package closers provides utilities for working with io.Closer implementations.
package closers

import (
	"io"
	"log"
)

// Close the io.Closer. Logging the error if it is non-nil.
func Close(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Printf("failed to close %T: %v", c, err)
	}
}
