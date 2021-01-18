package model

import "net"

type Real struct {
	ID              string   `json:"id"`
	Addr            net.IP   `json:"addr"`
	Port            int      `json:"port"`
	HealthcheckAddr net.IP   `json:"hcAddr"`
	HealthcheckPort int      `json:"hcPort"`
	ServiceID       string   `json:"serviceId,omitempty"`
	Service         *Service `json:"service,omitempty"`
}
