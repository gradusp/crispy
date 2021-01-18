package usecase

import (
	"context"
	"net"

	"github.com/gradusp/crispy/internal/model"
	"github.com/gradusp/crispy/internal/service"
)

type ServiceUsecase struct {
	serviceRepo service.Repository
}

func NewServiceUsecase(sr service.Repository) *ServiceUsecase {
	return &ServiceUsecase{
		serviceRepo: sr,
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

	return suc.serviceRepo.Create(ctx, cl, bs)
}

func (suc ServiceUsecase) Get(ctx context.Context) ([]*model.Service, error) {
	return suc.serviceRepo.Get(ctx)
}

func (suc ServiceUsecase) GetByID(ctx context.Context, id string) (*model.Service, error) {
	s := &model.Service{ID: id}
	return suc.serviceRepo.GetByID(ctx, s)
}

//func (suc ServiceUsecase) Update(ctx context.Context, id string) error {
//	panic("implement me")
//}

func (suc ServiceUsecase) Delete(ctx context.Context, id string) error {
	s := &model.Service{ID: id}
	return suc.serviceRepo.Delete(ctx, s)
}
