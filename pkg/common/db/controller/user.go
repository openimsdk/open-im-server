package controller

import (
	"Open_IM/pkg/common/db/relation"
	"Open_IM/pkg/common/db/table"
	"context"
	"gorm.io/gorm"
)

type UserInterface interface {
	//获取指定用户的信息 如果有记录未找到 也返回错误
	Find(ctx context.Context, userIDs []string) (users []*table.UserModel, err error)
	Create(ctx context.Context, users []*table.UserModel) error
	Update(ctx context.Context, users []*table.UserModel) (err error)
	UpdateByMap(ctx context.Context, userID string, args map[string]interface{}) (err error)
	GetByName(ctx context.Context, userName string, showNumber, pageNumber int32) (users []*table.UserModel, count int64, err error)
	GetByNameAndID(ctx context.Context, content string, showNumber, pageNumber int32) (users []*table.UserModel, count int64, err error)
	Get(ctx context.Context, showNumber, pageNumber int32) (users []*table.UserModel, count int64, err error)
	//userIDs是否存在 只要有一个存在就为true
	IsExist(ctx context.Context, userIDs []string) (exist bool, err error)
}

type UserController struct {
	database UserDatabaseInterface
}

func (u *UserController) Find(ctx context.Context, userIDs []string) (users []*table.UserModel, err error) {
	return u.database.Find(ctx, userIDs)
}
func (u *UserController) Create(ctx context.Context, users []*table.UserModel) error {
	return u.database.Create(ctx, users)
}
func (u *UserController) Take(ctx context.Context, userID string) (user *table.UserModel, err error) {
	return u.database.Take(ctx, userID)
}
func (u *UserController) Update(ctx context.Context, users []*table.UserModel) (err error) {
	return u.database.Update(ctx, users)
}
func (u *UserController) UpdateByMap(ctx context.Context, userID string, args map[string]interface{}) (err error) {
	return u.database.UpdateByMap(ctx, userID, args)
}
func (u *UserController) GetByName(ctx context.Context, userName string, showNumber, pageNumber int32) (users []*table.UserModel, count int64, err error) {
	return u.database.GetByName(ctx, userName, showNumber, pageNumber)
}
func (u *UserController) GetByNameAndID(ctx context.Context, content string, showNumber, pageNumber int32) (users []*table.UserModel, count int64, err error) {
	return u.database.GetByNameAndID(ctx, content, showNumber, pageNumber)
}
func (u *UserController) Get(ctx context.Context, showNumber, pageNumber int32) (users []*table.UserModel, count int64, err error) {
	return u.database.Get(ctx, showNumber, pageNumber)
}
func NewUserController(db *gorm.DB) *UserController {
	controller := &UserController{database: newUserDatabase(db)}
	return controller
}

type UserDatabaseInterface interface {
	Find(ctx context.Context, userIDs []string) (users []*table.UserModel, err error)
	Create(ctx context.Context, users []*table.UserModel) error
	Take(ctx context.Context, userID string) (user *table.UserModel, err error)
	Update(ctx context.Context, users []*table.UserModel) (err error)
	UpdateByMap(ctx context.Context, userID string, args map[string]interface{}) (err error)
	GetByName(ctx context.Context, userName string, showNumber, pageNumber int32) (users []*table.UserModel, count int64, err error)
	GetByNameAndID(ctx context.Context, content string, showNumber, pageNumber int32) (users []*table.UserModel, count int64, err error)
	Get(ctx context.Context, showNumber, pageNumber int32) (users []*table.UserModel, count int64, err error)
}

type UserDatabase struct {
	sqlDB *relation.UserGorm
}

func newUserDatabase(db *gorm.DB) *UserDatabase {
	sqlDB := relation.NewUserGorm(db)
	database := &UserDatabase{
		sqlDB: sqlDB,
	}
	return database
}

func (u *UserDatabase) Find(ctx context.Context, userIDs []string) (users []*table.UserModel, err error) {
	return u.sqlDB.Find(ctx, userIDs)
}

func (u *UserDatabase) Create(ctx context.Context, users []*table.UserModel) error {
	return u.sqlDB.Create(ctx, users)
}
func (u *UserDatabase) Take(ctx context.Context, userID string) (user *table.UserModel, err error) {
	return u.sqlDB.Take(ctx, userID)
}
func (u *UserDatabase) Update(ctx context.Context, users []*table.UserModel) (err error) {
	return u.sqlDB.Update(ctx, users)
}
func (u *UserDatabase) UpdateByMap(ctx context.Context, userID string, args map[string]interface{}) (err error) {
	return u.sqlDB.UpdateByMap(ctx, userID, args)
}
func (u *UserDatabase) GetByName(ctx context.Context, userName string, showNumber, pageNumber int32) (users []*table.UserModel, count int64, err error) {
	return u.sqlDB.GetByName(ctx, userName, showNumber, pageNumber)
}
func (u *UserDatabase) GetByNameAndID(ctx context.Context, content string, showNumber, pageNumber int32) (users []*table.UserModel, count int64, err error) {
	return u.sqlDB.GetByNameAndID(ctx, content, showNumber, pageNumber)
}
func (u *UserDatabase) Get(ctx context.Context, showNumber, pageNumber int32) (users []*table.UserModel, count int64, err error) {
	return u.sqlDB.Get(ctx, showNumber, pageNumber)
}
