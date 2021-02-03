package model

import "time"

type Audit struct {
	ID   string    `json:"id"`
	Time time.Time `json:"time"`
	What string    `json:"what"`
	Who  string    `json:"who"`
}
