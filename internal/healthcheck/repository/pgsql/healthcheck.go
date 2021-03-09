package pgsql

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"

	"github.com/gradusp/crispy/internal/model"
)

type HealthcheckPostgresRepo struct {
	log  *zap.SugaredLogger
	pool *pgxpool.Pool
}

func NewPgRepo(pool *pgxpool.Pool, l *zap.SugaredLogger) *HealthcheckPostgresRepo {
	return &HealthcheckPostgresRepo{
		log:  l,
		pool: pool,
	}
}

// TODO: refactor Healthcheck model since right now it is only unique by PK (id)

func (hcr *HealthcheckPostgresRepo) Create(ctx context.Context, s *model.Service, hc *model.Healthcheck) (*model.Healthcheck, error) {
	c, err := hcr.pool.Acquire(ctx)
	if err != nil {
		hcr.log.Error(err)
		return nil, err
	}
	defer c.Release()

	// TODO: advanced error checks

	query := `
insert into controller.healthchecks (hello_timer, response_timer, alive_threshold, dead_threshold, quorum, hysteresis,
                                     service_id)
values ($1, $2, $3, $4, $5, $6, $7)
returning id;`

	if err := c.QueryRow(ctx, query,
		hc.HelloTimer,
		hc.ResponseTimer,
		hc.AliveThreshold,
		hc.DeadThreshold,
		hc.Quorum,
		hc.Hysteresis,
		s.ID).
		Scan(&hc.ID); err != nil {
		hcr.log.Error(err)
		return nil, err
	}

	return hc, nil
}

func (hcr *HealthcheckPostgresRepo) Delete(ctx context.Context, hc *model.Healthcheck) error {
	c, err := hcr.pool.Acquire(ctx)
	if err != nil {
		hcr.log.Error(err)
		return err
	}
	defer c.Release()

	// TODO: advanced error checks

	r, err := c.Exec(ctx, "delete from controller.healthchecks where id=$1", hc.ID)
	if err != nil {
		hcr.log.Error(err)
		return err
	}
	if r.RowsAffected() != 1 {
		hcr.log.Debug("delete for non-existing Healthcheck ID")
		return nil
	}
	return nil
}
