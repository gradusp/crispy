package usecase

import (
	"context"

	"github.com/gradusp/crispy/internal/healthcheck"
	"github.com/gradusp/crispy/internal/model"
)

type HealthcheckUsecase struct {
	hcRepo healthcheck.Repository
}

func NewHealthcheckUsecase(r healthcheck.Repository) *HealthcheckUsecase {
	return &HealthcheckUsecase{
		hcRepo: r,
	}
}

func (hcuc HealthcheckUsecase) Create(ctx context.Context, sid string, ht, rt, ath, dth, q, h int) (*model.Healthcheck, error) {
	s := &model.Service{ID: sid}

	hc := &model.Healthcheck{
		HelloTimer:     ht,
		ResponseTimer:  rt,
		AliveThreshold: ath,
		DeadThreshold:  dth,
		Quorum:         q,
		Hysteresis:     h,
		ServiceID:      sid,
	}
	return hcuc.hcRepo.Create(ctx, s, hc)
}

func (hcuc HealthcheckUsecase) Delete(ctx context.Context, id string) error {
	hc := &model.Healthcheck{ID: id}
	return hcuc.hcRepo.Delete(ctx, hc)
}
