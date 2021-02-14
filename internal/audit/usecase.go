package audit

import (
	"context"

	"github.com/gradusp/crispy/internal/model"
)

type Usecase interface {
	Create(ctx context.Context, who, what string)
	Notify(ctx context.Context, audit *model.Audit)
}
