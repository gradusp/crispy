package real

import "errors"

var (
	ErrAlreadyExist = errors.New("real already exist")
	ErrNotFound     = errors.New("real not found")
	ErrWrongQuery   = errors.New("query for real is incorrect")
)
