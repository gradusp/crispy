package healthcheck

import (
	"context"

	"github.com/gradusp/crispy/internal/model"
)

type Repository interface {
	Create(ctx context.Context, s *model.Service, hc *model.Healthcheck) (*model.Healthcheck, error)

	//GetByID(ctx context.Context, hc *model.Healthcheck) (*model.Healthcheck, error)
	//GetByService(ctx context.Context, s *model.Service) ([]*model.Healthcheck, error)

	//Update(ctx context.Context, hc *model.Healthcheck) error

	Delete(ctx context.Context, hc *model.Healthcheck) error
}
