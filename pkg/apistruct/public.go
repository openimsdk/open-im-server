package apistruct

type ApiUserInfo struct {
	UserID      string `json:"userID" binding:"required,min=1,max=64" swaggo:"true,用户ID,"`
	Nickname    string `json:"nickname" binding:"omitempty,min=1,max=64" swaggo:"true,my id,19"`
	FaceURL     string `json:"faceURL" binding:"omitempty,max=1024"`
	Gender      int32  `json:"gender" binding:"omitempty,oneof=0 1 2"`
	PhoneNumber string `json:"phoneNumber" binding:"omitempty,max=32"`
	Birth       int64  `json:"birth" binding:"omitempty"`
	Email       string `json:"email" binding:"omitempty,max=64"`
	CreateTime  int64  `json:"createTime"`
	Ex          string `json:"ex" binding:"omitempty,max=1024"`
}

type GroupAddMemberInfo struct {
	UserID    string `json:"userID" binding:"required"`
	RoleLevel int32  `json:"roleLevel" binding:"required,oneof= 1 3"`
}
