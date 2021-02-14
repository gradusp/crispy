package model

type Cluster struct {
	Observable `json:"-"`
	ID         string `json:"id"`
	Name       string `json:"name"`
	Capacity   int64  `json:"capacity,omitempty"`
	Usage      int64  `json:"usage,omitempty"`
	ZoneID     string `json:"zoneId,omitempty"`
	Zone       *Zone  `json:"zone,omitempty"`
}
