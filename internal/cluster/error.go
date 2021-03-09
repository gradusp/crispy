package cluster

import "errors"

var (
	ErrAlreadyExist = errors.New("cluster with such name already exist")
	ErrHaveServices = errors.New("cluster contain related Services so it can't be deleted")
	ErrHaveNodes    = errors.New("cluster contain related Nodes so it can't be deleted")
	ErrNotFound     = errors.New("cluster not found")
)
