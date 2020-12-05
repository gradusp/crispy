package security_zone

import (
	"context"
	"github.com/gradusp/crispy/ctrl/model"
)

type Repository interface {
	Create(ctx context.Context, sz *model.SecurityZone) (*model.SecurityZone, error)

	Get(ctx context.Context) ([]*model.SecurityZone, error)
	GetByID(ctx context.Context, sz *model.SecurityZone) (*model.SecurityZone, error)

	Update(ctx context.Context, sz *model.SecurityZone) error

	Delete(ctx context.Context, sz *model.SecurityZone) error
}
