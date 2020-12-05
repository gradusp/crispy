package usecase

import (
	"context"
	"github.com/gradusp/crispy/ctrl/cluster"
	"github.com/gradusp/crispy/ctrl/model"
)

type ClusterUsecase struct {
	clusterRepo cluster.Repository
}

func NewClusterUsecase(clusterRepo cluster.Repository) *ClusterUsecase {
	return &ClusterUsecase{
		clusterRepo: clusterRepo,
	}
}

// TODO: usecase should accept only raw params and wrap it to models for repo?
func (cuc ClusterUsecase) Create(ctx context.Context, szid, name string, capacity int64) (*model.Cluster, error) {
	sz := &model.SecurityZone{
		ID: szid,
	}

	c := &model.Cluster{
		Name:     name,
		Capacity: capacity,
	}

	return cuc.clusterRepo.Create(ctx, sz, c)
}

func (cuc ClusterUsecase) Get(ctx context.Context) ([]*model.Cluster, error) {
	panic("usecase not implemented yet")
}
