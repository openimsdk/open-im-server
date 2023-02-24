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

type Api interface {
	OperationID() string
	OpUserID() string
	Context() context.Context
	Bind(req any) error
	Resp(resp any, err error)
}
