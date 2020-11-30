package model

type Cluster struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	SecurityZone *SecurityZone `json:"securityZone" pg:"rel:has-one"`
}

// https://pkg.go.dev/github.com/go-pg/pg/v10#example-DB.Model-HasOne
