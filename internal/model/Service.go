package model

import "net"

type Service struct {
	ID            string `json:"id"`
	BalancingType string `json:"balancingType"`
	RoutingType   string `json:"routingType"`
	Bandwidth     int    `json:"bandwidth"`
	Proto         string `json:"proto"`
	Addr          net.IP `json:"addr"`
	Port          int    `json:"port"`
	ClusterID     string `json:"clusterId,omitempty"`
	//Cluster       Cluster        `json:"cluster,omitempty"`
	Reals        []*Real        `json:"reals"`
	Healthchecks []*Healthcheck `json:"healthchecks"`
	//Orders        []Order  `pg:"many2many:controller.order_to_balancing_service"`
}
