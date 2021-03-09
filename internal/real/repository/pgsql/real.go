package pgsql

import (
	"context"
	"errors"
	"fmt"

	"github.com/gradusp/crispy/internal/model"
	real "github.com/gradusp/crispy/internal/real"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type RealPostgresRepo struct {
	log  *zap.SugaredLogger
	pool *pgxpool.Pool
}

func NewPgRepo(pool *pgxpool.Pool, l *zap.SugaredLogger) *RealPostgresRepo {
	return &RealPostgresRepo{
		log:  l,
		pool: pool,
	}
}

func (rr RealPostgresRepo) Create(ctx context.Context, s *model.Service, r *model.Real) (*model.Real, error) {
	c, err := rr.pool.Acquire(ctx)
	if err != nil {
		rr.log.Error(err)
		return nil, err
	}
	defer c.Release()

	// TODO: advanced error checks

	query := `
insert into controller.reals (addr, port, hc_addr, hc_port, service_id)
values ($1, $2, $3, $4, $5)
returning id;`

	if err := c.QueryRow(ctx, query, r.Addr, r.Port, r.HealthcheckAddr, r.HealthcheckPort, s.ID).Scan(&r.ID); err != nil {
		rr.log.Error(err)
		return nil, err
	}

	return r, nil
}

func (rr RealPostgresRepo) Delete(ctx context.Context, r *model.Real) error {
	c, err := rr.pool.Acquire(ctx)
	if err != nil {
		rr.log.Error(err)
		return err
	}
	defer c.Release()

	// TODO: advanced error checks

	res, err := c.Exec(ctx, "delete from controller.reals where id=$1", r.ID)
	if err != nil {
		rr.log.Error(err)
		return err
	}
	if res.RowsAffected() != 1 {
		rr.log.Debug("delete for non-existing Real ID")
		return nil
	}
	return nil
}
