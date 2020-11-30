package model

type SecurityZone struct {
	ID   string `json:"id" pg:"id,pk"`
	Name string `json:"name" pg:"name"`
}
