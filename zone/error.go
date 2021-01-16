package zone

import "errors"

var (
	ErrZoneAlreadyExist = errors.New("zone already exist")
	ErrZoneNotFound     = errors.New("zone not found")
	ErrZoneHaveClusters = errors.New("zone contains related Clusters so it can't be deleted")
)
