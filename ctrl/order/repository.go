package order

import (
	"context"
	"github.com/gradusp/crispy/ctrl/model"
)

type Repository interface {
	CreateOrder(ctx context.Context, bs *model.BalancingService, user *model.Order) error
	GetOrders(ctx context.Context) ([]*model.Order, error)
	DeleteOrder(ctx context.Context, id string) error
}
