package cluster

import "errors"

var (
	ErrClusterNotFound          = errors.New("cluster not found")
	ErrClusterAlreadyExist      = errors.New("cluster with such name already exist")
	ErrRequestedSecZoneNotFound = errors.New("there is no Security Zone with provided UUID")
)
