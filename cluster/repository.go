package cluster

import (
	"context"

	"github.com/gradusp/crispy/model"
)

type Repository interface {
	Create(ctx context.Context, zone *model.Zone, cluster *model.Cluster) (*model.Cluster, error)

	Get(ctx context.Context) ([]*model.Cluster, error)
	GetByID(ctx context.Context, cluster *model.Cluster) (*model.Cluster, error)

	Update(ctx context.Context, cluster *model.Cluster) error

	Delete(ctx context.Context, cluster *model.Cluster) error
}
