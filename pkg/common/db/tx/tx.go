package tx

import (
	"context"
	"github.com/OpenIMSDK/tools/tx"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewMongoTx(client *mongo.Client) tx.CtxTx {
	return &mongoTx{
		client: client,
	}
}

type mongoTx struct {
	client *mongo.Client
}

func (m *mongoTx) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	sess, err := m.client.StartSession()
	if err != nil {
		return err
	}
	_, err = sess.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		return nil, fn(ctx)
	})
	return err
}
