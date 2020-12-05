package pgsql

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/gradusp/crispy/ctrl/cluster"
	"github.com/gradusp/crispy/ctrl/model"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

func NewClusterRepo(db *pg.DB, kv *api.KV) *ClusterRepo {
	return &ClusterRepo{
		db: db,
		kv: kv,
	}
}

type ClusterRepo struct {
	db  orm.DB
	kv  *api.KV
	log *zap.Logger
}

func (cr *ClusterRepo) Create(ctx context.Context, sz *model.SecurityZone, c *model.Cluster) (*model.Cluster, error) {
	// TODO: missing trace logs here
	// TODO: error case `invalid security_zone_id`

	// Request SecurityZone from DB
	err := cr.db.Model(sz).WherePK().Select()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	c.SecurityZoneID = sz.ID

	fmt.Printf("%+v\n",sz)
	fmt.Printf("%+v\n",c)

	r, err := cr.db.Model(c).Where("name = ?", c.Name).Exists()
	if err != nil {
		panic(err)
	} else if r {
		if err = cr.db.Model(c).Where("name = ?", c.Name).Select(); err != nil {
			fmt.Println(err)
			return nil, err
		}
		return c, cluster.ErrClusterAlreadyExist
	}

	_, err = cr.db.Model(c).Insert()
	if err != nil {
		fmt.Println("REPO INSERT:",err)
		return nil, err
	}

	if err = cr.db.Model(c).Column("cluster.*").Where("cluster.name = ?", c.Name).Relation("SecurityZone").Select(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Printf("%+v\n",c)

	return c, nil
}

func (cr *ClusterRepo) Get(ctx context.Context) ([]*model.Cluster, error) {
	panic("repo is not implemented yet")
}
