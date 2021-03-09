package service

import "errors"

var (
	ErrAlreadyExist = errors.New("service with such name already exist")
	ErrNotFound     = errors.New("service not found")
)
