package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"errors"
	"github.com/jinzhu/gorm"
	"math/rand"
	"time"
)

/**
 * 批量生成邀请码
 */
func BatchCreateInvitationCodes(CodeNums int, CodeLen int) ([]string, error) {
	i := CodeNums
	var codes []string
	for {
		if i == 0 {
			break
		}
		code := CreateRandomString(CodeLen)
		invitation := new(db.Invitation)
		invitation.CreateTime = time.Now()
		invitation.InvitationCode = code
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
		codes = append(codes, code)
	}
	return codes, nil
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
func TryLockInvitationCode(Code string, UserID string) bool {
	Data := make(map[string]interface{}, 0)
	Data["user_id"] = UserID
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

func GetInvitationCode(code string) (*db.Invitation, error) {
	invitation := &db.Invitation{
		InvitationCode: code,
	}
	err := db.DB.MysqlDB.DefaultGormDB().Model(invitation).Find(invitation).Error
	if gorm.IsRecordNotFoundError(err) {
		return invitation, nil
	}
	return invitation, err
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

func GetInvitationCodes(pageNumber, showNumber, status int32) ([]db.Invitation, error) {
	var invitationList []db.Invitation
	err := db.DB.MysqlDB.DefaultGormDB().Model(db.Invitation{}).Limit(int(showNumber)).Offset(int(showNumber*(pageNumber-1))).Where("status=?", status).
		Order("create_time desc").Find(&invitationList).Error
	return invitationList, err
}
