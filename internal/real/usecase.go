package real

import (
	"context"
	"net"

	"github.com/gradusp/crispy/internal/model"
)

type Usecase interface {
	Create(ctx context.Context, sid string, a net.IP, p int) (*model.Real, error)

	Get(ctx context.Context, sid, a string) ([]*model.Real, error)
	GetByID(ctx context.Context, rid string) (*model.Real, error)

	Delete(ctx context.Context, rid string) error
}
