package usecase

import (
	"context"

	"github.com/gradusp/crispy/internal/cluster"
	"github.com/gradusp/crispy/internal/model"
)

type ClusterUsecase struct {
	clusterRepo cluster.Repository
}

func NewClusterUsecase(clusterRepo cluster.Repository) *ClusterUsecase {
	return &ClusterUsecase{
		clusterRepo: clusterRepo,
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

	return cuc.clusterRepo.Create(ctx, z, c)
}

func (cuc ClusterUsecase) Get(ctx context.Context) ([]*model.Cluster, error) {
	return cuc.clusterRepo.Get(ctx)
}

func (cuc ClusterUsecase) GetByID(ctx context.Context, id string) (*model.Cluster, error) {
	c := &model.Cluster{
		ID: id,
	}
	return cuc.clusterRepo.GetByID(ctx, c)
}

func (cuc ClusterUsecase) Update(ctx context.Context, id, name string, capacity int64) error {
	c := &model.Cluster{
		ID:       id,
		Name:     name,
		Capacity: capacity,
	}
	return cuc.clusterRepo.Update(ctx, c)
}

func (cuc ClusterUsecase) Delete(ctx context.Context, id string) error {
	c := &model.Cluster{
		ID: id,
	}
	return cuc.clusterRepo.Delete(ctx, c)
}
