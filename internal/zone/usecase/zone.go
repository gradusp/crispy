package usecase

import (
	"context"

	"github.com/gradusp/crispy/internal/model"
	"github.com/gradusp/crispy/internal/zone"
)

type ZoneUsecase struct {
	r zone.Repository
}

func NewUsecase(r zone.Repository) *ZoneUsecase {
	return &ZoneUsecase{
		r: r,
	}
}

func (zuc ZoneUsecase) Create(ctx context.Context, name string) (*model.Zone, error) {
	z := &model.Zone{
		Name: name,
	}
	return zuc.r.Create(ctx, z)
}

func (zuc ZoneUsecase) Get(ctx context.Context) ([]*model.Zone, error) {
	return zuc.r.Get(ctx)
}

func (zuc ZoneUsecase) GetByID(ctx context.Context, id string) (*model.Zone, error) {
	sz := &model.Zone{
		ID: id,
	}
	return zuc.r.GetByID(ctx, sz)
}

func (zuc ZoneUsecase) Update(ctx context.Context, id, name string) error {
	sz := &model.Zone{
		ID:   id,
		Name: name,
	}
	return zuc.r.Update(ctx, sz)
}

func (zuc ZoneUsecase) Delete(ctx context.Context, id string) error {
	sz := &model.Zone{
		ID: id,
	}
	return zuc.r.Delete(ctx, sz)
}
