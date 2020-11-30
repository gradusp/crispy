package security_zone

import "errors"

var (
	ErrSecurityZoneNotFound = errors.New("security zone not found")
	ErrSecurityZoneAlreadyExist = errors.New("security zone already exist")
)
