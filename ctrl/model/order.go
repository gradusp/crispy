package model

import "time"

type Order struct {
	ID               string            `json:"id" pg:"id,pk"`
	OrderTypeID      string            `json:"orderType" pg:"order_type"`
	CreatedAt        time.Time         `json:"createdAt" pg:"created_at"`
	Source           string            `json:"source" pk:"source"`
	RawBody          string            `json:"rawBody" pk:"raw_body"`
	ServiceManagerID string            `json:"smId" pk:"sm_id"`
	BalancingService *BalancingService `json:"balancingService" pg:"rel:has-one"`
}
