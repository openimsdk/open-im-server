package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"errors"
	"math/rand"
	"time"
)

/**
 * 批量生成邀请码
 */
func BatchCreateInvitationCodes(CodeNums int, CodeLen int) error {
	i := CodeNums
	for {
		if i == 0 {
			break
		}
		invitation := new(db.Invitation)
		invitation.CreateTime = time.Now()
		invitation.InvitationCode = CreateRandomString(CodeLen)
		invitation.LastTime = time.Now()
		invitation.Status = 0
		invitation.UserID = ""
		result := db.DB.MysqlDB.DefaultGormDB().Table("invitations").Create(&invitation)
		if result.Error != nil {
			continue
		}
		if result.RowsAffected > 0 {
			i = i - 1
		}
	}
	return nil
}

/**
 * 检查邀请码
 */
func CheckInvitationCode(code string) error {
	var invitationCode db.Invitation
	err := db.DB.MysqlDB.DefaultGormDB().Table("invitations").Where("invitation_code=?", code).Take(&invitationCode).Error
	if err != nil {
		return err
	}
	if invitationCode.InvitationCode != code {
		return errors.New("邀请码不存在")
	}
	if invitationCode.Status != 0 {
		return errors.New("邀请码已经被使用")
	}
	return nil
}

/**
 * 尝试加锁模式解决邀请码抢占的问题
 */
func TryLockInvitationCode(Code string, UserId string) bool {
	Data := make(map[string]interface{}, 0)
	Data["user_id"] = UserId
	Data["status"] = 1
	Data["last_time"] = time.Now()
	result := db.DB.MysqlDB.DefaultGormDB().Table("invitations").Where("invitation_code=? and user_id=? and status=?", Code, "", 0).Updates(Data)
	if result.Error != nil {
		return false
	}
	return result.RowsAffected > 0
}

/**
 * 完成邀请码的状态
 */
func FinishInvitationCode(Code string, UserId string) bool {
	Data := make(map[string]interface{}, 0)
	Data["status"] = 2
	result := db.DB.MysqlDB.DefaultGormDB().Table("invitations").Where("invitation_code=? and user_id=? and status=?", Code, UserId, 1).Updates(Data)
	if result.Error != nil {
		return false
	}
	return result.RowsAffected > 0
}

func CreateRandomString(strlen int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < strlen; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
