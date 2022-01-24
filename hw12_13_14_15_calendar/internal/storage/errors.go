package storage

import "errors"

var (
	ErrOverlapID           = errors.New("overlap id of event")
	ErrTitle               = errors.New("title should be filled")
	ErrWithPlannedTime     = errors.New("planned time in the past")
	ErrMismatchedID        = errors.New("id should be equal")
	ErrEventIDNotFound     = errors.New("event id not found")
	ErrInvalidEventID      = errors.New("id not match uuid type")
	ErrTimeIsBusy          = errors.New("time is busy")
	ErrEmptyUserIDField    = errors.New("empty user id field")
	ErrUnkownUserID        = errors.New("unknown user")
	ErrInvalidUserID       = errors.New("id not match uuid type")
	ErrUnkownTypeOfStorage = errors.New("unknown type of storage")
)
