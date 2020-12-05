package model

type Cluster struct {
	ID             string        `json:"id" pg:"id,pk"`
	Name           string        `json:"name" pg:"name"`
	Capacity       int64         `json:"capacity" pg:"capacity"`
	Usage          int64         `json:"usage" pg:"usage"`
	SecurityZoneID string        `json:"securityZoneId" pg:"security_zone_id"`
	SecurityZone   *SecurityZone `json:"securityZone" pg:"rel:has-one"`
}

// https://pkg.go.dev/github.com/go-pg/pg/v10#example-DB.Model-HasOne
