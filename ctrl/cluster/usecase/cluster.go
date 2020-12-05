package usecase

import (
	"context"
	"fmt"
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

func (cuc ClusterUsecase) Create(ctx context.Context, sz *model.SecurityZone, name string, cap int64) (*model.Cluster, error) {
	c := &model.Cluster{
		Name: name,
	}
	fmt.Printf("CLUSTER_USECASE:24: %+v\n", c)

	return cuc.clusterRepo.Create(ctx, sz, c)
}

func (cuc ClusterUsecase) Get(ctx context.Context) ([]*model.Cluster, error) {
	panic("usecase not implemented yet")
}
