package usecase

import (
	"context"

	"github.com/gradusp/crispy/internal/model"

	"github.com/gradusp/crispy/internal/audit"
)

type AuditUsecase struct {
	auditRepo audit.Repository
}

func NewAuditUsecase(ar audit.Repository) *AuditUsecase {
	return &AuditUsecase{
		auditRepo: ar,
	}
}

func (auc AuditUsecase) Create(ctx context.Context, who, what string) {
	a := &model.Audit{
		Who:  who,
		What: what,
	}
	auc.auditRepo.Create(ctx, a)
}
