package service

import (
	"context"
	"net"

	"github.com/gradusp/crispy/internal/model"
)

type Usecase interface {
	Create(ctx context.Context, cid, rt, bt, proto string, a net.IP, bw, port int) (*model.Service, error)

	Get(ctx context.Context) ([]*model.Service, error)
	GetByID(ctx context.Context, id string) (*model.Service, error)

	//Update(ctx context.Context, id string) error

	Delete(ctx context.Context, id string) error
}
