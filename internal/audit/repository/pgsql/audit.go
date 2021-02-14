package pgsql

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"

	"github.com/gradusp/crispy/internal/model"
)

type AuditRepo struct {
	pool *pgxpool.Pool
	log  *zap.SugaredLogger
}

func NewAuditRepo(pool *pgxpool.Pool, log *zap.SugaredLogger) *AuditRepo {
	return &AuditRepo{
		pool: pool,
		log:  log,
	}
}

func (ar *AuditRepo) Create(ctx context.Context, a *model.Audit) {
	c, err := ar.pool.Acquire(ctx)
	if err != nil {
		ar.log.Error(err)
	}
	defer c.Release()

	query := `insert into controller.audit (entity, action, who, what) values ($1, $2, $3, $4) returning id;`
	if err := c.QueryRow(ctx, query, a.Entity, a.Action, a.Who, a.What).Scan(&a.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			ar.log.Errorw("issue with audit on create", "error_body", pgErr.Message, "error_code", pgErr.Code)
		}
		ar.log.Error(err)
	}
}
