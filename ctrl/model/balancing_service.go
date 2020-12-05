package model

import "net"

type BalancingType int
type RoutingType int

const (
	NAT BalancingType = 0 << iota
	Tunnel
	GRETunnel

	RoundRobin RoutingType = 0 << iota
	SourceIP
	SourceIPPort
	LeastConnection
	URI
	Random
	RDPCookie
)

type BalancingService struct {
	ID            string        `json:"id" pg:"id,pk"`
	BalancingType BalancingType `json:"balancingType" pg:"balancing_type"`
	RoutingType   RoutingType   `json:"routingType" pg:"routing_type"`
	Proto         string        `json:"proto" pg:"proto"`
	Addr          net.IPAddr    `json:"addr" pg:"addr"`
	Port          int           `json:"port" pg:"port"`
	Cluster       *Cluster      `json:"cluster" pg:"rel:has-one"`
}

//type BalancingService struct {
//	IP              string `json:"ip"`
//	Port            int    `json:"port"`
//	BalanceType     string `json:"balanceType"`
//	RoutingType     string `json:"routingType"`
//	Protocol        string `json:"protocol"`
//	Quorum          int    `json:"quorum"`
//	Hysteresis      int    `json:"hysteresis"`
//	HealthcheckType string `json:"healthcheckType"`
//	HelloTimer      string `json:"helloTimer"`
//	ResponseTimer   string `json:"responseTimer"`
//	AliveThreshold  int    `json:"aliveThreshold"`
//	deadThreshold   int    `json:"deadThreshold"`
//}

//type BalancingService struct {
//	BalancingServiceID string
//	OrderID            string
//	VipID              string
//	RoutingTypeID      int
//	BalancingType      string
//}
