package real

import (
	"context"
	"net"

	"github.com/gradusp/crispy/internal/model"
)

type Usecase interface {
	Create(ctx context.Context, sid string, a, hca net.IP, p, hcp int) (*model.Real, error)

	//GetByBalancingService(ctx context.Context, id string) ([]*model.Real, error)

	//UpdateHealthcheckAddress(ctx context.Context, id, hca string) error

	Delete(ctx context.Context, id string) error
}
