package usecase

import (
	"context"
	"github.com/gradusp/crispy/ctrl/model"
)

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