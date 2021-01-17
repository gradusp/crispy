package pgsql

import (
	"context"
	"errors"

	"github.com/hashicorp/consul/api"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"

	"github.com/gradusp/crispy/cluster"
	"github.com/gradusp/crispy/model"
)

type ClusterRepo struct {
	pool *pgxpool.Pool
	kv   *api.KV
	log  *zap.SugaredLogger
}

func NewClusterRepo(pool *pgxpool.Pool, kv *api.KV, l *zap.SugaredLogger) *ClusterRepo {
	return &ClusterRepo{
		pool: pool,
		kv:   kv,
		log:  l,
	}
}

func (cr *ClusterRepo) Create(ctx context.Context, sz *model.Zone, cl *model.Cluster) (*model.Cluster, error) {
	c, err := cr.pool.Acquire(ctx)
	if err != nil {
		cr.log.Error(err)
		return nil, err
	}
	defer c.Release()

	// @khannz: decided to not implement existence check here,
	// since name constraint @ DB would give feedback

	// TODO: missing trace logs here?
	// TODO: error case `invalid zone_id`

	query := `insert into controller.clusters (name, capacity, zone_id) values ($1, $2, $3) returning id;`
	if err := c.QueryRow(ctx, query, cl.Name, cl.Capacity, sz.ID).Scan(&cl.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch {
			case pgErr.Code == "23505":
				cr.log.Warnw("issue with cluster on create",
					"error_body", pgErr.Message,
					"error_code", pgErr.Code,
				)
				return cl, cluster.ErrAlreadyExist
			default:
				cr.log.Errorw("issue with cluster on create",
					"error_body", pgErr.Message,
					"error_code", pgErr.Code,
				)
				return nil, err
			}
		}
		cr.log.Error(err)
		return nil, err
	}

	// since we already got zone_id from usecase, we can just query name for the whole object
	if err := c.QueryRow(ctx, "select name from controller.zones where id=$1;", sz.ID).Scan(&sz.Name); err != nil {
		// TODO: make it more complex
		// (not so critical right after previous block)
		cr.log.Error(err)
		return nil, err
	}
	cl.Zone = sz
	return cl, nil
}

func (cr *ClusterRepo) Get(ctx context.Context) ([]*model.Cluster, error) {
	c, err := cr.pool.Acquire(ctx)
	if err != nil {
		cr.log.Error(err)
		return nil, err
	}
	defer c.Release()

	query := `select c.id, c.name, c.capacity, coalesce(sum(s.bandwidth),0), z.id, z.name
from controller.clusters c
         left join controller.services s on c.id = s.cluster_id
         left join controller.zones z on c.zone_id = z.id
group by c.id, z.id;`
	clusters, err := c.Query(ctx, query)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			cr.log.Errorw("can't get clusters",
				"error_body", pgErr.Message,
				"error_code", pgErr.Code,
			)
			return nil, cluster.ErrNotFound
		}
		cr.log.Error("clusters can't be selected", err)
		return nil, err
	}

	var r []*model.Cluster
	for clusters.Next() {
		var (
			c model.Cluster
			z model.Zone
		)
		err = clusters.Scan(&c.ID, &c.Name, &c.Capacity, &c.Usage, &z.ID, &z.Name)
		if err != nil {
			cr.log.Error(err)
			return nil, err
		}
		c.Zone = &z
		r = append(r, &c)
	}
	err = clusters.Err()

	return r, err
}

func (cr *ClusterRepo) GetByID(ctx context.Context, cl *model.Cluster) (*model.Cluster, error) {
	c, err := cr.pool.Acquire(ctx)
	if err != nil {
		cr.log.Error(err)
		return nil, err
	}
	defer c.Release()

	query := `select c.id, c.name, c.capacity, coalesce(sum(s.bandwidth),0), z.id, z.name
from controller.clusters c
         left join controller.services s on c.id = s.cluster_id
         left join controller.zones z on c.zone_id = z.id
where c.id=$1
group by c.id, z.id;`

	var z model.Zone
	if err = c.QueryRow(ctx, query, cl.ID).Scan(&cl.ID, &cl.Name, &cl.Capacity, &cl.Usage, &z.ID, &z.Name); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, cluster.ErrNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			cr.log.Errorw("can't get zone",
				"error_body", pgErr.Message,
				"error_code", pgErr.Code,
			)
			return nil, cluster.ErrNotFound
		}
		return nil, err
	}
	cl.Zone = &z
	return cl, err
}

func (cr *ClusterRepo) Update(ctx context.Context, cl *model.Cluster) error {
	c, err := cr.pool.Acquire(ctx)
	if err != nil {
		cr.log.Error(err)
		return nil
	}
	defer c.Release()

	// TODO: current code can't change only changed things - working with whole object via handler>usecase>_here_
	r, err := c.Exec(ctx, "update controller.clusters set name=$2, capacity=$3 where id=$1;", cl.ID, cl.Name, cl.Capacity)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch {
			case pgErr.Code == "23505":
				cr.log.Warnw("issue with cluster on update",
					"error_body", pgErr.Message,
					"error_code", pgErr.Code,
				)
				return cluster.ErrAlreadyExist
			default:
				cr.log.Errorw("issue with cluster on update",
					"error_body", pgErr.Message,
					"error_code", pgErr.Code,
				)
				return cluster.ErrNotFound
			}
		}
		cr.log.Error(err)
		return cluster.ErrNotFound
	}
	if r.RowsAffected() != 1 {
		cr.log.Debug("update for non-existing cluster ID")
		return cluster.ErrNotFound
	}

	// FIXME: error case `name already exist` (name is unique in db)
	// FIXME: check if existing ID used or 400
	// TODO: test case `change name`
	// TODO: test case `change ZoneID`
	// TODO: test case `change capacity`

	return nil
}

func (cr *ClusterRepo) Delete(ctx context.Context, cl *model.Cluster) error {
	c, err := cr.pool.Acquire(ctx)
	if err != nil {
		cr.log.Error(err)
		return err
	}
	defer c.Release()

	r, err := c.Exec(ctx, "delete from controller.clusters where id=$1", cl.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch {
			case pgErr.Code == "23503":
				cr.log.Warnw("issue with cluster on delete",
					"error_body", pgErr.Message,
					"error_code", pgErr.Code,
				)
				return cluster.ErrHaveServices
			default:
				cr.log.Errorw("issue with cluster on delete",
					"error_body", pgErr.Message,
					"error_code", pgErr.Code,
				)
				return err
			}
		}
		cr.log.Error(err)
		return err
	}
	if r.RowsAffected() != 1 {
		cr.log.Debug("delete for non-existing cluster ID")
		return nil
	}
	return nil
}
