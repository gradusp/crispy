package pgsql

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"

	"github.com/gradusp/crispy/ctrl/model"
	"github.com/gradusp/crispy/ctrl/securityzone"
)

func NewSecurityzoneRepo(db *pg.DB, kv *api.KV, l *zap.SugaredLogger) *SecurityzonePostgresRepo {
	return &SecurityzonePostgresRepo{
		db: db,
		//kv:  kv,
		log: l,
	}
}

type SecurityzonePostgresRepo struct {
	db orm.DB
	//kv  *api.KV
	log *zap.SugaredLogger
}

func (szr *SecurityzonePostgresRepo) Create(ctx context.Context, sz *model.SecurityZone) (*model.SecurityZone, error) {
	// TODO: trace logs here
	r, err := szr.db.Model(sz).Where("name = ?", sz.Name).Exists()
	if err != nil {
		panic(err)
	} else if r {
		// FIXME: sloppyReassign
		if err = szr.db.Model(sz).Where("name = ?", sz.Name).Select(); err != nil {
			return nil, err
		}
		return sz, securityzone.ErrSecurityzoneAlreadyExist
	}

	_, err = szr.db.Model(sz).Insert()
	if err != nil {
		return nil, err
	}

	// FIXME: sloppyReassign
	if err = szr.db.Model(sz).Where("name = ?", sz.Name).Select(); err != nil {
		return nil, err
	}

	// PUT a new KV pair
	//p := &api.KVPair{Key: "lbos/", Value: []byte("1000")}
	//_, err = szr.kv.Put(p, nil)
	//if err != nil {
	//	panic(err)
	//}
	return sz, nil
}

func (szr *SecurityzonePostgresRepo) Get(ctx context.Context) ([]*model.SecurityZone, error) {
	var r []*model.SecurityZone
	err := szr.db.ModelContext(ctx, &r).Select()
	if err != nil {
		szr.log.Error("securityZones can't be selected", err)
	}
	return r, err
}

func (szr *SecurityzonePostgresRepo) GetByID(ctx context.Context, sz *model.SecurityZone) (*model.SecurityZone, error) {
	// FIXME: error case `no UUID found` (500>400)
	err := szr.db.ModelContext(ctx, sz).WherePK().Select()
	if err != nil {
		szr.log.Error("securityZones can't be selected by ID", err)
	}
	return sz, err
}

func (szr *SecurityzonePostgresRepo) Update(ctx context.Context, sz *model.SecurityZone) error {
	// FIXME: error case `name already exist` (name is unique in db)
	_, err := szr.db.ModelContext(ctx, sz).Where("id = ?", sz.ID).Update()
	return err
}

func (szr *SecurityzonePostgresRepo) Delete(ctx context.Context, sz *model.SecurityZone) error {
	_, err := szr.db.ModelContext(ctx, sz).WherePK().Delete()
	return err
}
