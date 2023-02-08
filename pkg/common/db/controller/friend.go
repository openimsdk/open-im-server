package controller

import (
	"Open_IM/pkg/common/constant"
	relation1 "Open_IM/pkg/common/db/relation"
	"Open_IM/pkg/common/db/table/relation"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"gorm.io/gorm"
)

type FriendInterface interface {
	//  检查user2是否在user1的好友列表中(inUser1Friends==true) 检查user1是否在user2的好友列表中(inUser2Friends==true)
	CheckIn(ctx context.Context, user1, user2 string) (inUser1Friends bool, inUser2Friends bool, err error)
	//  增加或者更新好友申请 如果之前有记录则更新，没有记录则新增
	AddFriendRequest(ctx context.Context, fromUserID, toUserID string, reqMsg string, ex string) (err error)
	//  (1)先判断是否在好友表 （在不在都不返回错误） (2)对于不在好友列表的 插入即可
	BecomeFriends(ctx context.Context, ownerUserID string, friendUserIDs []string, addSource int32, OperatorUserID string) (err error)
	//  拒绝好友申请 (1)检查是否有申请记录且为未处理状态 （没有记录返回错误） (2)修改申请记录 已拒绝
	RefuseFriendRequest(ctx context.Context, friendRequest *relation.FriendRequestModel) (err error)
	//  同意好友申请  (1)检查是否有申请记录且为未处理状态 （没有记录返回错误） (2)检查是否好友（不返回错误）   (3) 不是好友则建立双向好友关系  （4）修改申请记录 已同意
	AgreeFriendRequest(ctx context.Context, friendRequest *relation.FriendRequestModel) (err error)
	//  删除好友  外部判断是否好友关系
	Delete(ctx context.Context, ownerUserID string, friendUserIDs []string) (err error)
	//  更新好友备注 零值也支持
	UpdateRemark(ctx context.Context, ownerUserID, friendUserID, remark string) (err error)
	//  获取ownerUserID的好友列表 无结果不返回错误
	FindOwnerFriends(ctx context.Context, ownerUserID string, pageNumber, showNumber int32) (friends []*relation.FriendModel, total int64, err error)
	//  friendUserID在哪些人的好友列表中
	FindInWhoseFriends(ctx context.Context, friendUserID string, pageNumber, showNumber int32) (friends []*relation.FriendModel, total int64, err error)
	//  获取我发出去的好友申请  无结果不返回错误
	FindFriendRequestFromMe(ctx context.Context, userID string, pageNumber, showNumber int32) (friends []*relation.FriendRequestModel, total int64, err error)
	//  获取我收到的的好友申请 无结果不返回错误
	FindFriendRequestToMe(ctx context.Context, userID string, pageNumber, showNumber int32) (friends []*relation.FriendRequestModel, total int64, err error)
	//  获取某人指定好友的信息 如果有一个不存在也返回错误
	FindFriends(ctx context.Context, ownerUserID string, friendUserIDs []string) (friends []*relation.FriendModel, err error)
}

type FriendController struct {
	database FriendDatabaseInterface
}

func NewFriendController(db *gorm.DB) *FriendController {
	return &FriendController{database: NewFriendDatabase(db)}
}

// 检查user2是否在user1的好友列表中(inUser1Friends==true) 检查user1是否在user2的好友列表中(inUser2Friends==true)
func (f *FriendController) CheckIn(ctx context.Context, user1, user2 string) (inUser1Friends bool, inUser2Friends bool, err error) {
	return f.database.CheckIn(ctx, user1, user2)
}

// AddFriendRequest 增加或者更新好友申请
func (f *FriendController) AddFriendRequest(ctx context.Context, fromUserID, toUserID string, reqMsg string, ex string) (err error) {
	return f.database.AddFriendRequest(ctx, fromUserID, toUserID, reqMsg, ex)
}

// BecomeFriend 先判断是否在好友表，如果在则不插入
func (f *FriendController) BecomeFriends(ctx context.Context, ownerUserID string, friendUserIDs []string, addSource int32, OperatorUserID string) (err error) {
	return f.database.BecomeFriends(ctx, ownerUserID, friendUserIDs, addSource, OperatorUserID)
}

// RefuseFriendRequest 拒绝好友申请
func (f *FriendController) RefuseFriendRequest(ctx context.Context, friendRequest *relation.FriendRequestModel) (err error) {
	return f.database.RefuseFriendRequest(ctx, friendRequest)
}

