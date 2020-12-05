package cluster

import (
	"context"
	"github.com/gradusp/crispy/ctrl/model"
)

type Usecase interface {
	Create(ctx context.Context, sz *model.SecurityZone, name string, cap int64) (*model.Cluster, error)

	Get(ctx context.Context) ([]*model.Cluster, error)
	//GetByID(ctx context.Context, id string) (*model.Cluster, error)
	//
	//Update(ctx context.Context, sz *model.SecurityZone, id, name string) error
	//
	//Delete(ctx context.Context, id string) error
}
