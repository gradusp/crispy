package model

type Healthcheck struct {
	ID             string   `json:"id"`
	HelloTimer     int      `json:"helloTimer"`
	ResponseTimer  int      `json:"responseTimer"`
	AliveThreshold int      `json:"aliveThreshold"`
	DeadThreshold  int      `json:"deadThreshold"`
	Quorum         int      `json:"quorum"`
	Hysteresis     int      `json:"hysteresis"`
	ServiceID      string   `json:"serviceId,omitempty"`
	Service        *Service `json:"service,omitempty"`
}
