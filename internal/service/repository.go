package service

import (
	"context"

	"github.com/gradusp/crispy/internal/model"
)

type Repository interface {
	Create(ctx context.Context, cl *model.Cluster, s *model.Service) (*model.Service, error)

	Get(ctx context.Context) ([]*model.Service, error)
	GetByID(ctx context.Context, s *model.Service) (*model.Service, error)

	//Update(ctx context.Context) error

	Delete(ctx context.Context, s *model.Service) error
}
