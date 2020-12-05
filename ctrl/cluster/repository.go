package cluster

import (
	"context"
	"github.com/gradusp/crispy/ctrl/model"
)

type Repository interface {
	Create(ctx context.Context, sz *model.SecurityZone, c *model.Cluster) (*model.Cluster, error)

	Get(ctx context.Context) ([]*model.Cluster, error)
	GetByID(ctx context.Context, c *model.Cluster) (*model.Cluster, error)

	Update(ctx context.Context, sz *model.SecurityZone, c *model.Cluster) error

	Delete(ctx context.Context, c *model.Cluster) error
}
