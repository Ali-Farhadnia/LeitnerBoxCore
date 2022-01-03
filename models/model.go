package models

import "time"

type Cart struct {
	ID         string
	Box        int
	Data       []byte
	CreateTime time.Time
}
