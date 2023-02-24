package api2rpc

import (
	"context"
)

type Ignore struct{}

type ApiBind[A, B any] interface {
	OperationID() string
	OpUserID() (string, error)
	Bind(*A) error
	Context() context.Context
	Resp(resp *B, err error)
}
