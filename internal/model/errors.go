package model

import "errors"

var (
	ErrBadRequest       = errors.New("bad request")
	ErrBadRequestUserID = errors.New("invalid user id")
	ErrNotFound         = errors.New("not found")
	ErrInternalServer   = errors.New("internal server")
)
