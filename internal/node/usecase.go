package node

import (
	"context"
	"net"

	"github.com/gradusp/crispy/internal/model"
)

type Usecase interface {
	Create(ctx context.Context, cid, h string, a net.IP) (*model.Node, error)

	Get(ctx context.Context, cid, a string) ([]*model.Node, error)
	GetByID(ctx context.Context, nid int) (*model.Node, error)

	Delete(ctx context.Context, nid int) error
}
