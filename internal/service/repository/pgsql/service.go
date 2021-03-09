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
	"github.com/gradusp/crispy/internal/service"
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

func (r *PgRepo) Create(ctx context.Context, s *model.Service) (*model.Service, error) {
	c, err := r.pool.Acquire(ctx)
	if err != nil {
		r.log.Error(err)
		return nil, err
	}
	defer c.Release()

	if rowExists(ctx, c, "select id from controller.services where proto=$1 and addr=$2 and port=$3", s.Proto, s.Addr.To4(), s.Port) {
		err := c.QueryRow(ctx, "select id from controller.services where proto=$1 and addr=$2 and port=$3;", s.Proto, s.Addr.To4(), s.Port).Scan(&s.ID) //nolint:lll
		if err != nil {
			r.log.Error(err) // TODO: better error handling
			return nil, err
		}
		return s, service.ErrAlreadyExist
	}

	query := `
insert into controller.services (cluster_id, routing_type, balancing_type, bandwidth, proto, addr, port)
values ($1, $2, $3, $4, $5, $6, $7)
returning id;`
	if err := c.
		QueryRow(ctx, query, s.ClusterID, s.RoutingType, s.BalancingType, s.Bandwidth, s.Proto, s.Addr.To4(), s.Port).
		Scan(&s.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch {
			case pgErr.Code == "23505":
				r.log.Debugw("service already exist",
					"error_body", pgErr.Message,
					"error_code", pgErr.Code)
				return s, service.ErrAlreadyExist
			default:
				r.log.Errorw("issue with service on create",
					"error_body", pgErr.Message,
					"error_code", pgErr.Code)
				return nil, err
			}
		}
		r.log.Error(err)
		return nil, err
	}
	return s, nil
}

func (r *PgRepo) Get(ctx context.Context) ([]*model.Service, error) {
	var services []*model.Service

	c, err := r.pool.Acquire(ctx)
	if err != nil {
		r.log.Error(err)
		return nil, err
	}
	defer c.Release()

	q := `select id, cluster_id, routing_type, balancing_type, bandwidth, proto, addr, port from controller.services`
	rows, err := c.Query(ctx, q)
	if err != nil {
		r.log.Error(err)
		return nil, err
	}
	for rows.Next() {
		var s model.Service
		err = rows.Scan(&s.ID, &s.ClusterID, &s.RoutingType, &s.BalancingType, &s.Bandwidth, &s.Proto, &s.Addr, &s.Port)
		if err != nil {
			r.log.Error(err)
			return nil, err
		}
		services = append(services, &s)
	}
	err = rows.Err()

	return services, err
}

func (r *PgRepo) GetByID(ctx context.Context, s *model.Service) (*model.Service, error) {
	c, err := r.pool.Acquire(ctx)
	if err != nil {
		r.log.Error(err)
		return nil, err
	}
	defer c.Release()

	query := `select cluster_id, routing_type, balancing_type, bandwidth, proto, addr, port from controller.services where id=$1;`
	if err = c.QueryRow(ctx, query, s.ID).Scan(&s.ClusterID, &s.RoutingType, &s.BalancingType, &s.Bandwidth, &s.Proto, &s.Addr, &s.Port); err != nil { //nolint:lll
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, service.ErrNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			r.log.Errorw("can't get service",
				"error_body", pgErr.Message,
				"error_code", pgErr.Code,
			)
			return nil, service.ErrNotFound
		}
		return nil, err
	}
	return s, nil
}

// TODO: implement
//func (r *PgRepo) Update(ctx context.Context) error {
//	panic("implement my repo")
//}

func (r *PgRepo) Delete(ctx context.Context, s *model.Service) error {
	c, err := r.pool.Acquire(ctx)
	if err != nil {
		r.log.Error(err)
		return err
	}
	defer c.Release()

	res, err := c.Exec(ctx, "delete from controller.services where id=$1", s.ID)
	if err != nil {
		r.log.Error(err)
		return err
	}
	if res.RowsAffected() != 1 {
		r.log.Debug("delete for non-existing Service ID")
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
