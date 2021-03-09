package usecase

import (
	"context"
	"fmt"
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

func (ruc RealUsecase) Create(ctx context.Context, sid string, a net.IP, p int) (*model.Real, error) {
	r := &model.Real{
		Addr:      a,
		Port:      p,
		ServiceID: sid,
	}

	return ruc.r.Create(ctx, r)
}

// TODO: implement checks for SQL INJECTIONS (regexp?)
// TODO: check 'sid' for UUID
// TODO: check 'a' for net.IP
func (ruc RealUsecase) Get(ctx context.Context, sid, a string) ([]*model.Real, error) {
	switch {
	case sid != "" && a != "":
		return nil, real.ErrWrongQuery
	case sid != "":
		q := fmt.Sprintf("where service_id='%s'", sid)
		return ruc.r.GetByField(ctx, q)
	case a != "":
		q := fmt.Sprintf("where addr='%s'", net.ParseIP(a))
		return ruc.r.GetByField(ctx, q)
	default:
		return ruc.r.GetByField(ctx, "")
	}
}

func (ruc RealUsecase) GetByID(ctx context.Context, rid string) (*model.Real, error) {
	r := &model.Real{ID: rid}
	return ruc.r.GetByID(ctx, r)
}

func (ruc RealUsecase) Delete(ctx context.Context, rid string) error {
	r := &model.Real{ID: rid}
	return ruc.r.Delete(ctx, r)
}
