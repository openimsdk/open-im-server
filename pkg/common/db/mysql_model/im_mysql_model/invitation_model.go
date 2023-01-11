package im_mysql_model

import (
	"errors"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

var InvitationDB *gorm.DB

type Invitation struct {
	InvitationCode string    `gorm:"column:invitation_code;primary_key;type:varchar(32)"`
	CreateTime     time.Time `gorm:"column:create_time"`
	UserID         string    `gorm:"column:user_id;index:userID"`
	LastTime       time.Time `gorm:"column:last_time"`
	Status         int32     `gorm:"column:status"`
}

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
		invitation := new(Invitation)
		invitation.CreateTime = time.Now()
		invitation.InvitationCode = code
		invitation.LastTime = time.Now()
		invitation.Status = 0
		invitation.UserID = ""
		result := InvitationDB.Table("invitations").Create(&invitation)
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
	var invitationCode Invitation
	err := InvitationDB.Table("invitations").Where("invitation_code=?", code).Take(&invitationCode).Error
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
	result := InvitationDB.Table("invitations").Where("invitation_code=? and user_id=? and status=?", Code, "", 0).Updates(Data)
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
	result := InvitationDB.Table("invitations").Where("invitation_code=? and user_id=? and status=?", Code, UserId, 1).Updates(Data)
	if result.Error != nil {
		return false
	}
	return result.RowsAffected > 0
}

func GetInvitationCode(code string) (*Invitation, error) {
	invitation := &Invitation{
		InvitationCode: code,
	}
	err := InvitationDB.Model(invitation).Find(invitation).Error
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

func GetInvitationCodes(showNumber, pageNumber, status int32) ([]Invitation, int64, error) {
	var invitationList []Invitation
	db := InvitationDB.Model(Invitation{}).Where("status=?", status)
	var count int64
	err := db.Count(&count).Error
	err = db.Limit(int(showNumber)).Offset(int(showNumber * (pageNumber - 1))).
		Order("create_time desc").Find(&invitationList).Error
	return invitationList, count, err
}
