package node

import "errors"

var (
	ErrAlreadyExist = errors.New("node already exist")
	ErrNotFound     = errors.New("node not found")
	ErrWrongQuery   = errors.New("query for node is incorrect")
)
