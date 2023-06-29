package controller

import (
	"context"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/tx"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"gorm.io/gorm"
)

type FriendDatabase interface {
	// 检查user2是否在user1的好友列表中(inUser1Friends==true) 检查user1是否在user2的好友列表中(inUser2Friends==true)
	CheckIn(ctx context.Context, user1, user2 string) (inUser1Friends bool, inUser2Friends bool, err error)
	// 增加或者更新好友申请
	AddFriendRequest(ctx context.Context, fromUserID, toUserID string, reqMsg string, ex string) (err error)
	// 先判断是否在好友表，如果在则不插入
	BecomeFriends(ctx context.Context, ownerUserID string, friendUserIDs []string, addSource int32) (err error)
	// 拒绝好友申请
	RefuseFriendRequest(ctx context.Context, friendRequest *relation.FriendRequestModel) (err error)
	// 同意好友申请
	AgreeFriendRequest(ctx context.Context, friendRequest *relation.FriendRequestModel) (err error)
	// 删除好友
	Delete(ctx context.Context, ownerUserID string, friendUserIDs []string) (err error)
	// 更新好友备注
	UpdateRemark(ctx context.Context, ownerUserID, friendUserID, remark string) (err error)
	// 获取ownerUserID的好友列表
	PageOwnerFriends(ctx context.Context, ownerUserID string, pageNumber, showNumber int32) (friends []*relation.FriendModel, total int64, err error)
	// friendUserID在哪些人的好友列表中
	PageInWhoseFriends(ctx context.Context, friendUserID string, pageNumber, showNumber int32) (friends []*relation.FriendModel, total int64, err error)
	// 获取我发出去的好友申请
	PageFriendRequestFromMe(ctx context.Context, userID string, pageNumber, showNumber int32) (friends []*relation.FriendRequestModel, total int64, err error)
	// 获取我收到的的好友申请
	PageFriendRequestToMe(ctx context.Context, userID string, pageNumber, showNumber int32) (friends []*relation.FriendRequestModel, total int64, err error)
	// 获取某人指定好友的信息
	FindFriendsWithError(ctx context.Context, ownerUserID string, friendUserIDs []string) (friends []*relation.FriendModel, err error)
	FindFriendUserIDs(ctx context.Context, ownerUserID string) (friendUserIDs []string, err error)
}

type friendDatabase struct {
	friend        relation.FriendModelInterface
	friendRequest relation.FriendRequestModelInterface
	tx            tx.Tx
	cache         cache.FriendCache
}

func NewFriendDatabase(friend relation.FriendModelInterface, friendRequest relation.FriendRequestModelInterface, cache cache.FriendCache, tx tx.Tx) FriendDatabase {
	return &friendDatabase{friend: friend, friendRequest: friendRequest, cache: cache, tx: tx}
}

// ok 检查user2是否在user1的好友列表中(inUser1Friends==true) 检查user1是否在user2的好友列表中(inUser2Friends==true)
func (f *friendDatabase) CheckIn(ctx context.Context, userID1, userID2 string) (inUser1Friends bool, inUser2Friends bool, err error) {
	userID1FriendIDs, err := f.cache.GetFriendIDs(ctx, userID1)
	if err != nil {
		return
	}
	userID2FriendIDs, err := f.cache.GetFriendIDs(ctx, userID2)
	if err != nil {
		return
	}
	return utils.IsContain(userID2, userID1FriendIDs), utils.IsContain(userID1, userID2FriendIDs), nil
}

// 增加或者更新好友申请 如果之前有记录则更新，没有记录则新增
func (f *friendDatabase) AddFriendRequest(ctx context.Context, fromUserID, toUserID string, reqMsg string, ex string) (err error) {
	return f.tx.Transaction(func(tx any) error {
		_, err := f.friendRequest.NewTx(tx).Take(ctx, fromUserID, toUserID)
		//有db错误
		if err != nil && errs.Unwrap(err) != gorm.ErrRecordNotFound {
			return err
		}
		//无错误 则更新
		if err == nil {
			m := make(map[string]interface{}, 1)
			m["handle_result"] = 0
			m["handle_msg"] = ""
			m["req_msg"] = reqMsg
			m["ex"] = ex
			m["create_time"] = time.Now()
			if err := f.friendRequest.NewTx(tx).UpdateByMap(ctx, fromUserID, toUserID, m); err != nil {
				return err
			}
			return nil
		}
		//gorm.ErrRecordNotFound 错误，则新增
		if err := f.friendRequest.NewTx(tx).Create(ctx, []*relation.FriendRequestModel{{FromUserID: fromUserID, ToUserID: toUserID, ReqMsg: reqMsg, Ex: ex, CreateTime: time.Now(), HandleTime: time.Unix(0, 0)}}); err != nil {
			return err
		}
		return nil
	})
}

