package model

import "net"

type Healthchecks struct {
	HTTP []*HealthcheckHTTP `json:"http"`
	TCP  []*HealthcheckTCP  `json:"tcp"`
	UDP  []*HealthcheckUDP  `json:"udp"`
	ICMP []*HealthcheckICMP `json:"icmp"`
}

// HealthcheckBase is not used directly in code but designed for composite purpose
type HealthcheckBase struct {
	ID             string   `json:"id"`
	ServiceID      string   `json:"serviceId,omitempty"`
	Service        *Service `json:"service,omitempty"`
	Addr           net.IP   `json:"addr"`
	HelloTimer     int      `json:"helloTimer"`
	ResponseTimer  int      `json:"responseTimer"`
	AliveThreshold int      `json:"aliveThreshold"`
	DeadThreshold  int      `json:"deadThreshold"`
	Quorum         int      `json:"quorum"`
	Hysteresis     int      `json:"hysteresis"`
}

type HealthcheckICMP struct {
	HealthcheckBase
	Addr net.IP `json:"addr"`
}

type HealthcheckTCP struct {
	HealthcheckICMP
	Port int `json:"port"`
}

type HealthcheckUDP struct {
	HealthcheckICMP
	Port int `json:"port"`
}

type HealthcheckHTTP struct {
	HealthcheckTCP
	URI          string `json:"uri"`
	ResponseCode int    `json:"responseCode"`
}
