package model

type Cluster struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Capacity int64  `json:"capacity"`
	Usage    int64  `json:"usage"`
	ZoneID   string `json:"zoneId,omitempty"`
	Zone     *Zone  `json:"zone,omitempty"`
}
