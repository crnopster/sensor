package metric

import "github.com/google/uuid"

import "time"

// Metric .
type Metric struct {
	UUID      uuid.UUID
	Data      float32
	Type      string
	Timestamp time.Time
}
