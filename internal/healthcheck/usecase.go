package healthcheck

import (
	"context"

	"github.com/gradusp/crispy/internal/model"
)

type Usecase interface {
	Create(ctx context.Context, sid string, ht, rt, ath, dth, q, h int) (*model.Healthcheck, error)

	//GetByID(ctx context.Context, id string) (*model.Healthcheck, error)
	//GetByService(ctx context.Context, sid string) ([]*model.Healthcheck, error)

	//Update(ctx context.Context, id string, ht, rt, ath, dth, q, h int) error

	Delete(ctx context.Context, id string) error
}
