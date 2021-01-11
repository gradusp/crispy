package model

import "net"

type Real struct {
	tableName          struct{}          `pg:"controller.reals"`
	BalancingService   *BalancingService `json:"balancingService" pg:"rel:has-one"`
	BalancingServiceID string            `json:"-"`
	ID                 string            `json:"id" pg:"id,pk"`
	Addr               net.IP            `json:"addr" pg:"addr"`
	Port               int               `json:"port" pg:"port"`
}
