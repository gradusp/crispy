package cluster

import "errors"

var (
	ErrClusterNotFound = errors.New("cluster not found")
	ErrClusterAlreadyExist = errors.New("cluster already exist")
)
