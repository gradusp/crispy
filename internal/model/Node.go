package model

import "net"

type Node struct {
	Observable `json:"-"`
	ID         int    `json:"id"`
	ClusterID  string `json:"clusterId,omitempty"`
	Addr       net.IP `json:"addr"`
	Hostname   string `json:"hostname"`
}
