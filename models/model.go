package models

import "time"

// Card main model.
type Card struct {
	ID         string
	Box        int
	Data       []byte
	CreateTime *time.Time
}
