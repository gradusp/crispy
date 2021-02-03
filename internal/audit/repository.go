package audit

import (
	"context"

	"github.com/gradusp/crispy/internal/model"
)

type Repository interface {
	Create(ctx context.Context, a *model.Audit)
}
