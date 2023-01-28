package controller

import (
	"Open_IM/pkg/common/db/mysql"
	"context"
)

type UserModel struct {
	db *relation.User
}

func NewGroupUser(ctx context.Context) *UserModel {
	var userModel UserModel
	userModel.db = relation.NewUserDB()
	return &userModel
}

func (u *UserModel) Find(ctx context.Context, userIDs []string) (users []*relation.User, err error) {
	return u.db.Find(ctx, userIDs)
}

func (u *UserModel) Create(ctx context.Context, users []*relation.User) error {
	return u.db.Create(ctx, users)
}
