package usecase

import (
	"context"

	"github.com/gradusp/crispy/model"
)

func (szuc SecurityZoneUsecase) Create(ctx context.Context, name string) (*model.SecurityZone, error) {
	sz := &model.SecurityZone{
		Name: name,
	}

	return szuc.securityZoneRepo.Create(ctx, sz)
}
