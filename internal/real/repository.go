package real

import (
	"context"

	"github.com/gradusp/crispy/internal/model"
)

type Repository interface {
	Create(ctx context.Context, real *model.Real) (*model.Real, error)

	// GetAll(ctx context.Context) ([]*model.Real, error)
	// GetByServiceID(ctx context.Context, sid string) ([]*model.Real, error)
	// GetByAddr(ctx context.Context, a string) ([]*model.Real, error)
	GetByField(ctx context.Context, where string) ([]*model.Real, error)
	GetByID(ctx context.Context, real *model.Real) (*model.Real, error)

	Delete(ctx context.Context, real *model.Real) error
}
