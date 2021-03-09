package pgsql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"

	"github.com/gradusp/crispy/internal/model"
	real "github.com/gradusp/crispy/internal/real"
)

type PgRepo struct {
	log  *zap.SugaredLogger
	pool *pgxpool.Pool
}

func NewPgRepo(pool *pgxpool.Pool, l *zap.SugaredLogger) *PgRepo {
	return &PgRepo{
		log:  l,
		pool: pool,
	}
}

func (r PgRepo) Create(ctx context.Context, rl *model.Real) (*model.Real, error) {
	c, err := r.pool.Acquire(ctx)
	if err != nil {
		r.log.Error(err)
		return nil, err
	}
	defer c.Release()

	if rowExists(ctx, c, "select id from controller.reals where addr=$1 and port=$2", rl.Addr, rl.Port) {
		err := c.QueryRow(ctx, "select id from controller.reals where addr=$1 and port=$2;", rl.Addr, rl.Port).Scan(&rl.ID)
		if err != nil {
			r.log.Error(err)
			return nil, err
		}
		return nil, real.ErrAlreadyExist
	}

	query := `insert into controller.reals (service_id, addr, port) values ($1, $2, $3) returning id;`
	if err := c.QueryRow(ctx, query, rl.ServiceID, rl.Addr.String(), rl.Port).Scan(&rl.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch {
			// FIXME: switch is not needed with only default but can be useful for quick handle of pgsql errors
			default:
				r.log.Errorw("issue with real on create",
					"error_body", pgErr.Message,
					"error_code", pgErr.Code,
				)
				return nil, err
			}
		}
		r.log.Error(err)
		return nil, err
	}

	return rl, nil
}

func (r PgRepo) GetByField(ctx context.Context, where string) ([]*model.Real, error) {
	var reals []*model.Real
	c, err := r.pool.Acquire(ctx)
	if err != nil {
		r.log.Error(err)
		return nil, err
	}
	defer c.Release()

	q := fmt.Sprintf("select id, service_id, addr, port from controller.reals %s;", where)
	rows, err := c.Query(ctx, q)
	if err != nil {
		r.log.Error("reals can't be selected ", err) // TODO: better error handling
		return nil, err
	}

	for rows.Next() {
		var rl model.Real
		err = rows.Scan(&rl.ID, &rl.ServiceID, &rl.Addr, &rl.Port)
		if err != nil {
			return nil, err
		}
		reals = append(reals, &rl)
	}

	return reals, nil
}

func (r PgRepo) GetByID(ctx context.Context, rl *model.Real) (*model.Real, error) {
	c, err := r.pool.Acquire(ctx)
	if err != nil {
		r.log.Error(err)
		return nil, err
	}
	defer c.Release()

	query := `select id, service_id, addr, port from controller.reals where id=$1;`
	if err = c.QueryRow(ctx, query, rl.ID).Scan(&rl.ID, &rl.ServiceID, &rl.Addr, &rl.Port); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, real.ErrNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			r.log.Errorw("can't get real",
				"error_body", pgErr.Message,
				"error_code", pgErr.Code,
			)
			return nil, real.ErrNotFound
		}
		return nil, err
	}
	return rl, nil
}

func (r PgRepo) Delete(ctx context.Context, rl *model.Real) error {
	c, err := r.pool.Acquire(ctx)
	if err != nil {
		r.log.Error(err)
		return err
	}
	defer c.Release()

	// TODO: advanced error checks

	res, err := c.Exec(ctx, "delete from controller.reals where id=$1", rl.ID)
	if err != nil {
		r.log.Error(err)
		return err
	}
	if res.RowsAffected() != 1 {
		r.log.Debug("delete for non-existing Real ID")
		return nil
	}

	return nil
}

// TODO: make it DRY (right now it repeats in every repo)
func rowExists(ctx context.Context, c *pgxpool.Conn, q string, args ...interface{}) bool {
	var exists bool

	query := fmt.Sprintf("select exists (%s)", q)
	err := c.QueryRow(ctx, query, args...).Scan(&exists)
	if err != nil && err != pgx.ErrNoRows {
		panic(err)
	}
	return exists
}
