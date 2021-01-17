package cluster

import (
	"context"

	"github.com/gradusp/crispy/internal/model"
)

type Usecase interface {
	Create(ctx context.Context, zid, name string, capacity int64) (*model.Cluster, error)

	Get(ctx context.Context) ([]*model.Cluster, error)
	GetByID(ctx context.Context, id string) (*model.Cluster, error)

	Update(ctx context.Context, id, name string, capacity int64) error

	Delete(ctx context.Context, id string) error
}
