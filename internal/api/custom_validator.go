package api

import (
	"github.com/go-playground/validator/v10"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
)

func RequiredIf(fl validator.FieldLevel) bool {
	sessionType := fl.Parent().FieldByName("SessionType").Int()
	switch sessionType {
	case constant.SingleChatType, constant.NotificationChatType:
		if fl.FieldName() == "RecvID" {
			return fl.Field().String() != ""
		}
	case constant.GroupChatType, constant.SuperGroupChatType:
		if fl.FieldName() == "GroupID" {
			return fl.Field().String() != ""
		}
	default:
		return true
	}
	return true
}
