package usecase

import (
	"context"
	"net"

	"github.com/gradusp/crispy/internal/model"
	"github.com/gradusp/crispy/internal/real"
)

type RealUsecase struct {
	r real.Repository
}

func NewUsecase(r real.Repository) *RealUsecase {
	return &RealUsecase{
		r: r,
	}
}

func (ruc RealUsecase) Create(ctx context.Context, sid string, a, hca net.IP, p, hcp int) (*model.Real, error) {
	s := &model.Service{ID: sid}

	r := &model.Real{
		Addr:            a,
		Port:            p,
		HealthcheckAddr: hca,
		HealthcheckPort: hcp,
	}

	return ruc.r.Create(ctx, r)
}

func (ruc RealUsecase) Delete(ctx context.Context, id string) error {
	r := &model.Real{ID: id}

	return ruc.rRepo.Delete(ctx, r)
}
