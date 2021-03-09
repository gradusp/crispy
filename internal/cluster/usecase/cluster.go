package usecase

import (
	"context"

	"github.com/gradusp/crispy/internal/cluster"
	"github.com/gradusp/crispy/internal/model"
)

type ClusterUsecase struct {
	r cluster.Repository
}

func NewUsecase(r cluster.Repository) *ClusterUsecase {
	return &ClusterUsecase{
		r: r,
	}
}

func (cuc ClusterUsecase) Create(ctx context.Context, zid, name string, capacity int64) (*model.Cluster, error) {
	z := &model.Zone{
		ID: zid,
	}

	c := &model.Cluster{
		Name:     name,
		Capacity: capacity,
	}

	return cuc.r.Create(ctx, z, c)
}

func (cuc ClusterUsecase) Get(ctx context.Context) ([]*model.Cluster, error) {
	return cuc.r.Get(ctx)
}

func (cuc ClusterUsecase) GetByID(ctx context.Context, id string) (*model.Cluster, error) {
	c := &model.Cluster{
		ID: id,
	}
	return cuc.r.GetByID(ctx, c)
}

func (cuc ClusterUsecase) Update(ctx context.Context, id, name string, capacity int64) error {
	c := &model.Cluster{
		ID:       id,
		Name:     name,
		Capacity: capacity,
	}
	return cuc.r.Update(ctx, c)
}

func (cuc ClusterUsecase) Delete(ctx context.Context, id string) error {
	c := &model.Cluster{
		ID: id,
	}
	return cuc.r.Delete(ctx, c)
}
