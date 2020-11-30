package model

type Server struct {
	ID         string `json:"id"`
	Hostname   string `json:"hostname"`
	Address    string `json:"address"`
	Datacenter string `json:"datacenter"`
}
