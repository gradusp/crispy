package usecase

import (
	"context"
	"net"

	"github.com/gradusp/crispy/internal/model"
	"github.com/gradusp/crispy/internal/service"
)

type ServiceUsecase struct {
	r service.Repository
}

func NewUsecase(r service.Repository) *ServiceUsecase {
	return &ServiceUsecase{
		r: r,
	}
}

func (suc ServiceUsecase) Create(
	ctx context.Context,
	cid, rt, bt, proto string,
	a net.IP,
	bw, port int,
) (*model.Service, error) {
	cl := &model.Cluster{ID: cid}

	bs := &model.Service{
		BalancingType: bt,
		RoutingType:   rt,
		Bandwidth:     bw,
		Proto:         proto,
		Addr:          a,
		Port:          port,
		ClusterID:     cid,
	}

	return suc.r.Create(ctx, cl, bs)
}

func (suc ServiceUsecase) Get(ctx context.Context) ([]*model.Service, error) {
	return suc.r.Get(ctx)
}

func (suc ServiceUsecase) GetByID(ctx context.Context, id string) (*model.Service, error) {
	s := &model.Service{ID: id}
	return suc.r.GetByID(ctx, s)
}

//func (suc ServiceUsecase) Update(ctx context.Context, id string) error {
//	panic("implement me")
//}

func (suc ServiceUsecase) Delete(ctx context.Context, id string) error {
	s := &model.Service{ID: id}
	return suc.r.Delete(ctx, s)
}