// AgreeFriendRequest 同意好友申请
func (f *FriendController) AgreeFriendRequest(ctx context.Context, friendRequest *relation.FriendRequestModel) (err error) {
	return f.database.AgreeFriendRequest(ctx, friendRequest)
}

// Delete 删除好友
func (f *FriendController) Delete(ctx context.Context, ownerUserID string, friendUserIDs []string) (err error) {
	return f.database.Delete(ctx, ownerUserID, friendUserIDs)
}

// UpdateRemark 更新好友备注
func (f *FriendController) UpdateRemark(ctx context.Context, ownerUserID, friendUserID, remark string) (err error) {
	return f.database.UpdateRemark(ctx, ownerUserID, friendUserID, remark)
}

// FindOwnerFriends 获取ownerUserID的好友列表
func (f *FriendController) FindOwnerFriends(ctx context.Context, ownerUserID string, pageNumber, showNumber int32) (friends []*relation.FriendModel, total int64, err error) {
	return f.database.FindOwnerFriends(ctx, ownerUserID, pageNumber, showNumber)
}

// FindInWhoseFriends friendUserID在哪些人的好友列表中
func (f *FriendController) FindInWhoseFriends(ctx context.Context, friendUserID string, pageNumber, showNumber int32) (friends []*relation.FriendModel, total int64, err error) {
	return f.database.FindInWhoseFriends(ctx, friendUserID, pageNumber, showNumber)
}

// FindFriendRequestFromMe 获取我发出去的好友申请
func (f *FriendController) FindFriendRequestFromMe(ctx context.Context, userID string, pageNumber, showNumber int32) (friends []*relation.FriendRequestModel, total int64, err error) {
	return f.database.FindFriendRequestFromMe(ctx, userID, pageNumber, showNumber)
}

// FindFriendRequestToMe 获取我收到的的好友申请
func (f *FriendController) FindFriendRequestToMe(ctx context.Context, userID string, pageNumber, showNumber int32) (friends []*relation.FriendRequestModel, total int64, err error) {
	return f.database.FindFriendRequestToMe(ctx, userID, pageNumber, showNumber)
}

// FindFriends 获取某人指定好友的信息
func (f *FriendController) FindFriends(ctx context.Context, ownerUserID string, friendUserIDs []string) (friends []*relation.FriendModel, err error) {
	return f.database.FindFriends(ctx, ownerUserID, friendUserIDs)
}

type FriendDatabaseInterface interface {
	// 检查user2是否在user1的好友列表中(inUser1Friends==true) 检查user1是否在user2的好友列表中(inUser2Friends==true)
	CheckIn(ctx context.Context, user1, user2 string) (inUser1Friends bool, inUser2Friends bool, err error)
	// 增加或者更新好友申请
	AddFriendRequest(ctx context.Context, fromUserID, toUserID string, reqMsg string, ex string) (err error)
	// 先判断是否在好友表，如果在则不插入
	BecomeFriends(ctx context.Context, ownerUserID string, friendUserIDs []string, addSource int32, OperatorUserID string) (err error)
	// 拒绝好友申请
	RefuseFriendRequest(ctx context.Context, friendRequest *relation.FriendRequestModel) (err error)
	// 同意好友申请
	AgreeFriendRequest(ctx context.Context, friendRequest *relation.FriendRequestModel) (err error)
	// 删除好友
	Delete(ctx context.Context, ownerUserID string, friendUserIDs []string) (err error)
	// 更新好友备注
	UpdateRemark(ctx context.Context, ownerUserID, friendUserID, remark string) (err error)
	// 获取ownerUserID的好友列表
	FindOwnerFriends(ctx context.Context, ownerUserID string, pageNumber, showNumber int32) (friends []*relation.FriendModel, total int64, err error)
	// friendUserID在哪些人的好友列表中
	FindInWhoseFriends(ctx context.Context, friendUserID string, pageNumber, showNumber int32) (friends []*relation.FriendModel, total int64, err error)
	// 获取我发出去的好友申请
	FindFriendRequestFromMe(ctx context.Context, userID string, pageNumber, showNumber int32) (friends []*relation.FriendRequestModel, total int64, err error)
	// 获取我收到的的好友申请
	FindFriendRequestToMe(ctx context.Context, userID string, pageNumber, showNumber int32) (friends []*relation.FriendRequestModel, total int64, err error)
	// 获取某人指定好友的信息
	FindFriends(ctx context.Context, ownerUserID string, friendUserIDs []string) (friends []*relation.FriendModel, err error)
}

type FriendDatabase struct {
	friend        *relation1.FriendGorm
	friendRequest *relation1.FriendRequestGorm
}

