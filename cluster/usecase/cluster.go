package usecase

import (
	"github.com/gradusp/crispy/cluster"
)

type ClusterUsecase struct {
	clusterRepo cluster.Repository
}

func NewClusterUsecase(clusterRepo cluster.Repository) *ClusterUsecase {
	return &ClusterUsecase{
		clusterRepo: clusterRepo,
	}
}
