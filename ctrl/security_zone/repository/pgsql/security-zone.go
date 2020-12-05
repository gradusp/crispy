package pgsql

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/gradusp/crispy/ctrl/model"
	"github.com/gradusp/crispy/ctrl/security_zone"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"log"
)

func NewSecurityZoneRepo(db *pg.DB, kv *api.KV) *SecurityZoneRepo {
	return &SecurityZoneRepo{
		db: db,
		kv: kv,
	}
}

type SecurityZoneRepo struct {
	db  orm.DB
	kv  *api.KV
	log *zap.Logger
}

func (szr *SecurityZoneRepo) Create(ctx context.Context, sz *model.SecurityZone) (*model.SecurityZone, error) {
	// TODO: trace logs here
	r, err := szr.db.Model(sz).Where("name = ?", sz.Name).Exists()
	if err != nil {
		panic(err)
	} else if r {
		if err = szr.db.Model(sz).Where("name = ?", sz.Name).Select(); err != nil {
			return nil, err
		}
		return sz, security_zone.ErrSecurityZoneAlreadyExist
	}

	_, err = szr.db.Model(sz).Insert()
	if err != nil {
		return nil, err
	}

	if err = szr.db.Model(sz).Where("name = ?", sz.Name).Select(); err != nil {
		return nil, err
	}

	// PUT a new KV pair
	p := &api.KVPair{Key: "lbos/", Value: []byte("1000")}
	_, err = szr.kv.Put(p, nil)
	if err != nil {
		panic(err)
	}

	return sz, nil
}

func (szr *SecurityZoneRepo) Get(ctx context.Context) ([]*model.SecurityZone, error) {
	var r []*model.SecurityZone
	err := szr.db.Model(&r).Select()
	if err != nil {
		log.Print(err)
	}

	return r, err
}

func (szr *SecurityZoneRepo) GetByID(ctx context.Context, sz *model.SecurityZone) (*model.SecurityZone, error) {
	err := szr.db.Model(sz).Where("id = ?", sz.ID).Select()

	return sz, err
}

func (szr *SecurityZoneRepo) Update(ctx context.Context, sz *model.SecurityZone) error {
	_, err := szr.db.Model(sz).Where("id = ?", sz.ID).Update()

	return err
}

func (szr *SecurityZoneRepo) Delete(ctx context.Context, sz *model.SecurityZone) error {
	_, err := szr.db.Model(sz).WherePK().Delete()

	return err
}
