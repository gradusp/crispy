package usecase

import (
	"context"

	"github.com/gradusp/crispy/internal/model"
	"github.com/gradusp/crispy/internal/zone"
)

type ZoneUsecase struct {
	zoneRepo zone.Repository
}

func NewZoneUseCase(zoneRepo zone.Repository) *ZoneUsecase {
	return &ZoneUsecase{
		zoneRepo: zoneRepo,
	}
}

func (szuc ZoneUsecase) Create(ctx context.Context, name string) (*model.Zone, error) {
	sz := &model.Zone{
		Name: name,
	}
	return szuc.zoneRepo.Create(ctx, sz)
}

func (szuc ZoneUsecase) Get(ctx context.Context) ([]*model.Zone, error) {
	return szuc.zoneRepo.Get(ctx)
}

func (szuc ZoneUsecase) GetByID(ctx context.Context, id string) (*model.Zone, error) {
	sz := &model.Zone{
		ID: id,
	}
	return szuc.zoneRepo.GetByID(ctx, sz)
}

func (szuc ZoneUsecase) Update(ctx context.Context, id, name string) error {
	sz := &model.Zone{
		ID:   id,
		Name: name,
	}
	return szuc.zoneRepo.Update(ctx, sz)
}

func (szuc ZoneUsecase) Delete(ctx context.Context, id string) error {
	sz := &model.Zone{
		ID: id,
	}
	return szuc.zoneRepo.Delete(ctx, sz)
}
