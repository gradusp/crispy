package model

import "net"

type Service struct {
	tableName     struct{}       `json:",omitempty" pg:"controller.balancing_services"`
	ID            string         `json:"id" pg:"id,pk"`
	BalancingType string         `json:"balancingType" pg:"balancing_type"`
	RoutingType   string         `json:"routingType" pg:"routing_type"`
	Bandwidth     int            `json:"bandwidth" pg:"bandwidth"`
	Proto         string         `json:"proto" pg:"proto"`
	Addr          net.IP         `json:"addr" pg:"addr"`
	Port          int            `json:"port" pg:"port"`
	ClusterID     string         `json:"-" pg:"cluster_id"`
	Cluster       Cluster        `json:"cluster" pg:"rel:has-one"`
	Reals         []*Real        `json:"reals" pg:"rel:has-many"`
	Healthchecks  []*Healthcheck `json:"healthchecks" pg:"rel:has-many"`
	//Orders        []Order  `pg:"many2many:controller.order_to_balancing_service"`
}

// https://pkg.go.dev/github.com/go-pg/pg/v10#example-DB.Model-HasOne
// https://github.com/go-pg/pg/issues/229 TODO: how to use db schema
