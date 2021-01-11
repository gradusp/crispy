package balancingservice

import (
	"context"
	"github.com/gradusp/crispy/ctrl/model"
	"net"
)

type Usecase interface {
	Create(ctx context.Context, cid, rt, bt, proto string, a net.IP, bw, port int) (*model.BalancingService, error)

	Get(ctx context.Context) ([]*model.BalancingService, error)
	GetByID(ctx context.Context, id string) (*model.BalancingService, error)

	Update(ctx context.Context, id string) error

	Delete(ctx context.Context, id string) error
}