func NewFriendDatabase(db *gorm.DB) *FriendDatabase {
	return &FriendDatabase{friend: relation1.NewFriendGorm(db), friendRequest: relation1.NewFriendRequestGorm(db)}
}

// ok 检查user2是否在user1的好友列表中(inUser1Friends==true) 检查user1是否在user2的好友列表中(inUser2Friends==true)
func (f *FriendDatabase) CheckIn(ctx context.Context, userID1, userID2 string) (inUser1Friends bool, inUser2Friends bool, err error) {
	friends, err := f.friend.FindUserState(ctx, userID1, userID2)
	if err != nil {
		return false, false, err
	}
	for _, v := range friends {
		if v.OwnerUserID == userID1 && v.FriendUserID == userID2 {
			inUser1Friends = true
		}
		if v.OwnerUserID == userID2 && v.FriendUserID == userID1 {
			inUser2Friends = true
		}
	}
	return
}

// 增加或者更新好友申请 如果之前有记录则更新，没有记录则新增
func (f *FriendDatabase) AddFriendRequest(ctx context.Context, fromUserID, toUserID string, reqMsg string, ex string) (err error) {
	return f.friendRequest.DB.Transaction(func(tx *gorm.DB) error {
		_, err := f.friendRequest.Take(ctx, fromUserID, toUserID, tx)
		//有db错误
		if err != nil && errors.Unwrap(err) != gorm.ErrRecordNotFound {
			return err
		}
		//无错误 则更新
		if err == nil {
			m := make(map[string]interface{}, 1)
			m["handle_result"] = 0
			m["handle_msg"] = ""
			m["req_msg"] = reqMsg
			m["ex"] = ex
			if err := f.friendRequest.UpdateByMap(ctx, fromUserID, toUserID, m, tx); err != nil {
				return err
			}
			return nil
		}
		//gorm.ErrRecordNotFound 错误，则新增
		if err := f.friendRequest.Create(ctx, []*relation.FriendRequestModel{&relation.FriendRequestModel{FromUserID: fromUserID, ToUserID: toUserID, ReqMsg: reqMsg, Ex: ex}}, tx); err != nil {
			return err
		}
		return nil
	})
}

// (1)先判断是否在好友表 （在不在都不返回错误） (2)对于不在好友列表的 插入即可
func (f *FriendDatabase) BecomeFriends(ctx context.Context, ownerUserID string, friendUserIDs []string, addSource int32, OperatorUserID string) (err error) {
	return f.friend.DB.Transaction(func(tx *gorm.DB) error {
		//先find 找出重复的 去掉重复的
		fs1, err := f.friend.FindFriends(ctx, ownerUserID, friendUserIDs, tx)
		if err != nil {
			return err
		}
		for _, v := range friendUserIDs {
			fs1 = append(fs1, &relation.FriendModel{OwnerUserID: ownerUserID, FriendUserID: v, AddSource: addSource, OperatorUserID: OperatorUserID})
		}
		fs11 := utils.DistinctAny(fs1, func(e *relation.FriendModel) string {
			return e.FriendUserID
		})

		err = f.friend.Create(ctx, fs11, tx)
		if err != nil {
			return err
		}

		fs2, err := f.friend.FindReversalFriends(ctx, ownerUserID, friendUserIDs, tx)
		if err != nil {
			return err
		}
		for _, v := range friendUserIDs {
			fs2 = append(fs2, &relation.FriendModel{OwnerUserID: v, FriendUserID: ownerUserID, AddSource: addSource, OperatorUserID: OperatorUserID})
		}
		fs22 := utils.DistinctAny(fs2, func(e *relation.FriendModel) string {
			return e.OwnerUserID
		})
		err = f.friend.Create(ctx, fs22, tx)
		if err != nil {
			return err
		}
		return nil
	})
}

// 拒绝好友申请 (1)检查是否有申请记录且为未处理状态 （没有记录返回错误） (2)修改申请记录 已拒绝
func (f *FriendDatabase) RefuseFriendRequest(ctx context.Context, friendRequest *relation.FriendRequestModel) (err error) {
	_, err = f.friendRequest.Take(ctx, friendRequest.FromUserID, friendRequest.ToUserID)
	if err != nil {
		return err
	}
	friendRequest.HandleResult = constant.FriendResponseRefuse
	err = f.friendRequest.Update(ctx, []*relation.FriendRequestModel{friendRequest})
	if err != nil {
		return err
	}
	return nil
}

