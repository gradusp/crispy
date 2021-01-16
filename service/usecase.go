package service

import (
	"context"
	"net"

	"github.com/gradusp/crispy/model"
)

type Usecase interface {
	Create(ctx context.Context, cid, rt, bt, proto string, a net.IP, bw, port int) (*model.BalancingService, error)

	Get(ctx context.Context) ([]*model.BalancingService, error)
	GetByID(ctx context.Context, id string) (*model.BalancingService, error)

	Update(ctx context.Context, id string) error

	Delete(ctx context.Context, id string) error
}
