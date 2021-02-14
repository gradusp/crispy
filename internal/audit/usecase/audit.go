package usecase

import (
	"context"

	"github.com/gradusp/crispy/internal/audit"
	"github.com/gradusp/crispy/internal/model"
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

func (auc AuditUsecase) Notify(ctx context.Context, a *model.Audit) {
	auc.auditRepo.Create(ctx, a)
}
