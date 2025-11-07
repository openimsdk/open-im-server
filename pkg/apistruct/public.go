package apistruct

type GroupAddMemberInfo struct {
	UserID    string `json:"userID"    binding:"required"`
	RoleLevel int32  `json:"roleLevel" binding:"required,oneof= 1 3"`
}
