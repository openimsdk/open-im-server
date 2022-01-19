package base_info

//UserID               string   `protobuf:"bytes,1,opt,name=UserID" json:"UserID,omitempty"`
//	Nickname             string   `protobuf:"bytes,2,opt,name=Nickname" json:"Nickname,omitempty"`
//	FaceUrl              string   `protobuf:"bytes,3,opt,name=FaceUrl" json:"FaceUrl,omitempty"`
//	Gender               int32    `protobuf:"varint,4,opt,name=Gender" json:"Gender,omitempty"`
//	PhoneNumber          string   `protobuf:"bytes,5,opt,name=PhoneNumber" json:"PhoneNumber,omitempty"`
//	Birth                string   `protobuf:"bytes,6,opt,name=Birth" json:"Birth,omitempty"`
//	Email                string   `protobuf:"bytes,7,opt,name=Email" json:"Email,omitempty"`
//	Ex                   string   `protobuf:"bytes,8,opt,name=Ex" json:"Ex,omitempty"`

type UserRegisterReq struct {
	Secret   string `json:"secret" binding:"required,max=32"`
	Platform int32  `json:"platform" binding:"required,min=1,max=7"`
	ApiUserInfo
	OperationID string `json:"operationID" binding:"required"`
}

type UserTokenInfo struct {
	UserID      string `json:"userID"`
	Token       string `json:"token"`
	ExpiredTime int64  `json:"expiredTime"`
}
type UserRegisterResp struct {
	CommResp
	UserToken UserTokenInfo `json:"data"`
}

type UserTokenReq struct {
	Secret      string `json:"secret" binding:"required,max=32"`
	Platform    int32  `json:"platform" binding:"required,min=1,max=8"`
	UserID      string `json:"userID" binding:"required,min=1,max=64"`
	OperationID string `json:"operationID" binding:"required"`
}

type UserTokenResp struct {
	CommResp
	UserToken UserTokenInfo `json:"data"`
}
