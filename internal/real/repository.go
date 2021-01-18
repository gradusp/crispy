package real

import (
	"context"

	"github.com/gradusp/crispy/internal/model"
)

type Repository interface {
	Create(ctx context.Context, s *model.Service, r *model.Real) (*model.Real, error)

	//GetByID(ctx context.Context, hc *model.Real) (*model.Real, error)
	//GetByService(ctx context.Context, s *model.Service) ([]*model.Real, error)

	//Update(ctx context.Context, r *model.Real) error

	Delete(ctx context.Context, r *model.Real) error
}
