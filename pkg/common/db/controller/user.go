package controller

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/relation"
	relation2 "Open_IM/pkg/common/db/table/relation"
	"context"
	"gorm.io/gorm"
)

type UserInterface interface {
	//获取指定用户的信息 如有userID未找到 也返回错误
	Find(ctx context.Context, userIDs []string) (users []*relation2.UserModel, err error)
	//插入多条 外部保证userID 不重复 且在db中不存在
	Create(ctx context.Context, users []*relation2.UserModel) (err error)
	//更新（非零值） 外部保证userID存在
	Update(ctx context.Context, users []*relation2.UserModel) (err error)
	//更新（零值） 外部保证userID存在
	UpdateByMap(ctx context.Context, userID string, args map[string]interface{}) (err error)
	//获取，如果没找到，不返回错误
	Get(ctx context.Context, pageNumber, showNumber int32) (users []*relation2.UserModel, count int64, err error)
	//userIDs是否存在 只要有一个存在就为true
	IsExist(ctx context.Context, userIDs []string) (exist bool, err error)
}

type UserController struct {
	database UserDatabaseInterface
}

func (u *UserController) Find(ctx context.Context, userIDs []string) (users []*relation2.UserModel, err error) {
	return u.database.Find(ctx, userIDs)
}
func (u *UserController) Create(ctx context.Context, users []*relation2.UserModel) error {
	return u.database.Create(ctx, users)
}

func (u *UserController) Update(ctx context.Context, users []*relation2.UserModel) (err error) {
	return u.database.Update(ctx, users)
}
func (u *UserController) UpdateByMap(ctx context.Context, userID string, args map[string]interface{}) (err error) {
	return u.database.UpdateByMap(ctx, userID, args)
}

func (u *UserController) Get(ctx context.Context, pageNumber, showNumber int32) (users []*relation2.UserModel, count int64, err error) {
	return u.database.Get(ctx, pageNumber, showNumber)
}

func (u *UserController) IsExist(ctx context.Context, userIDs []string) (exist bool, err error) {
	return u.database.IsExist(ctx, userIDs)
}
func NewUserController(db *gorm.DB) *UserController {
	controller := &UserController{database: newUserDatabase(db)}
	return controller
}

type UserDatabaseInterface interface {
	Find(ctx context.Context, userIDs []string) (users []*relation2.UserModel, err error)
	Create(ctx context.Context, users []*relation2.UserModel) error
	Update(ctx context.Context, users []*relation2.UserModel) (err error)
	UpdateByMap(ctx context.Context, userID string, args map[string]interface{}) (err error)
	Get(ctx context.Context, pageNumber, showNumber int32) (users []*relation2.UserModel, count int64, err error)
	IsExist(ctx context.Context, userIDs []string) (exist bool, err error)
}

type UserDatabase struct {
	user *relation.UserGorm
}

func newUserDatabase(db *gorm.DB) *UserDatabase {
	sqlDB := relation.NewUserGorm(db)
	database := &UserDatabase{
		user: sqlDB,
	}
	return database
}

// 获取指定用户的信息 如有userID未找到 也返回错误
func (u *UserDatabase) Find(ctx context.Context, userIDs []string) (users []*relation2.UserModel, err error) {
	users, err = u.user.Find(ctx, userIDs)
	if err != nil {
		return
	}
	if len(users) != len(userIDs) {
		err = constant.ErrRecordNotFound.Wrap()
	}
	return
}

// 插入多条 外部保证userID 不重复 且在db中不存在
func (u *UserDatabase) Create(ctx context.Context, users []*relation2.UserModel) (err error) {
	return u.user.Create(ctx, users)
}

// 更新（非零值） 外部保证userID存在
func (u *UserDatabase) Update(ctx context.Context, users []*relation2.UserModel) (err error) {
	return u.user.Update(ctx, users)
}

// 更新（零值） 外部保证userID存在
func (u *UserDatabase) UpdateByMap(ctx context.Context, userID string, args map[string]interface{}) (err error) {
	return u.user.UpdateByMap(ctx, userID, args)
}

// 获取，如果没找到，不返回错误
func (u *UserDatabase) Get(ctx context.Context, showNumber, pageNumber int32) (users []*relation2.UserModel, count int64, err error) {
	return u.user.Get(ctx, showNumber, pageNumber)
}

// userIDs是否存在 只要有一个存在就为true
func (u *UserDatabase) IsExist(ctx context.Context, userIDs []string) (exist bool, err error) {
	users, err := u.user.Find(ctx, userIDs)
	if err != nil {
		return false, err
	}
	if len(users) > 0 {
		return true, nil
	}
	return false, nil
}
