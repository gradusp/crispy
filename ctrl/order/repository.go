package order

import (
	"context"
	"github.com/gradusp/crispy/ctrl/models"
)

type Repository interface {
	CreateOrder(ctx context.Context, bs *models.BalancingService, user *models.Order) error
	GetOrders(ctx context.Context) ([]*models.Order, error)
	DeleteOrder(ctx context.Context, id string) error
}
