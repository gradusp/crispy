package usecase

import (
	"context"
	"github.com/gradusp/crispy/ctrl/model"
	"github.com/gradusp/crispy/ctrl/security_zone"
)

type SecurityZoneUsecase struct {
	securityZoneRepo security_zone.Repository
}

func NewSecurityZoneUseCase(securityZoneRepo security_zone.Repository) *SecurityZoneUsecase {
	return &SecurityZoneUsecase{
		securityZoneRepo: securityZoneRepo,
	}
}

func (szuc SecurityZoneUsecase) Create(ctx context.Context, name string) (*model.SecurityZone, error) {
	sz := &model.SecurityZone{
		Name: name,
	}

	return szuc.securityZoneRepo.Create(ctx, sz)
}

func (szuc SecurityZoneUsecase) Get(ctx context.Context) ([]*model.SecurityZone, error) {
	return szuc.securityZoneRepo.Get(ctx)
}

func (szuc SecurityZoneUsecase) Update(ctx context.Context, id, name string) error {
	sz := &model.SecurityZone{
		ID:   id,
		Name: name,
	}

	return szuc.securityZoneRepo.Update(ctx, sz)
}

func (szuc SecurityZoneUsecase) Delete(ctx context.Context, id string) error {
	sz := &model.SecurityZone{
		ID:   id,
	}

	return szuc.securityZoneRepo.Delete(ctx, sz)
}
