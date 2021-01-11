package usecase

import (
	"context"
	"fmt"
	"net"

	"github.com/gradusp/crispy/ctrl/balancingservice"
	"github.com/gradusp/crispy/ctrl/model"
)

type BalancingserviceUsecase struct {
	balancingserviceRepo balancingservice.Repository
}

func NewBalancingserviceUsecase(balancingserviceRepo balancingservice.Repository) *BalancingserviceUsecase {
	return &BalancingserviceUsecase{
		balancingserviceRepo: balancingserviceRepo,
	}
}

func (bsuc BalancingserviceUsecase) Create(
	ctx context.Context,
	cid, rt, bt, proto string,
	a net.IP,
	bw, port int,
) (*model.BalancingService, error) {
	bs := &model.BalancingService{
		BalancingType: bt,
		RoutingType:   rt,
		Bandwidth:     bw,
		Proto:         proto,
		Addr:          a,
		Port:          port,
		ClusterID:     cid,
	}
	fmt.Printf("USECASE:20:CreateBalancingservice: %+v\n", bs) // FIXME

	return nil, nil
}

func (bsuc BalancingserviceUsecase) Get(ctx context.Context) ([]*model.BalancingService, error) {
	return nil, nil
	//panic("implement me")
}

func (bsuc BalancingserviceUsecase) GetByID(ctx context.Context, id string) (*model.BalancingService, error) {
	panic("implement me")
}

func (bsuc BalancingserviceUsecase) Update(ctx context.Context, id string) error {
	panic("implement me")
}

func (bsuc BalancingserviceUsecase) Delete(ctx context.Context, id string) error {
	panic("implement me")
}
