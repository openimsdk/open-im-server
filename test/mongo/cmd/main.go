package main

import (
	mongo2 "Open_IM/test/mongo"
	"context"
	"flag"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:37017/openIM/?maxPoolSize=100")
	var err error
	mongo2.Client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}
	err = mongo2.Client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to MongoDB!")
}

func main() {
	userID := flag.String("userID", "", "userID")
	flag.Parse()
	fmt.Println("userID:", *userID)
	mongo2.GetUserAllChat(*userID)
}
