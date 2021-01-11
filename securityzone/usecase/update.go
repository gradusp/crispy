package usecase

import (
	"context"

	"github.com/gradusp/crispy/model"
)

func (szuc SecurityZoneUsecase) Update(ctx context.Context, id, name string) error {
	sz := &model.SecurityZone{
		ID:   id,
		Name: name,
	}
	return szuc.securityZoneRepo.Update(ctx, sz)
}
