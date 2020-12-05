package model

import "net"

type Real struct {
	ID               string            `json:"id" pg:"id,pk"`
	Addr             net.IPAddr        `json:"ip" pg:"addr"`
	Port             int               `json:"port" pg:"port"`
	HealthcheckAddr  string            `json:"healthcheckAddr" pg:"hc_addr"`
	BalancingService *BalancingService `json:"balancingService" pg:"rel:has-one"`
}
