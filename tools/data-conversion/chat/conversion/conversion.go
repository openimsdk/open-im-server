package conversion

import (
	v2 "github.com/openimsdk/open-im-server/v3/tools/data-conversion/chat/v2"
	"github.com/openimsdk/open-im-server/v3/tools/data-conversion/chat/v3/admin"
	"github.com/openimsdk/open-im-server/v3/tools/data-conversion/chat/v3/chat"
)

// ########## chat ##########

func Account(v v2.Account) chat.Account {
	return chat.Account{
		UserID:         v.UserID,
		Password:       v.Password,
		CreateTime:     v.CreateTime,
		ChangeTime:     v.ChangeTime,
		OperatorUserID: v.OperatorUserID,
	}
}

func Attribute(v v2.Attribute) chat.Attribute {
	return chat.Attribute{
		UserID:           v.UserID,
		Account:          v.Account,
		PhoneNumber:      v.PhoneNumber,
		AreaCode:         v.AreaCode,
		Email:            v.Email,
		Nickname:         v.Nickname,
		FaceURL:          v.FaceURL,
		Gender:           v.Gender,
		CreateTime:       v.CreateTime,
		ChangeTime:       v.ChangeTime,
		BirthTime:        v.BirthTime,
		Level:            v.Level,
		AllowVibration:   v.AllowVibration,
		AllowBeep:        v.AllowBeep,
		AllowAddFriend:   v.AllowAddFriend,
		GlobalRecvMsgOpt: 0,
	}
}

func Register(v v2.Register) chat.Register {
	return chat.Register{
		UserID:      v.UserID,
		DeviceID:    v.DeviceID,
		IP:          v.IP,
		Platform:    v.Platform,
		AccountType: v.AccountType,
		Mode:        v.Mode,
		CreateTime:  v.CreateTime,
	}
}

func UserLoginRecord(v v2.UserLoginRecord) chat.UserLoginRecord {
	return chat.UserLoginRecord{
		UserID:    v.UserID,
		LoginTime: v.LoginTime,
		IP:        v.IP,
		DeviceID:  v.DeviceID,
		Platform:  v.Platform,
	}
}

// ########## admin ##########

func Admin(v v2.Admin) admin.Admin {
	return admin.Admin{
		Account:    v.Account,
		Password:   v.Password,
		FaceURL:    v.FaceURL,
		Nickname:   v.Nickname,
		UserID:     v.UserID,
		Level:      v.Level,
		CreateTime: v.CreateTime,
	}
}

func Applet(v v2.Applet) admin.Applet {
	return admin.Applet{
		ID:         v.ID,
		Name:       v.Name,
		AppID:      v.AppID,
		Icon:       v.Icon,
		URL:        v.URL,
		MD5:        v.MD5,
		Size:       v.Size,
		Version:    v.Version,
		Priority:   v.Priority,
		Status:     v.Status,
		CreateTime: v.CreateTime,
	}
}

func ForbiddenAccount(v v2.ForbiddenAccount) admin.ForbiddenAccount {
	return admin.ForbiddenAccount{
		UserID:         v.UserID,
		Reason:         v.Reason,
		OperatorUserID: v.OperatorUserID,
		CreateTime:     v.CreateTime,
	}
}

func InvitationRegister(v v2.InvitationRegister) admin.InvitationRegister {
	return admin.InvitationRegister{
		InvitationCode: v.InvitationCode,
		UsedByUserID:   v.UsedByUserID,
		CreateTime:     v.CreateTime,
	}
}

func IPForbidden(v v2.IPForbidden) admin.IPForbidden {
	return admin.IPForbidden{
		IP:            v.IP,
		LimitRegister: v.LimitRegister > 0,
		LimitLogin:    v.LimitLogin > 0,
		CreateTime:    v.CreateTime,
	}
}

func LimitUserLoginIP(v v2.LimitUserLoginIP) admin.LimitUserLoginIP {
	return admin.LimitUserLoginIP{
		UserID:     v.UserID,
		IP:         v.IP,
		CreateTime: v.CreateTime,
	}
}

func RegisterAddFriend(v v2.RegisterAddFriend) admin.RegisterAddFriend {
	return admin.RegisterAddFriend{
		UserID:     v.UserID,
		CreateTime: v.CreateTime,
	}
}

func RegisterAddGroup(v v2.RegisterAddGroup) admin.RegisterAddGroup {
	return admin.RegisterAddGroup{
		GroupID:    v.GroupID,
		CreateTime: v.CreateTime,
	}
}
