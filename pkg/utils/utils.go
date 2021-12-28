package utils

import (
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/token_verify"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"time"
)

// copy a by b  b->a
func CopyStructFields(a interface{}, b interface{}, fields ...string) (err error) {
	at := reflect.TypeOf(a)
	av := reflect.ValueOf(a)
	bt := reflect.TypeOf(b)
	bv := reflect.ValueOf(b)

	if at.Kind() != reflect.Ptr {
		err = fmt.Errorf("a must be a struct pointer")
		return err
	}
	av = reflect.ValueOf(av.Interface())

	_fields := make([]string, 0)
	if len(fields) > 0 {
		_fields = fields
	} else {
		for i := 0; i < bv.NumField(); i++ {
			_fields = append(_fields, bt.Field(i).Name)
		}
	}

	if len(_fields) == 0 {
		err = fmt.Errorf("no fields to copy")
		return err
	}

	for i := 0; i < len(_fields); i++ {
		name := _fields[i]
		f := av.Elem().FieldByName(name)
		bValue := bv.FieldByName(name)

		if f.IsValid() && f.Kind() == bValue.Kind() {
			f.Set(bValue)
		}
	}
	return nil
}

func OperationIDGenerator() string {
	return strconv.FormatInt(time.Now().UnixNano()+int64(rand.Uint32()), 10)
}

func FriendOpenIMCopyDB(dst *imdb.Friend, src open_im_sdk.FriendInfo) {
	CopyStructFields(dst, src)
	dst.FriendUserID = src.FriendUser.UserID
}

func FriendDBCopyOpenIM(dst *open_im_sdk.FriendInfo, src imdb.Friend) {
	CopyStructFields(dst, src)
	user, _ := imdb.GetUserByUserID(src.FriendUserID)
	if user != nil {
		CopyStructFields(dst.FriendUser, user)
	}
}

//
func FriendRequestOpenIMCopyDB(dst *imdb.FriendRequest, src open_im_sdk.FriendRequest) {
	CopyStructFields(dst, src)
}

func FriendRequestDBCopyOpenIM(dst *open_im_sdk.FriendRequest, src imdb.FriendRequest) {
	CopyStructFields(dst, src)
}

func GroupOpenIMCopyDB(dst *imdb.Group, src open_im_sdk.GroupInfo) {
	CopyStructFields(dst, src)
}

func GroupDBCopyOpenIM(dst *open_im_sdk.GroupInfo, src imdb.Group) {
	CopyStructFields(dst, src)
	user, _ := imdb.GetGroupOwnerInfoByGroupID(src.GroupID)
	if user != nil {
		dst.OwnerUserID = user.UserID
	}
	dst.MemberCount = imdb.GetGroupMemberNumByGroupID(src.GroupID)
}

func GroupMemberOpenIMCopyDB(dst *imdb.GroupMember, src open_im_sdk.GroupMemberFullInfo) {
	CopyStructFields(dst, src)
}

func GroupMemberDBCopyOpenIM(dst *open_im_sdk.GroupMemberFullInfo, src imdb.GroupMember) {
	CopyStructFields(dst, src)
	if token_verify.IsMangerUserID(src.UserID) {
		u, _ := imdb.GetUserByUserID(src.UserID)
		if u != nil {
			CopyStructFields(dst, u)
		}
		dst.AppMangerLevel = 1
	}
}

func GroupRequestOpenIMCopyDB(dst *imdb.GroupRequest, src open_im_sdk.GroupRequest) {
	CopyStructFields(dst, src)
}

func GroupRequestDBCopyOpenIM(dst *open_im_sdk.GroupRequest, src imdb.GroupRequest) {
	CopyStructFields(dst, src)
}

func UserOpenIMCopyDB(dst *imdb.User, src open_im_sdk.UserInfo) {
	CopyStructFields(dst, src)
}

func UserDBCopyOpenIM(dst *open_im_sdk.UserInfo, src imdb.User) {
	CopyStructFields(dst, src)
}

func BlackOpenIMCopyDB(dst *imdb.Black, src open_im_sdk.BlackInfo) {
	CopyStructFields(dst, src)
	dst.BlockUserID = src.BlackUserInfo.UserID
}

func BlackDBCopyOpenIM(dst *open_im_sdk.BlackInfo, src imdb.Black) {
	CopyStructFields(dst, src)
	user, _ := imdb.GetUserByUserID(src.BlockUserID)
	if user != nil {
		CopyStructFields(dst.BlackUserInfo, user)
	}
}
