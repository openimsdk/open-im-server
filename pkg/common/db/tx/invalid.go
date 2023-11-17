package tx

import (
	"context"
	"github.com/OpenIMSDK/tools/tx"
)

func NewInvalidTx() tx.CtxTx {
	return invalid{}
}

type invalid struct{}

func (m invalid) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}
