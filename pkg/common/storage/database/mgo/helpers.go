package mgo

import (
	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/mongo"
)

func IsNotFound(err error) bool {
	return errs.Unwrap(err) == mongo.ErrNoDocuments
}