// (1)先判断是否在好友表 （在不在都不返回错误） (2)对于不在好友列表的 插入即可
func (f *friendDatabase) BecomeFriends(ctx context.Context, ownerUserID string, friendUserIDs []string, addSource int32) (err error) {
	cache := f.cache.NewCache()
	if err := f.tx.Transaction(func(tx any) error {
		//先find 找出重复的 去掉重复的
		fs1, err := f.friend.NewTx(tx).FindFriends(ctx, ownerUserID, friendUserIDs)
		if err != nil {
			return err
		}
		opUserID := mcontext.GetOperationID(ctx)
		for _, v := range friendUserIDs {
			fs1 = append(fs1, &relation.FriendModel{OwnerUserID: ownerUserID, FriendUserID: v, AddSource: addSource, OperatorUserID: opUserID})
		}
		fs11 := utils.DistinctAny(fs1, func(e *relation.FriendModel) string {
			return e.FriendUserID
		})

		err = f.friend.NewTx(tx).Create(ctx, fs11)
		if err != nil {
			return err
		}
		fs2, err := f.friend.NewTx(tx).FindReversalFriends(ctx, ownerUserID, friendUserIDs)
		if err != nil {
			return err
		}
		var newFriendIDs []string
		for _, v := range friendUserIDs {
			fs2 = append(fs2, &relation.FriendModel{OwnerUserID: v, FriendUserID: ownerUserID, AddSource: addSource, OperatorUserID: opUserID})
			newFriendIDs = append(newFriendIDs, v)
		}
		fs22 := utils.DistinctAny(fs2, func(e *relation.FriendModel) string {
			return e.OwnerUserID
		})
		err = f.friend.NewTx(tx).Create(ctx, fs22)
		if err != nil {
			return err
		}
		newFriendIDs = append(newFriendIDs, ownerUserID)
		cache = cache.DelFriendIDs(newFriendIDs...)
		return nil
	}); err != nil {
		return nil
	}
	return cache.ExecDel(ctx)
}

// 拒绝好友申请 (1)检查是否有申请记录且为未处理状态 （没有记录返回错误） (2)修改申请记录 已拒绝
func (f *friendDatabase) RefuseFriendRequest(ctx context.Context, friendRequest *relation.FriendRequestModel) (err error) {
	fr, err := f.friendRequest.Take(ctx, friendRequest.FromUserID, friendRequest.ToUserID)
	if err != nil {
		return err
	}
	if fr.HandleResult != 0 {
		return errs.ErrArgs.Wrap("the friend request has been processed")
	}
	friendRequest.HandleResult = constant.FriendResponseRefuse
	friendRequest.HandleTime = time.Now()
	err = f.friendRequest.Update(ctx, friendRequest)
	if err != nil {
		return err
	}
	return nil
}

