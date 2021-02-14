package model

type Zone struct {
	Observable `json:"-"`
	ID         string `json:"id"`
	Name       string `json:"name"`
}
