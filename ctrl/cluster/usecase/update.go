package usecase

import (
	"context"
	"github.com/gradusp/crispy/ctrl/model"
)

func (cuc ClusterUsecase) Update(ctx context.Context, szid, id, name string, capacity int64) error {
	sz := &model.SecurityZone{
		ID: szid,
	}

	c := &model.Cluster{
		ID:       id,
		Name:     name,
		Capacity: capacity,
	}

	return cuc.clusterRepo.Update(ctx, sz, c)
}
