package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/openimsdk/protocol/constant"
)

// RequiredIf validates if the specified field is required based on the session type.
func RequiredIf(fl validator.FieldLevel) bool {
	sessionType := fl.Parent().FieldByName("SessionType").Int()

	switch sessionType {
	case constant.SingleChatType, constant.NotificationChatType:
		return fl.FieldName() != "RecvID" || fl.Field().String() != ""
	case constant.WriteGroupChatType, constant.ReadGroupChatType:
		return fl.FieldName() != "GroupID" || fl.Field().String() != ""
	default:
		return true
	}
}
