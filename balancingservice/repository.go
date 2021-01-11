package balancingservice

import (
	"context"

	"github.com/gradusp/crispy/model"
)

type Repository interface {
	Create(ctx context.Context) (*model.BalancingService, error)

	Get(ctx context.Context) ([]*model.BalancingService, error)
	GetByID(ctx context.Context) (*model.BalancingService, error)

	Update(ctx context.Context) error

	Delete(ctx context.Context) error
}
