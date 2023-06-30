package tx

import "context"

type Tx interface {
	Transaction(fn func(tx any) error) error
}

type CtxTx interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
