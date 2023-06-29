package controller

import (
	"context"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/tx"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
)

type UserDatabase interface {
	//获取指定用户的信息 如有userID未找到 也返回错误
	FindWithError(ctx context.Context, userIDs []string) (users []*relation.UserModel, err error)
	//获取指定用户的信息 如有userID未找到 不返回错误
	Find(ctx context.Context, userIDs []string) (users []*relation.UserModel, err error)
	//插入多条 外部保证userID 不重复 且在db中不存在
	Create(ctx context.Context, users []*relation.UserModel) (err error)
	//更新（非零值） 外部保证userID存在
	Update(ctx context.Context, user *relation.UserModel) (err error)
	//更新（零值） 外部保证userID存在
	UpdateByMap(ctx context.Context, userID string, args map[string]interface{}) (err error)
	//如果没找到，不返回错误
	Page(ctx context.Context, pageNumber, showNumber int32) (users []*relation.UserModel, count int64, err error)
	//只要有一个存在就为true
	IsExist(ctx context.Context, userIDs []string) (exist bool, err error)
	//获取所有用户ID
	GetAllUserID(ctx context.Context) ([]string, error)
	//函数内部先查询db中是否存在，存在则什么都不做；不存在则插入
	InitOnce(ctx context.Context, users []*relation.UserModel) (err error)
	// 获取用户总数
	CountTotal(ctx context.Context) (int64, error)
	// 获取范围内用户增量
	CountRangeEverydayTotal(ctx context.Context, start time.Time, end time.Time) (map[string]int64, error)
}

type userDatabase struct {
	userDB relation.UserModelInterface
	cache  cache.UserCache
	tx     tx.Tx
}

func NewUserDatabase(userDB relation.UserModelInterface, cache cache.UserCache, tx tx.Tx) UserDatabase {
	return &userDatabase{userDB: userDB, cache: cache, tx: tx}
}

func (u *userDatabase) InitOnce(ctx context.Context, users []*relation.UserModel) (err error) {
	userIDs := utils.Slice(users, func(e *relation.UserModel) string {
		return e.UserID
	})
	result, err := u.userDB.Find(ctx, userIDs)
	if err != nil {
		return err
	}
	miss := utils.SliceAnySub(users, result, func(e *relation.UserModel) string { return e.UserID })
	if len(miss) > 0 {
		_ = u.userDB.Create(ctx, miss)
	}
	return nil
}

// 获取指定用户的信息 如有userID未找到 也返回错误
func (u *userDatabase) FindWithError(ctx context.Context, userIDs []string) (users []*relation.UserModel, err error) {
	users, err = u.cache.GetUsersInfo(ctx, userIDs)
	if err != nil {
		return
	}
	if len(users) != len(userIDs) {
		err = errs.ErrRecordNotFound.Wrap("userID not found")
	}
	return
}

// 获取指定用户的信息 如有userID未找到 不返回错误
func (u *userDatabase) Find(ctx context.Context, userIDs []string) (users []*relation.UserModel, err error) {
	users, err = u.cache.GetUsersInfo(ctx, userIDs)
	return
}

// 插入多条 外部保证userID 不重复 且在db中不存在
func (u *userDatabase) Create(ctx context.Context, users []*relation.UserModel) (err error) {
	if err := u.tx.Transaction(func(tx any) error {
		err = u.userDB.Create(ctx, users)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	var userIDs []string
	for _, user := range users {
		userIDs = append(userIDs, user.UserID)
	}
	return u.cache.DelUsersInfo(userIDs...).ExecDel(ctx)
}

// 更新（非零值） 外部保证userID存在
func (u *userDatabase) Update(ctx context.Context, user *relation.UserModel) (err error) {
	if err := u.userDB.Update(ctx, user); err != nil {
		return err
	}
	return u.cache.DelUsersInfo(user.UserID).ExecDel(ctx)
}

// 更新（零值） 外部保证userID存在
func (u *userDatabase) UpdateByMap(ctx context.Context, userID string, args map[string]interface{}) (err error) {
	if err := u.userDB.UpdateByMap(ctx, userID, args); err != nil {
		return err
	}
	return u.cache.DelUsersInfo(userID).ExecDel(ctx)
}

// 获取，如果没找到，不返回错误
func (u *userDatabase) Page(ctx context.Context, pageNumber, showNumber int32) (users []*relation.UserModel, count int64, err error) {
	return u.userDB.Page(ctx, pageNumber, showNumber)
}

// userIDs是否存在 只要有一个存在就为true
func (u *userDatabase) IsExist(ctx context.Context, userIDs []string) (exist bool, err error) {
	users, err := u.userDB.Find(ctx, userIDs)
	if err != nil {
		return false, err
	}
	if len(users) > 0 {
		return true, nil
	}
	return false, nil
}

func (u *userDatabase) GetAllUserID(ctx context.Context) (userIDs []string, err error) {
	return u.userDB.GetAllUserID(ctx)
}

func (u *userDatabase) CountTotal(ctx context.Context) (count int64, err error) {
	return u.userDB.CountTotal(ctx)
}

func (u *userDatabase) CountRangeEverydayTotal(ctx context.Context, start time.Time, end time.Time) (map[string]int64, error) {
	return u.userDB.CountRangeEverydayTotal(ctx, start, end)
}
