package pgsql

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/gradusp/crispy/ctrl/model"
	"github.com/gradusp/crispy/ctrl/security_zone"
	"go.uber.org/zap"
)

func NewSecurityZoneRepo(db *pg.DB) *SecurityZoneRepo {
	return &SecurityZoneRepo{
		db: db,
	}
}

type SecurityZoneRepo struct {
	db  orm.DB
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
	return sz, nil
}

func (szr *SecurityZoneRepo) Get(ctx context.Context) ([]*model.SecurityZone, error) {
	var r []*model.SecurityZone
	err := szr.db.Model(&r).Select()

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
