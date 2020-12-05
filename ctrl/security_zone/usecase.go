package security_zone

import (
	"context"
	"github.com/gradusp/crispy/ctrl/model"
)

type Usecase interface {
	Create(ctx context.Context, name string) (*model.SecurityZone, error)

	Get(ctx context.Context) ([]*model.SecurityZone, error)

	Update(ctx context.Context, id, name string) error

	Delete(ctx context.Context, id string) error
}
