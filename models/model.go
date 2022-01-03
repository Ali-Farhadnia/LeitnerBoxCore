package models

import "time"

type CreateTime struct {
	Year  int
	Month time.Month
	Day   int
	Hour  int
}
type Cart struct {
	ID         string
	Box        int
	Data       []byte
	CreateTime CreateTime
}
