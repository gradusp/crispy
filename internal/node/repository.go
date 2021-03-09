package node

import (
	"context"

	"github.com/gradusp/crispy/internal/model"
)

type Repository interface {
	Create(ctx context.Context, n *model.Node) (*model.Node, error)

	GetNodeByField(ctx context.Context, where string) ([]*model.Node, error)
	GetByID(ctx context.Context, n *model.Node) (*model.Node, error)

	Delete(ctx context.Context, n *model.Node) error
}