// 同意好友申请  (1)检查是否有申请记录且为未处理状态 （没有记录返回错误） (2)检查是否好友（不返回错误）   (3) 不是好友则建立双向好友关系  （4）修改申请记录 已同意
func (f *FriendDatabase) AgreeFriendRequest(ctx context.Context, friendRequest *relation.FriendRequestModel) (err error) {
	return f.friend.DB.Transaction(func(tx *gorm.DB) error {
		_, err = f.friendRequest.Take(ctx, friendRequest.FromUserID, friendRequest.ToUserID)
		if err != nil {
			return err
		}
		friendRequest.HandlerUserID = friendRequest.FromUserID
		friendRequest.HandleResult = constant.FriendResponseAgree
		err = f.friendRequest.Update(ctx, []*relation.FriendRequestModel{friendRequest}, tx)
		if err != nil {
			return err
		}

		ownerUserID := friendRequest.FromUserID
		friendUserIDs := []string{friendRequest.ToUserID}
		addSource := int32(constant.BecomeFriendByApply)
		OperatorUserID := friendRequest.FromUserID
		//先find 找出重复的 去掉重复的
		fs1, err := f.friend.FindFriends(ctx, ownerUserID, friendUserIDs, tx)
		if err != nil {
			return err
		}
		for _, v := range friendUserIDs {
			fs1 = append(fs1, &relation.FriendModel{OwnerUserID: ownerUserID, FriendUserID: v, AddSource: addSource, OperatorUserID: OperatorUserID})
		}
		fs11 := utils.DistinctAny(fs1, func(e *relation.FriendModel) string {
			return e.FriendUserID
		})

		err = f.friend.Create(ctx, fs11, tx)
		if err != nil {
			return err
		}

		fs2, err := f.friend.FindReversalFriends(ctx, ownerUserID, friendUserIDs, tx)
		if err != nil {
			return err
		}
		for _, v := range friendUserIDs {
			fs2 = append(fs2, &relation.FriendModel{OwnerUserID: v, FriendUserID: ownerUserID, AddSource: addSource, OperatorUserID: OperatorUserID})
		}
		fs22 := utils.DistinctAny(fs2, func(e *relation.FriendModel) string {
			return e.OwnerUserID
		})
		err = f.friend.Create(ctx, fs22, tx)
		if err != nil {
			return err
		}
		return nil
	})
}

// 删除好友  外部判断是否好友关系
func (f *FriendDatabase) Delete(ctx context.Context, ownerUserID string, friendUserIDs []string) (err error) {
	return f.friend.Delete(ctx, ownerUserID, friendUserIDs)
}

// 更新好友备注 零值也支持
func (f *FriendDatabase) UpdateRemark(ctx context.Context, ownerUserID, friendUserID, remark string) (err error) {
	return f.friend.UpdateRemark(ctx, ownerUserID, friendUserID, remark)
}

// 获取ownerUserID的好友列表 无结果不返回错误
func (f *FriendDatabase) FindOwnerFriends(ctx context.Context, ownerUserID string, pageNumber, showNumber int32) (friends []*relation.FriendModel, total int64, err error) {
	return f.friend.FindOwnerFriends(ctx, ownerUserID, pageNumber, showNumber)
}

// friendUserID在哪些人的好友列表中
func (f *FriendDatabase) FindInWhoseFriends(ctx context.Context, friendUserID string, pageNumber, showNumber int32) (friends []*relation.FriendModel, total int64, err error) {
	return f.friend.FindInWhoseFriends(ctx, friendUserID, pageNumber, showNumber)
}

// 获取我发出去的好友申请  无结果不返回错误
func (f *FriendDatabase) FindFriendRequestFromMe(ctx context.Context, userID string, pageNumber, showNumber int32) (friends []*relation.FriendRequestModel, total int64, err error) {
	return f.friendRequest.FindFromUserID(ctx, userID, pageNumber, showNumber)
}

// 获取我收到的的好友申请 无结果不返回错误
func (f *FriendDatabase) FindFriendRequestToMe(ctx context.Context, userID string, pageNumber, showNumber int32) (friends []*relation.FriendRequestModel, total int64, err error) {
	return f.friendRequest.FindToUserID(ctx, userID, pageNumber, showNumber)
}

// 获取某人指定好友的信息 如果有好友不存在，也返回错误
func (f *FriendDatabase) FindFriends(ctx context.Context, ownerUserID string, friendUserIDs []string) (friends []*relation.FriendModel, err error) {
	friends, err = f.friend.FindFriends(ctx, ownerUserID, friendUserIDs)
	if err != nil {
		return
	}
	if len(friends) != len(friendUserIDs) {
		err = constant.ErrRecordNotFound.Wrap()
	}
	return
}
