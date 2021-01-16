package zone

import (
	"context"

	"github.com/gradusp/crispy/model"
)

type Usecase interface {
	Create(ctx context.Context, name string) (*model.Zone, error)

	Get(ctx context.Context) ([]*model.Zone, error)
	GetByID(ctx context.Context, id string) (*model.Zone, error)

	Update(ctx context.Context, id, name string) error

	Delete(ctx context.Context, id string) error
}