// AgreeFriendRequest 同意好友申请  (1)检查是否有申请记录且为未处理状态 （没有记录返回错误） (2)检查是否好友（不返回错误）   (3) 建立双向好友关系（存在的忽略）
func (f *friendDatabase) AgreeFriendRequest(ctx context.Context, friendRequest *relation.FriendRequestModel) (err error) {
	return f.tx.Transaction(func(tx any) error {
		fr, err := f.friendRequest.NewTx(tx).Take(ctx, friendRequest.FromUserID, friendRequest.ToUserID)
		if err != nil {
			return err
		}
		if fr.HandleResult != 0 {
			return errs.ErrArgs.Wrap("the friend request has been processed")
		}
		friendRequest.HandlerUserID = mcontext.GetOpUserID(ctx)
		friendRequest.HandleResult = constant.FriendResponseAgree
		friendRequest.HandleTime = time.Now()
		err = f.friendRequest.NewTx(tx).Update(ctx, friendRequest)
		if err != nil {
			return err
		}
		exists, err := f.friend.NewTx(tx).FindUserState(ctx, friendRequest.FromUserID, friendRequest.ToUserID)
		if err != nil {
			return err
		}
		existsMap := utils.SliceSet(utils.Slice(exists, func(friend *relation.FriendModel) [2]string {
			return [...]string{friend.OwnerUserID, friend.FriendUserID} // 自己 - 好友
		}))
		var adds []*relation.FriendModel
		if _, ok := existsMap[[...]string{friendRequest.ToUserID, friendRequest.FromUserID}]; !ok { // 自己 - 好友
			adds = append(adds, &relation.FriendModel{OwnerUserID: friendRequest.ToUserID, FriendUserID: friendRequest.FromUserID, AddSource: int32(constant.BecomeFriendByApply), OperatorUserID: friendRequest.FromUserID})
		}
		if _, ok := existsMap[[...]string{friendRequest.FromUserID, friendRequest.ToUserID}]; !ok { // 好友 - 自己
			adds = append(adds, &relation.FriendModel{OwnerUserID: friendRequest.FromUserID, FriendUserID: friendRequest.ToUserID, AddSource: int32(constant.BecomeFriendByApply), OperatorUserID: friendRequest.FromUserID})
		}
		if len(adds) > 0 {
			if err := f.friend.NewTx(tx).Create(ctx, adds); err != nil {
				return err
			}
		}
		return f.cache.DelFriendIDs(friendRequest.ToUserID, friendRequest.FromUserID).ExecDel(ctx)
	})
}

// 删除好友  外部判断是否好友关系
func (f *friendDatabase) Delete(ctx context.Context, ownerUserID string, friendUserIDs []string) (err error) {
	if err := f.friend.Delete(ctx, ownerUserID, friendUserIDs); err != nil {
		return err
	}
	return f.cache.DelFriendIDs(append(friendUserIDs, ownerUserID)...).ExecDel(ctx)
}

// 更新好友备注 零值也支持
func (f *friendDatabase) UpdateRemark(ctx context.Context, ownerUserID, friendUserID, remark string) (err error) {
	if err := f.friend.UpdateRemark(ctx, ownerUserID, friendUserID, remark); err != nil {
		return err
	}
	return f.cache.DelFriend(ownerUserID, friendUserID).ExecDel(ctx)
}

// 获取ownerUserID的好友列表 无结果不返回错误
func (f *friendDatabase) PageOwnerFriends(ctx context.Context, ownerUserID string, pageNumber, showNumber int32) (friends []*relation.FriendModel, total int64, err error) {
	return f.friend.FindOwnerFriends(ctx, ownerUserID, pageNumber, showNumber)
}

// friendUserID在哪些人的好友列表中
func (f *friendDatabase) PageInWhoseFriends(ctx context.Context, friendUserID string, pageNumber, showNumber int32) (friends []*relation.FriendModel, total int64, err error) {
	return f.friend.FindInWhoseFriends(ctx, friendUserID, pageNumber, showNumber)
}

// 获取我发出去的好友申请  无结果不返回错误
func (f *friendDatabase) PageFriendRequestFromMe(ctx context.Context, userID string, pageNumber, showNumber int32) (friends []*relation.FriendRequestModel, total int64, err error) {
	return f.friendRequest.FindFromUserID(ctx, userID, pageNumber, showNumber)
}

// 获取我收到的的好友申请 无结果不返回错误
func (f *friendDatabase) PageFriendRequestToMe(ctx context.Context, userID string, pageNumber, showNumber int32) (friends []*relation.FriendRequestModel, total int64, err error) {
	return f.friendRequest.FindToUserID(ctx, userID, pageNumber, showNumber)
}

// 获取某人指定好友的信息 如果有好友不存在，也返回错误
func (f *friendDatabase) FindFriendsWithError(ctx context.Context, ownerUserID string, friendUserIDs []string) (friends []*relation.FriendModel, err error) {
	friends, err = f.friend.FindFriends(ctx, ownerUserID, friendUserIDs)
	if err != nil {
		return
	}
	if len(friends) != len(friendUserIDs) {
		err = errs.ErrRecordNotFound.Wrap()
	}
	return
}

func (f *friendDatabase) FindFriendUserIDs(ctx context.Context, ownerUserID string) (friendUserIDs []string, err error) {
	return f.cache.GetFriendIDs(ctx, ownerUserID)
}
