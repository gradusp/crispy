package healthcheck

import (
	"context"
	"github.com/gradusp/crispy/ctrl/model"
)

type Usecase interface {
	Create(ctx context.Context, bsid string, ht, rt, ath, dth, q, h int) (*model.Healthcheck, error)

	GetByBalancingService(ctx context.Context, bsid string) ([]*model.BalancingService, error)

	Update(ctx context.Context, id string, ht, rt, ath, dth, q, h int) error

	Delete(ctx context.Context, id string) error
}
