package usecase

import (
	"context"
	"github.com/gradusp/crispy/ctrl/model"
)

func (szuc SecurityZoneUsecase) Get(ctx context.Context) ([]*model.SecurityZone, error) {
	return szuc.securityZoneRepo.Get(ctx)
}

func (szuc SecurityZoneUsecase) GetByID(ctx context.Context, id string) (*model.SecurityZone, error) {
	sz := &model.SecurityZone{
		ID: id,
	}
	return szuc.securityZoneRepo.GetByID(ctx, sz)
}
