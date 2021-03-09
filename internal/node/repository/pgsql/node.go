package pgsql

import (
	"context"
	"errors"
	"fmt"

	"github.com/gradusp/crispy/internal/node"

	"github.com/jackc/pgx/v4"

	"github.com/gradusp/crispy/internal/model"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type Repo struct {
	log  *zap.SugaredLogger
	pool *pgxpool.Pool
}

func NewPgRepo(pool *pgxpool.Pool, l *zap.SugaredLogger) *Repo {
	return &Repo{
		log:  l,
		pool: pool,
	}
}

func (r Repo) Create(ctx context.Context, n *model.Node) (*model.Node, error) {
	c, err := r.pool.Acquire(ctx)
	if err != nil {
		r.log.Error(err)
		return nil, err
	}
	defer c.Release()

	if rowExists(ctx, c, "select id from controller.nodes where addr=$1", n.Addr.String()) {
		err := c.QueryRow(ctx, "select id from controller.nodes where addr=$1;", n.Addr.String()).Scan(&n.ID)
		if err != nil {
			r.log.Error(err)
			return nil, err
		}
		return n, node.ErrAlreadyExist
	}

	query := `insert into controller.nodes (cluster_id, addr, hostname) values ($1, $2, $3) returning id;`
	if err := c.QueryRow(ctx, query, n.ClusterID, n.Addr.String(), n.Hostname).Scan(&n.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch {
			// FIXME: switch is not needed with only default but can be useful for quick handle of pgsql errors
			default:
				r.log.Errorw("issue with node on create",
					"error_body", pgErr.Message,
					"error_code", pgErr.Code,
				)
				return nil, err
			}
		}
		r.log.Error(err)
		return nil, err
	}

	return n, nil
}

//func (r Repo) GetAll(ctx context.Context) ([]*model.Node, error) {
//	c, err := r.pool.Acquire(ctx)
//	if err != nil {
//		r.log.Error(err)
//		return nil, err
//	}
//	defer c.Release()
//
//	var reals []*model.Real
//	query := `select id, service_id, addr, port from controller.reals;`
//	if rows, err := c.Query(ctx, query); err != nil {
//		r.log.Error("reals can't be selected", err) // TODO: better error handling
//		return nil, err
//	} else {
//		for rows.Next() {
//			var r model.Real
//			if err := rows.Scan(&r.ID, &r.ServiceID, &r.Addr, &r.Port); err != nil {
//				return nil, err
//			} else {
//				reals = append(reals, &r)
//			}
//		}
//	}
//	return reals, nil
//}

func (r Repo) GetByID(ctx context.Context, n *model.Node) (*model.Node, error) {
	c, err := r.pool.Acquire(ctx)
	if err != nil {
		r.log.Error(err)
		return nil, err
	}
	defer c.Release()

	query := `select id, cluster_id, addr, hostname from controller.nodes where id=$1;`
	if err = c.QueryRow(ctx, query, n.ID).Scan(&n.ID, &n.ClusterID, &n.Addr, &n.Hostname); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, node.ErrNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			r.log.Errorw("can't get real",
				"error_body", pgErr.Message,
				"error_code", pgErr.Code,
			)
			return nil, node.ErrNotFound
		}
		return nil, err
	}
	return n, nil
}

func (r Repo) GetNodeByField(ctx context.Context, where string) ([]*model.Node, error) {
	var nodes []*model.Node
	c, err := r.pool.Acquire(ctx)
	if err != nil {
		r.log.Error(err)
		return nil, err
	}
	defer c.Release()

	q := fmt.Sprintf("select id, cluster_id, addr, hostname from controller.nodes %s;", where)
	fmt.Printf("\n%s\n\n", q)
	rows, err := c.Query(ctx, q)
	if err != nil {
		r.log.Error("nodes can't be selected ", err) // TODO: better error handling
		return nil, err
	}

	for rows.Next() {
		var n model.Node
		err = rows.Scan(&n.ID, &n.ClusterID, &n.Addr, &n.Hostname)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, &n)
	}

	return nodes, nil
}

// TODO: refactor GetByAddr & GetByServiceID
//func (r Repo) GetByAddr(ctx context.Context, a string) ([]*model.Node, error) {
//	c, err := r.pool.Acquire(ctx)
//	if err != nil {
//		r.log.Error(err)
//		return nil, err
//	}
//	defer c.Release()
//
//	query := `select id, service_id, addr, port from controller.reals where addr=$1;`
//	rows, err := c.Query(ctx, query, a)
//	if err != nil {
//		r.log.Error("reals can't be selected", err) // TODO: better error handling
//		return nil, err
//	}
//
//	var reals []*model.Real
//	for rows.Next() {
//		var r model.Real
//		err = rows.Scan(&r.ID, &r.ServiceID, &r.Addr, &r.Port)
//		if err != nil {
//			return nil, err
//		}
//		reals = append(reals, &r)
//	}
//
//	return reals, nil
//}

// TODO: refactor GetByAddr & GetByServiceID
//func (r Repo) GetByClusterID(ctx context.Context, cid string) ([]*model.Node, error) {
//	c, err := r.pool.Acquire(ctx)
//	if err != nil {
//		r.log.Error(err)
//		return nil, err
//	}
//	defer c.Release()
//
//	query := `select id, service_id, addr, port from controller.reals where service_id=$1;`
//	rows, err := c.Query(ctx, query, sid)
//	if err != nil {
//		r.log.Error("reals can't be selected", err) // TODO: better error handling
//		return nil, err
//	}
//
//	var reals []*model.Real
//	for rows.Next() {
//		var r model.Real
//		err = rows.Scan(&r.ID, &r.ServiceID, &r.Addr, &r.Port)
//		if err != nil {
//			return nil, err
//		}
//		reals = append(reals, &r)
//	}
//
//	return reals, nil
//}

func (r Repo) Delete(ctx context.Context, n *model.Node) error {
	c, err := r.pool.Acquire(ctx)
	if err != nil {
		r.log.Error(err)
		return err
	}
	defer c.Release()

	// TODO: advanced error checks

	res, err := c.Exec(ctx, "delete from controller.nodes where id=$1;", n.ID)
	if err != nil {
		r.log.Error(err)
		return err
	}
	if res.RowsAffected() != 1 {
		r.log.Debug("delete for non-existing Node ID")
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
