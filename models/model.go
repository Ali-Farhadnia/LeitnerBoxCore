package models

import "time"

type Card struct {
	ID         string
	Box        int
	Data       []byte
	CreateTime *time.Time
}
