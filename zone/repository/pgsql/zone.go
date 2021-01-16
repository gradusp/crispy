package pgsql

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"

	"github.com/gradusp/crispy/model"
	"github.com/gradusp/crispy/zone"
)

type ZonePostgresRepo struct {
	//kv  *api.KV
	log  *zap.SugaredLogger
	pool *pgxpool.Pool
}

func NewZonePostgresRepo(pool *pgxpool.Pool, kv *api.KV, l *zap.SugaredLogger) *ZonePostgresRepo {
	return &ZonePostgresRepo{
		//kv:  kv,
		pool: pool,
		log:  l,
	}
}

func (zr *ZonePostgresRepo) Create(ctx context.Context, sz *model.Zone) (*model.Zone, error) {
	c, err := zr.pool.Acquire(ctx)
	if err != nil {
		zr.log.Error(err)
		return nil, err
	}
	defer c.Release()

	if rowExists(ctx, c, "select name from controller.zones where name=$1", sz.Name) {
		err := c.QueryRow(ctx, "select id from controller.zones where name=$1", sz.Name).Scan(&sz.ID)
		if err != nil {
			zr.log.Error(err)
			return nil, err
		}
		return sz, zone.ErrZoneAlreadyExist
	}

	if err := c.QueryRow(ctx, "insert into controller.zones (name) values ($1) returning id", sz.Name).Scan(&sz.ID); err != nil {
		zr.log.Error(err)
		return nil, err
	}

	// PUT a new KV pair
	//p := &api.KVPair{Key: "lbos/", Value: []byte("1000")}
	//_, err = zr.kv.Put(p, nil)
	//if err != nil {
	//	panic(err)
	//}
	return sz, nil
}

func (zr *ZonePostgresRepo) Get(ctx context.Context) ([]*model.Zone, error) {
	c, err := zr.pool.Acquire(ctx)
	if err != nil {
		zr.log.Error(err)
		return nil, err
	}
	defer c.Release()

	zones, err := c.Query(ctx, "select id, name from controller.zones")
	if err != nil {
		zr.log.Error("zones can't be selected", err) // FIXME: rework into advanced err logic (GetByID)
		return nil, err
	}

	var r []*model.Zone
	for zones.Next() {
		var zone model.Zone
		err = zones.Scan(&zone.ID, &zone.Name)
		if err != nil {
			return nil, err
		}
		r = append(r, &zone)
	}
	err = zones.Err()

	return r, err
}

func (zr *ZonePostgresRepo) GetByID(ctx context.Context, sz *model.Zone) (*model.Zone, error) {
	// TODO: investigate if following part could be externalized
	c, err := zr.pool.Acquire(ctx)
	if err != nil {
		zr.log.Error(err)
		return nil, err
	}
	defer c.Release()

	if err = c.QueryRow(ctx, "select name from controller.zones where id=$1", sz.ID).Scan(&sz.Name); err != nil {
		// since it is not exact error, just specific condition of data,
		// we detect it via default const of pgx lib
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, zone.ErrZoneNotFound
		}

		// own (and big) pool of native pgsql errors parsed here
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			zr.log.Errorw("can't get zone",
				"error_body", pgErr.Message,
				"error_code", pgErr.Code,
			)
			return nil, zone.ErrZoneNotFound
		}
		return nil, err
	}
	return sz, err
}

func (zr *ZonePostgresRepo) Update(ctx context.Context, sz *model.Zone) error {
	c, err := zr.pool.Acquire(ctx)
	if err != nil {
		zr.log.Error(err)
		return err
	}
	defer c.Release()

	r, err := c.Exec(ctx, "update controller.zones set name=$2 where id=$1", sz.ID, sz.Name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch {
			case pgErr.Code == "23505":
				zr.log.Warnw("issue with zone on update",
					"error_body", pgErr.Message,
					"error_code", pgErr.Code,
				)
				return zone.ErrZoneAlreadyExist
			default:
				zr.log.Errorw("issue with zone on update",
					"error_body", pgErr.Message,
					"error_code", pgErr.Code,
				)
				return zone.ErrZoneNotFound
			}
		}
		zr.log.Error(err)
		return zone.ErrZoneNotFound
	}
	if r.RowsAffected() != 1 {
		zr.log.Debug("update for non-existing zone ID")
		return zone.ErrZoneNotFound
	}
	return nil
}

func (zr *ZonePostgresRepo) Delete(ctx context.Context, sz *model.Zone) error {
	c, err := zr.pool.Acquire(ctx)
	if err != nil {
		zr.log.Error(err)
		return err
	}
	defer c.Release()

	r, err := c.Exec(ctx, "delete from controller.zones where id=$1", sz.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch {
			case pgErr.Code == "23503":
				zr.log.Warnw("issue with zone on delete",
					"error_body", pgErr.Message,
					"error_code", pgErr.Code,
				)
				return zone.ErrZoneHaveClusters
			default:
				zr.log.Errorw("issue with zone on delete",
					"error_body", pgErr.Message,
					"error_code", pgErr.Code,
				)
				return err
			}
		}
		zr.log.Error(err)
		return err
	}
	if r.RowsAffected() != 1 {
		zr.log.Debug("delete for non-existing zone ID")
		return nil
	}
	return nil
}

func rowExists(ctx context.Context, c *pgxpool.Conn, q string, args ...interface{}) bool {
	var exists bool

	query := fmt.Sprintf("select exists (%s)", q)
	err := c.QueryRow(ctx, query, args...).Scan(&exists)
	if err != nil && err != pgx.ErrNoRows {
		panic(err)
	}
	return exists
}
