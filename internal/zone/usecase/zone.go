package usecase

import (
	"context"

	"github.com/gradusp/crispy/internal/model"
	"github.com/gradusp/crispy/internal/zone"
)

type ZoneUsecase struct {
	zoneRepo zone.Repository
}

func NewZoneUsecase(zoneRepo zone.Repository) *ZoneUsecase {
	return &ZoneUsecase{
		zoneRepo: zoneRepo,
	}
}

func (zuc ZoneUsecase) Create(ctx context.Context, name string) (*model.Zone, error) {
	z := &model.Zone{
		Name: name,
	}
	return zuc.zoneRepo.Create(ctx, z)
}

func (zuc ZoneUsecase) Get(ctx context.Context) ([]*model.Zone, error) {
	return zuc.zoneRepo.Get(ctx)
}

func (zuc ZoneUsecase) GetByID(ctx context.Context, id string) (*model.Zone, error) {
	sz := &model.Zone{
		ID: id,
	}
	return zuc.zoneRepo.GetByID(ctx, sz)
}

func (zuc ZoneUsecase) Update(ctx context.Context, id, name string) error {
	sz := &model.Zone{
		ID:   id,
		Name: name,
	}
	return zuc.zoneRepo.Update(ctx, sz)
}

func (zuc ZoneUsecase) Delete(ctx context.Context, id string) error {
	sz := &model.Zone{
		ID: id,
	}
	return zuc.zoneRepo.Delete(ctx, sz)
}
