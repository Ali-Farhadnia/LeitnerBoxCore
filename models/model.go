package models

import "time"

// Card main model.
type Card struct {
	ID         string
	Box        int
	Data       []byte
	CreateTime *time.Time
}

// NewCard return empty card.
func NewCard() *Card {
	return &Card{ID: "", Box: 0, Data: nil, CreateTime: &time.Time{}}
}
