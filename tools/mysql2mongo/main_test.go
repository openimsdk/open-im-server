package main

import (
	"errors"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/mgo"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
	"testing"
)

func getColl1(obj any) (_ *mongo.Collection, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New("not found")
		}
	}()
	stu := reflect.ValueOf(obj).Elem()
	typ := reflect.TypeOf(&mongo.Collection{}).String()
	for i := 0; i < stu.NumField(); i++ {
		field := stu.Field(i)
		if field.Type().String() == typ {
			return (*mongo.Collection)(field.UnsafePointer()), nil
		}
	}
	return nil, errors.New("not found")
}

func TestName(t *testing.T) {
	coll, err := getColl1(&mgo.GroupMgo{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(coll)

}
