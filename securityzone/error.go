package securityzone

import "errors"

var (
	ErrSecurityZoneNotFound     = errors.New("security zone not found")
	ErrSecurityzoneAlreadyExist = errors.New("security zone already exist")
)
