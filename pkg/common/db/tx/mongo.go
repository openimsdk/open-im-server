package tx

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewMongo(client *mongo.Client) CtxTx {
	return &_Mongo{
		client: client,
	}
}

type _Mongo struct {
	client *mongo.Client
}

func (m *_Mongo) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	sess, err := m.client.StartSession()
	if err != nil {
		return err
	}
	sCtx := mongo.NewSessionContext(ctx, sess)
	defer sess.EndSession(sCtx)
	if err := fn(sCtx); err != nil {
		_ = sess.AbortTransaction(sCtx)
		return err
	}
	return utils.Wrap(sess.CommitTransaction(sCtx), "")
}
