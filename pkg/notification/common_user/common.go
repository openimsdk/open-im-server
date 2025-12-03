package common_user

type CommonUser interface {
	GetNickname() string
	GetFaceURL() string
	GetUserID() string
	GetEx() string
}

type CommonGroup interface {
	GetNickname() string
	GetFaceURL() string
	GetGroupID() string
	GetEx() string
}
