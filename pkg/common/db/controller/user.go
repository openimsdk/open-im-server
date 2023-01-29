package controller

import (
	"Open_IM/pkg/common/db/relation"
	"context"
	"gorm.io/gorm"
)

type UserInterface interface {
	Find(ctx context.Context, userIDs []string) (users []*relation.User, err error)
	Create(ctx context.Context, users []*relation.User) error
	Take(ctx context.Context, userID string) (user *relation.User, err error)
	Update(ctx context.Context, users []*relation.User) (err error)
	UpdateByMap(ctx context.Context, userID string, args map[string]interface{}) (err error)
	GetByName(ctx context.Context, userName string, showNumber, pageNumber int32) (users []*relation.User, err error)
}

type UserController struct {
	database UserDatabaseInterface
}

func (u *UserController) Find(ctx context.Context, userIDs []string) (users []*relation.User, err error) {
	return u.database.Find(ctx, userIDs)
}
func (u *UserController) Create(ctx context.Context, users []*relation.User) error {
	return u.database.Create(ctx, users)
}
func (u *UserController) Take(ctx context.Context, userID string) (user *relation.User, err error) {
	return u.database.Take(ctx, userID)
}
func (u *UserController) Update(ctx context.Context, users []*relation.User) (err error) {
	return u.database.Update(ctx, users)
}
func (u *UserController) UpdateByMap(ctx context.Context, userID string, args map[string]interface{}) (err error) {
	return u.database.UpdateByMap(ctx, userID, args)
}
func (u *UserController) GetByName(ctx context.Context, userName string, showNumber, pageNumber int32) (users []*relation.User, err error) {
	return u.database.GetByName(ctx, userName, showNumber, pageNumber)
}

func NewUserController(db *gorm.DB) UserInterface {
	controller := &UserController{database: newUserDatabase(db)}
	return controller
}

type UserDatabaseInterface interface {
	Find(ctx context.Context, userIDs []string) (users []*relation.User, err error)
	Create(ctx context.Context, users []*relation.User) error
	Take(ctx context.Context, userID string) (user *relation.User, err error)
	Update(ctx context.Context, users []*relation.User) (err error)
	UpdateByMap(ctx context.Context, userID string, args map[string]interface{}) (err error)
	GetByName(ctx context.Context, userName string, showNumber, pageNumber int32) (users []*relation.User, err error)
}

type UserDatabase struct {
	sqlDB *relation.User
}

func newUserDatabase(db *gorm.DB) UserDatabaseInterface {
	sqlDB := relation.NewUserDB(db)
	database := &UserDatabase{
		sqlDB: sqlDB,
	}
	return database
}

func (u *UserDatabase) Find(ctx context.Context, userIDs []string) (users []*relation.User, err error) {
	return u.sqlDB.Find(ctx, userIDs)
}

func (u *UserDatabase) Create(ctx context.Context, users []*relation.User) error {
	return u.sqlDB.Create(ctx, users)
}
func (u *UserDatabase) Take(ctx context.Context, userID string) (user *relation.User, err error) {
	return u.sqlDB.Take(ctx, userID)
}
func (u *UserDatabase) Update(ctx context.Context, users []*relation.User) (err error) {
	return u.sqlDB.Update(ctx, users)
}
func (u *UserDatabase) UpdateByMap(ctx context.Context, userID string, args map[string]interface{}) (err error) {
	return u.sqlDB.UpdateByMap(ctx, userID, args)
}
func (u *UserDatabase) GetByName(ctx context.Context, userName string, showNumber, pageNumber int32) (users []*relation.User, err error) {
	return u.sqlDB.GetByName(ctx, userName, showNumber, pageNumber)
}
