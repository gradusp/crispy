package usecase

import (
	"context"

	"github.com/gradusp/crispy/model"
)

func (cuc ClusterUsecase) Get(ctx context.Context) ([]*model.Cluster, error) {
	return cuc.clusterRepo.Get(ctx)
}

func (cuc ClusterUsecase) GetByID(ctx context.Context, id string) (*model.Cluster, error) {
	c := &model.Cluster{
		ID: id,
	}
	return cuc.clusterRepo.GetByID(ctx, c)
}
