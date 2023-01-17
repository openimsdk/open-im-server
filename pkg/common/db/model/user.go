package model

import (
	"Open_IM/pkg/common/db/mysql"
	"context"
)

type UserModel struct {
	db *mysql.User
}

func NewGroupUser(ctx context.Context) *UserModel {
	var userModel UserModel
	userModel.db = mysql.NewUserDB()
	return &userModel
}

func (u *UserModel) Find(ctx context.Context, userIDs []string) (users []*mysql.User, err error) {
	return u.db.Find(ctx, userIDs)
}

func (u *UserModel) Create(ctx context.Context, users []*mysql.User) error {
	return u.db.Create(ctx, users)
}
