package model

import "net"

type Real struct {
	ID        string `json:"id"`
	ServiceID string `json:"serviceId"`
	Addr      net.IP `json:"addr"`
	Port      int    `json:"port"`
}
