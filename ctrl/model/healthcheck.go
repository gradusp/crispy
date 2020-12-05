package model

type Healthcheck struct {
	ID               string            `json:"id" pg:"id,pk"`
	HelloTimer       int               `json:"helloTimer" pg:"hello_timer"`
	ResponseTimer    int               `json:"responseTimer" pg:"response_timer"`
	AliveThreshold   int               `json:"aliveThreshold" pg:"alive_threshold"`
	DeadThreshold    int               `json:"deadThreshold" pg:"dead_threshold"`
	Quorum           int               `json:"quorum" pg:"quorum"`
	Hysteresis       int               `json:"hysteresis" pg:"hysteresis"`
	BalancingService *BalancingService `json:"balancingService" pg:"rel:has-one"`
	//HealthcheckConfig *HealthcheckConfig `json:"healthcheckConfig" pg:"rel:has-one"`
}
