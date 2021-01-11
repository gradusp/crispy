package usecase

import (
	"context"

	"github.com/gradusp/crispy/model"
)

func (szuc SecurityZoneUsecase) Delete(ctx context.Context, id string) error {
	sz := &model.SecurityZone{
		ID: id,
	}

	return szuc.securityZoneRepo.Delete(ctx, sz)
}
