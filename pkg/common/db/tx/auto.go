package tx

import (
	"context"
	"github.com/OpenIMSDK/tools/tx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewAuto(ctx context.Context, cli *mongo.Client) (tx.CtxTx, error) {
	var res map[string]any
	if err := cli.Database("admin").RunCommand(ctx, bson.M{"isMaster": 1}).Decode(&res); err != nil {
		return nil, err
	}
	if _, ok := res["setName"]; ok {
		return NewMongoTx(cli), nil
	}
	return NewInvalidTx(), nil
}
