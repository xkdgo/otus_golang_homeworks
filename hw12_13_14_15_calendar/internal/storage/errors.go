package storage

import "errors"

var (
	ErrOverlapID       = errors.New("overlap id of event")
	ErrTitle           = errors.New("title should be filled")
	ErrWithPlannedTime = errors.New("planned time in the past")
)
