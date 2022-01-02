package models

import "time"

type Cart struct {
	Id         string
	Box        int
	Data       []byte
	CreateTime time.Time
}
