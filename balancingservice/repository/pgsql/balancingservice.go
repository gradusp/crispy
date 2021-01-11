package pgsql

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"

	"github.com/gradusp/crispy/model"
)

type BalancingserviceRepo struct {
	db  orm.DB
	kv  *api.KV
	log *zap.SugaredLogger
}

func NewBalancingserviceRepo(db *pg.DB, kv *api.KV, l *zap.SugaredLogger) *BalancingserviceRepo {
	return &BalancingserviceRepo{
		db:  db,
		kv:  kv,
		log: l,
	}
}

func (bsr *BalancingserviceRepo) Create(ctx context.Context) (*model.BalancingService, error) { // TODO: implement
	panic("implement me")
}

func (bsr *BalancingserviceRepo) Get(ctx context.Context) ([]*model.BalancingService, error) { // TODO: implement
	panic("implement me")
}

func (bsr *BalancingserviceRepo) GetByID(ctx context.Context) (*model.BalancingService, error) { // TODO: implement
	panic("implement me")
}

func (bsr *BalancingserviceRepo) Update(ctx context.Context) error { // TODO: implement
	panic("implement me")
}

func (bsr *BalancingserviceRepo) Delete(ctx context.Context) error { // TODO: implement
	panic("implement me")
}
