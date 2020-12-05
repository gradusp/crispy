package cluster

import (
	"context"
	"github.com/gradusp/crispy/ctrl/model"
)

type Usecase interface {
	// TODO: usecase should accept only raw params and wrap it to models for repo?
	Create(ctx context.Context, sz *model.SecurityZone, name string, capacity int64) (*model.Cluster, error)

	Get(ctx context.Context) ([]*model.Cluster, error)
	//GetByID(ctx context.Context, id string) (*model.Cluster, error)
	//
	//Update(ctx context.Context, sz *model.SecurityZone, id, name string) error
	//
	//Delete(ctx context.Context, id string) error
}
