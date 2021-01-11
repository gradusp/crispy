package consul

import (
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

type SecurityzoneConsulRepo struct {
	storage *api.KV
	logger  *zap.SugaredLogger
}

func NewSecurityzoneConsulRepo(kv *api.KV, l *zap.SugaredLogger) *SecurityzoneConsulRepo {
	return &SecurityzoneConsulRepo{
		storage: kv,
		logger:  l,
	}
}

func (szcr *SecurityzoneConsulRepo) Create() error {
	return nil
}
