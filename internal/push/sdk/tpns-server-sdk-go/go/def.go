package tpns

type CommonRspEnv string

const (
	// EnvProd
	EnvProd CommonRspEnv = "product"
	// EnvDev
	EnvDev CommonRspEnv = "dev"
)

type CommonRsp struct {
	// TODO: doc this
	Seq int64 `json:"seq"`

	PushID string `json:"push_id"`

	RetCode int `json:"ret_code"`

	Environment CommonRspEnv `json:"environment"`

	ErrMsg string `json:"err_msg,omitempty"`

	Result map[string]string `json:"result,omitempty"`
}

type AudienceType string

const (
	AdAll AudienceType = "all"

	AdTag AudienceType = "tag"

	AdToken AudienceType = "token"

	AdTokenList AudienceType = "token_list"

	AdAccount AudienceType = "account"

	AdAccountList AudienceType = "account_list"

	AdPackageAccount AudienceType = "package_account_push"

	AdPackageToken AudienceType = "package_token_push"
)

// MessageType push API message_type
type MessageType string

const (
	MsgTypeNotify MessageType = "notify"

	MsgTypeMessage MessageType = "message"
)

type Request struct {
	AudienceType AudienceType `json:"audience_type"`

	Message Message `json:"message"`

	MessageType MessageType `json:"message_type"`

	Tag []TagRule `json:"tag_rules,omitempty"`

	TokenList []string `json:"token_list,omitempty"`

	AccountList []string `json:"account_list,omitempty"`

	Environment CommonRspEnv `json:"environment,omitempty"`

	UploadId int `json:"upload_id,omitempty"`

	ExpireTime int `json:"expire_time,omitempty"`

	SendTime string `json:"send_time,omitempty"`

	MultiPkg bool `json:"multi_pkg,omitempty"`

	PlanId string `json:"plan_id,omitempty"`

	AccountPushType int `json:"account_push_type,omitempty"`

	PushSpeed int `json:"push_speed,omitempty"`

	CollapseId int `json:"collapse_id"`

	TPNSOnlinePushType int `json:"tpns_online_push_type"`

	ChannelRules []*ChannelDistributeRule `json:"channel_rules,omitempty"`

	LoopParam     *PushLoopParam `json:"loop_param,omitempty"`
	ForceCollapse bool           `json:"force_collapse"`
}

type TagListOperation string

type ChannelDistributeRule struct {
	ChannelName string `json:"channel"`
	Disable     bool   `json:"disable"`
}

type PushLoopParam struct {
	StartDate string `json:"startDate"`

	EndDate string `json:"endDate"`

	LoopType PushLoopType `json:"loopType"`

	LoopDayIndexs []uint32 `json:"loopDayIndexs"`

	DayTimes []string `json:"dayTimes"`
}

type PushLoopType int32

const (
	TagListOpAnd TagListOperation = "AND"

	TagListOpOr TagListOperation = "OR"
)

type TagType string

const (
	XGAutoProvince      TagType = "xg_auto_province"
	XGAutoActive        TagType = "xg_auto_active"
	XGUserDefine        TagType = "xg_user_define"
	XGAutoVersion       TagType = "xg_auto_version"
	XGAutoSdkversion    TagType = "xg_auto_sdkversion"
	XGAutoDevicebrand   TagType = "xg_auto_devicebrand"
	XGAutoDeviceversion TagType = "xg_auto_deviceversion"
	XGAutoCountry       TagType = "xg_auto_country"
)

type TagRule struct {
	TagItems []TagItem `json:"tag_items"`

	IsNot bool `json:"is_not"`

	Operator TagListOperation `json:"operator"`
}

type TagItem struct {
	// 标签
	Tags          []string         `json:"tags"`
	IsNot         bool             `json:"is_not"`
	TagsOperator  TagListOperation `json:"tags_operator"`
	ItemsOperator TagListOperation `json:"items_operator"`
	TagType       TagType          `json:"tag_type"`
}

type Message struct {
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`

	AcceptTime []AcceptTimeItem `json:"accept_time,omitempty"`

	Android *AndroidParams `json:"android,omitempty"`

	IOS *IOSParams `json:"ios,omitempty"`

	ThreadId string `json:"thread_id,omitempty"`

	ThreadSumtext string `json:"thread_sumtext,omitempty"`

	XGMediaResources string `json:"xg_media_resources,omitempty"`

	XGMediaAudioResources string `json:"xg_media_audio_resources,omitempty"`
}

type AcceptTimeItem struct {
	Start HourAndMin `json:"start,omitempty"`
	End   HourAndMin `json:"end,omitempty"`
}

type HourAndMin struct {
	Hour string `json:"hour,omitempty"`
	Min  string `json:"min,omitempty"`
}

type AndroidParams struct {
	BuilderId *int `json:"builder_id,omitempty"`

	Ring *int `json:"ring,omitempty"`

	RingRaw string `json:"ring_raw,omitempty"`

	Vibrate *int `json:"vibrate,omitempty"`

	Lights *int `json:"lights,omitempty"`

	Clearable *int `json:"clearable,omitempty"`

	IconType *int `json:"icon_type"`

	IconRes string `json:"icon_res,omitempty"`

	StyleId *int `json:"style_id,omitempty"`

	SmallIcon string `json:"small_icon,omitempty"`

	Action *Action `json:"action,omitempty"`

	CustomContent string `json:"custom_content,omitempty"`

	ShowType *int `json:"show_type,omitempty"`

	NChId string `json:"n_ch_id,omitempty"`

	NChName string `json:"n_ch_name,omitempty"`

	HwChId string `json:"hw_ch_id,omitempty"`

	XmChId string `json:"xm_ch_id,omitempty"`

	OppoChId string `json:"oppo_ch_id,omitempty"`

	VivoChId string `json:"vivo_ch_id,omitempty"`

	BadgeType *int `json:"badge_type,omitempty"`

	IconColor *int `json:"icon_color,omitempty"`
}

type Action struct {
	ActionType *int    `json:"action_type,omitempty"`
	Activity   string  `json:"activity"`
	AtyAttr    AtyAttr `json:"aty_attr,omitempty"`
	Intent     string  `json:"intent"`
	Browser    Browser `json:"browser,omitempty"`
}

type Browser struct {
	Url     string `json:"url,omitempty"`
	Confirm *int   `json:"confirm,omitempty"`
}

type AtyAttr struct {
	AttrIf *int `json:"if,omitempty"`
	Pf     *int `json:"pf,omitempty"`
}

type IOSParams struct {
	Aps *Aps `json:"aps,omitempty"`

	CustomContent string `json:"custom_content,omitempty"`
}

type Aps struct {
	Alert               map[string]string `json:"alert,omitempty"`
	BadgeType           *int              `json:"badge_type,omitempty"`
	Category            string            `json:"category,omitempty"`
	ContentAvailableInt *int              `json:"content-available,omitempty"`
	MutableContent      *int              `json:"mutable-content,omitempty"`
	Sound               string            `json:"sound,omitempty"`
}
