package usecase

import (
	"context"

	"github.com/gradusp/crispy/model"
)

func (cuc ClusterUsecase) Delete(ctx context.Context, id string) error {
	c := &model.Cluster{
		ID: id,
	}
	return cuc.clusterRepo.Delete(ctx, c)
}
