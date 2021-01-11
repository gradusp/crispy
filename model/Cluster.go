package model

type Cluster struct {
	tableName      struct{}      `pg:"controller.clusters"`
	SecurityZone   *SecurityZone `json:",omitempty" pg:"rel:has-one"`
	SecurityZoneID string        `json:"-"`
	ID             string        `json:"id" pg:"id,pk"`
	Name           string        `json:"name" pg:"name"`
	Capacity       int64         `json:"capacity" pg:"capacity"`
}
