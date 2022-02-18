package base_info

type MinioStorageCredentialReq struct {
	Action string `form:"Action";binding:"required"`
	DurationSeconds int `form:"DurationSeconds"`
	Version string `form:"Version"`
	Policy string
}

type MiniostorageCredentialResp struct {

}
