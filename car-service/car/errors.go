package car

import "errors"

var (
	ErrCarNotFound     = errors.New("order not found")
	ErrCmdRepository   = errors.New("unable to command repository")
	ErrQueryRepository = errors.New("unable to query repository")
)
