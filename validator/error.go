package validator

import "errors"

var (
	ErrorKnownBlock    = errors.New("block already known")
	ErrorUnknownParent = errors.New("unknown parent")
	IncorrectTxHash    = errors.New("incorrect transaction hash")
)
