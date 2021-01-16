package zone

import (
	"context"

	"github.com/gradusp/crispy/model"
)

type Repository interface {
	Create(ctx context.Context, sz *model.Zone) (*model.Zone, error)

	Get(ctx context.Context) ([]*model.Zone, error)
	GetByID(ctx context.Context, sz *model.Zone) (*model.Zone, error)

	Update(ctx context.Context, sz *model.Zone) error

	Delete(ctx context.Context, sz *model.Zone) error
}
