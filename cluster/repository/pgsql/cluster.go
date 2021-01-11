package pgsql

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"

	"github.com/gradusp/crispy/cluster"
	"github.com/gradusp/crispy/model"
)

func NewClusterRepo(db *pg.DB, kv *api.KV, l *zap.SugaredLogger) *ClusterRepo {
	return &ClusterRepo{
		db:  db,
		kv:  kv,
		log: l,
	}
}

type ClusterRepo struct {
	db  orm.DB
	kv  *api.KV
	log *zap.SugaredLogger
}

func (cr *ClusterRepo) Create(ctx context.Context, sz *model.SecurityZone, c *model.Cluster) (*model.Cluster, error) {
	// TODO: missing trace logs here?
	// TODO: error case `invalid security_zone_id`

	// Request SecurityZone from DB
	err := cr.db.Model(sz).WherePK().Select()
	if err != nil {
		cr.log.Error("securityZone for cluster not found", err)
		return nil, err
	}
	c.SecurityZoneID = sz.ID

	r, err := cr.db.Model(c).Where("name = ?", c.Name).Exists()
	if err != nil {
		cr.log.Error(err)
	} else if r {
		if err = cr.db.Model(c).Where("name = ?", c.Name).Select(); err != nil {
			cr.log.Error(err)
			return nil, err
		}
		return c, cluster.ErrClusterAlreadyExist
	}

	_, err = cr.db.Model(c).Insert()
	if err != nil {
		cr.log.Error(err)
		return nil, err
	}

	if err = cr.db.Model(c).Column("cluster.*").Where("cluster.name = ?", c.Name).Select(); err != nil {
		cr.log.Error(err)
		return nil, err
	}

	return c, nil
}

func (cr *ClusterRepo) Get(ctx context.Context) ([]*model.Cluster, error) {
	var r []*model.Cluster
	err := cr.db.Model(&r).Select()
	if err != nil {
		cr.log.Error(err)
		return nil, err
	}
	return r, nil
}

func (cr *ClusterRepo) GetByID(ctx context.Context, c *model.Cluster) (*model.Cluster, error) {
	err := cr.db.Model(c).WherePK().Select()
	return c, err
}

func (cr *ClusterRepo) Update(ctx context.Context, sz *model.SecurityZone, c *model.Cluster) error {
	// FIXME: error case `name already exist` (name is unique in db)
	// TODO: test case `change name`
	// TODO: test case `change SecurityZoneID`
	// TODO: test case `change capacity`

	// Select Security Zone from DB by ID so we can be sure such Security Zone is real
	err := cr.db.Model(sz).WherePK().Select()
	if err != nil {
		// TODO: discover better way to check what exact err comes from pg (next block maybe?)
		// FIXME: rework guessing that err is only when nothing returned
		return cluster.ErrRequestedSecZoneNotFound
	}
	c.SecurityZoneID = sz.ID

	_, err = cr.db.Model(c).WherePK().Update()
	//_, err = cr.db.Model(c).Where("id = ?", c.ID).Update()
	if err != nil {
		pgErr, ok := err.(pg.Error)
		if ok && pgErr.IntegrityViolation() {
			return cluster.ErrClusterAlreadyExist
		}
		panic(err)
		return err
	}

	return err
}

func (cr *ClusterRepo) Delete(ctx context.Context, c *model.Cluster) error {
	_, err := cr.db.Model(c).WherePK().Delete()
	return err
}
