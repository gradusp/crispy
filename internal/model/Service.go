package model

import "net"

type Service struct {
	Observable    `json:"-"`
	ClusterID     string `json:"clusterId"`
	ID            string `json:"id"`
	BalancingType string `json:"balancingType"`
	RoutingType   string `json:"routingType"`
	Bandwidth     int    `json:"bandwidth"`
	Proto         string `json:"proto"`
	Addr          net.IP `json:"addr"`
	Port          int    `json:"port"`
	// Reals         []*Real `json:"reals"`
	// *Healthchecks `json:"healthchecks"`
}
