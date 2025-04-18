package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/apistruct"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/protocol/auth"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/sdkws"
	pbuser "github.com/openimsdk/protocol/user"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/system/program"
)

// 测试第二期：
// 自动化测试第二期：需要多个会话 每个会话都会有新消息
// 需求：
// （1）整体系统达到百万个群组。
// （2）每个人同事有 100个10万人群， 1000个999人群。
// （3）每个同事1万好友。
// （4）每个群周期性发送消息。

// 如何测试：

// （1）先让之前的测试继续运行，达到万人就停止，这样就有一个万人群了，同事也有万人好友了。
// （2）创建10万个新用户。
// （3）同事和新用户，建立100个10万人群。每秒钟一个。
// （4）同事和新用户，创建1000个999人群。每秒钟一个。
// （5）等以上完成后，10万人群，每秒钟发送1条消息。999人群，每分钟发送1条消息。

// TODO
// 不要出现第一个邀请的是0，应该是从1开始邀请
// 最好是加一个检测，先判断这一批有没有在里面，有的话剔除
var (
	//  Use default userIDs List for testing, need to be created.
	TestTargetUserList = []string{
		// "<need-update-it>",
		"6760971175",
		"test_v3_u0",
		"test_v3_u1",
		"test_v3_u2",
		"test_v3_u3",
	}
	DefaultGroupID = "<need-update-it>" // Use default group ID for testing, need to be created.
)

var (
	ApiAddress string

	// API method
	GetAdminToken      = "/auth/get_admin_token"
	UserCheck          = "/user/account_check"
	CreateUser         = "/user/user_register"
	ImportFriend       = "/friend/import_friend"
	InviteToGroup      = "/group/invite_user_to_group"
	GetGroupMemberInfo = "/group/get_group_members_info"
	SendMsg            = "/msg/send_msg"
	CreateGroup        = "/group/create_group"
	GetUserToken       = "/auth/user_token"
)

const (
	MaxUser = 1000
	// MaxUser            = 100000
	Max100KGroup       = 100
	Max999Group        = 1000
	MaxInviteUserLimit = 999

	CreateUserTicker  = 1 * time.Second
	CreateGroupTicker = 1 * time.Second
	// Create100KGroupTicker    = 1 * time.Minute
	// Create999GroupTicker     = 1 * time.Minute
	Create100KGroupTicker    = 1 * time.Second
	Create999GroupTicker     = 1 * time.Second
	SendMsgTo100KGroupTicker = 1 * time.Second
	SendMsgTo999GroupTicker  = 1 * time.Minute
)

type BaseResp struct {
	ErrCode int             `json:"errCode"`
	ErrMsg  string          `json:"errMsg"`
	Data    json.RawMessage `json:"data"`
}

type StressTest struct {
	Conf                   *conf
	AdminUserID            string
	AdminToken             string
	DefaultGroupID         string
	DefaultUserID          string
	UserCounter            int
	CreateUserCounter      int
	Create100kGroupCounter int
	Create999GroupCounter  int
	MsgCounter             int
	CreatedUsers           []string
	CreatedGroups          []string
	Mutex                  sync.Mutex
	Ctx                    context.Context
	Cancel                 context.CancelFunc
	HttpClient             *http.Client
	Wg                     sync.WaitGroup
	Once                   sync.Once
}

type conf struct {
	Share config.Share
	Api   config.API
}

func initConfig(configDir string) (*config.Share, *config.API, error) {
	var (
		share     = &config.Share{}
		apiConfig = &config.API{}
	)

	err := config.Load(configDir, config.ShareFileName, config.EnvPrefixMap[config.ShareFileName], share)
	if err != nil {
		return nil, nil, err
	}

	err = config.Load(configDir, config.OpenIMAPICfgFileName, config.EnvPrefixMap[config.OpenIMAPICfgFileName], apiConfig)
	if err != nil {
		return nil, nil, err
	}

	return share, apiConfig, nil
}

// Post Request
func (st *StressTest) PostRequest(ctx context.Context, url string, reqbody any) ([]byte, error) {
	// Marshal body
	jsonBody, err := json.Marshal(reqbody)
	if err != nil {
		log.ZError(ctx, "Failed to marshal request body", err, "url", url, "reqbody", reqbody)
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("operationID", st.AdminUserID)
	if st.AdminToken != "" {
		req.Header.Set("token", st.AdminToken)
	}

	// log.ZInfo(ctx, "Header info is ", "Content-Type", "application/json", "operationID", st.AdminUserID, "token", st.AdminToken)

	resp, err := st.HttpClient.Do(req)
	if err != nil {
		log.ZError(ctx, "Failed to send request", err, "url", url, "reqbody", reqbody)
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.ZError(ctx, "Failed to read response body", err, "url", url)
		return nil, err
	}

	var baseResp BaseResp
	if err := json.Unmarshal(respBody, &baseResp); err != nil {
		log.ZError(ctx, "Failed to unmarshal response body", err, "url", url, "respBody", string(respBody))
		return nil, err
	}

	if baseResp.ErrCode != 0 {
		err = fmt.Errorf(baseResp.ErrMsg)
		log.ZError(ctx, "Failed to send request", err, "url", url, "reqbody", reqbody, "resp", baseResp)
		return nil, err
	}

	return baseResp.Data, nil
}

func (st *StressTest) GetAdminToken(ctx context.Context) (string, error) {
	req := auth.GetAdminTokenReq{
		Secret: st.Conf.Share.Secret,
		UserID: st.AdminUserID,
	}

	resp, err := st.PostRequest(ctx, ApiAddress+GetAdminToken, &req)
	if err != nil {
		return "", err
	}

	data := &auth.GetAdminTokenResp{}
	if err := json.Unmarshal(resp, &data); err != nil {
		return "", err
	}

	return data.Token, nil
}

func (st *StressTest) CheckUser(ctx context.Context, userIDs []string) ([]string, error) {
	req := pbuser.AccountCheckReq{
		CheckUserIDs: userIDs,
	}

	resp, err := st.PostRequest(ctx, ApiAddress+UserCheck, &req)
	if err != nil {
		return nil, err
	}

	data := &pbuser.AccountCheckResp{}
	if err := json.Unmarshal(resp, &data); err != nil {
		return nil, err
	}

	unRegisteredUserIDs := make([]string, 0)

	for _, res := range data.Results {
		if res.AccountStatus == constant.UnRegistered {
			unRegisteredUserIDs = append(unRegisteredUserIDs, res.UserID)
		}
	}

	return unRegisteredUserIDs, nil
}

func (st *StressTest) CreateUser(ctx context.Context, userID string) (string, error) {
	user := &sdkws.UserInfo{
		UserID:   userID,
		Nickname: userID,
	}

	req := pbuser.UserRegisterReq{
		Users: []*sdkws.UserInfo{user},
	}

	_, err := st.PostRequest(ctx, ApiAddress+CreateUser, &req)
	if err != nil {
		return "", err
	}

	st.UserCounter++
	return userID, nil
}

func (st *StressTest) CreateUserBatch(ctx context.Context, userIDs []string) error {
	// The method can import a large number of users at once.
	var userList []*sdkws.UserInfo

	defer st.Once.Do(
		func() {
			st.DefaultUserID = userIDs[0]
			fmt.Println("Default Send User Created ID:", st.DefaultUserID)
		})

	needUserIDs, err := st.CheckUser(ctx, userIDs)
	if err != nil {
		return err
	}

	for _, userID := range needUserIDs {
		user := &sdkws.UserInfo{
			UserID:   userID,
			Nickname: userID,
		}
		userList = append(userList, user)
	}

	req := pbuser.UserRegisterReq{
		Users: userList,
	}

	_, err = st.PostRequest(ctx, ApiAddress+CreateUser, &req)
	if err != nil {
		return err
	}

	st.UserCounter += len(userList)
	return nil
}

// func (st *StressTest) ImportFriend(ctx context.Context, userID string) error {
// 	req := relation.ImportFriendReq{
// 		OwnerUserID:   userID,
// 		FriendUserIDs: TestTargetUserList,
// 	}

// 	_, err := st.PostRequest(ctx, ApiAddress+ImportFriend, &req)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (st *StressTest) GetGroupMembersInfo(ctx context.Context, groupID string, userIDs []string) ([]string, error) {
	req := group.GetGroupMembersInfoReq{
		GroupID: groupID,
		UserIDs: userIDs,
	}

	needInviteUserIDs := make([]string, 0)

	resp, err := st.PostRequest(ctx, ApiAddress+GetGroupMemberInfo, &req)
	if err != nil {
		return nil, err
	}

	data := &group.GetGroupMembersInfoResp{}
	if err := json.Unmarshal(resp, &data); err != nil {
		return nil, err
	}

	return needInviteUserIDs, nil
}

func (st *StressTest) InviteToGroup(ctx context.Context, groupID string, userIDs []string) error {
	req := group.InviteUserToGroupReq{
		GroupID:        groupID,
		InvitedUserIDs: userIDs,
	}
	_, err := st.PostRequest(ctx, ApiAddress+InviteToGroup, &req)
	if err != nil {
		return err
	}

	return nil
}

func (st *StressTest) SendMsg(ctx context.Context, userID string, groupID string) error {
	contentObj := map[string]any{
		"content": fmt.Sprintf("index %d. The current time is %s", st.MsgCounter, time.Now().Format("2006-01-02 15:04:05.000")),
	}

	req := &apistruct.SendMsgReq{
		SendMsg: apistruct.SendMsg{
			SendID:         userID,
			SenderNickname: userID,
			GroupID:        groupID,
			ContentType:    constant.Text,
			SessionType:    constant.ReadGroupChatType,
			Content:        contentObj,
		},
	}

	_, err := st.PostRequest(ctx, ApiAddress+SendMsg, &req)
	if err != nil {
		log.ZError(ctx, "Failed to send message", err, "userID", userID, "req", &req)
		return err
	}

	st.MsgCounter++

	return nil
}

// Max userIDs number is 1000
func (st *StressTest) CreateGroup(ctx context.Context, groupID string, userID string, userIDsList []string) (string, error) {
	// groupID := fmt.Sprintf("StressTestGroup_v2_%d_%s", st.GroupCounter, time.Now().Format("20060102150405"))

	groupInfo := &sdkws.GroupInfo{
		GroupID:   groupID,
		GroupName: groupID,
		GroupType: constant.WorkingGroup,
	}

	req := group.CreateGroupReq{
		OwnerUserID:   userID,
		MemberUserIDs: userIDsList,
		GroupInfo:     groupInfo,
	}

	resp := group.CreateGroupResp{}

	response, err := st.PostRequest(ctx, ApiAddress+CreateGroup, &req)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(response, &resp); err != nil {
		return "", err
	}

	// st.GroupCounter++

	return resp.GroupInfo.GroupID, nil
}

func main() {
	var configPath string
	// defaultConfigDir := filepath.Join("..", "..", "..", "..", "..", "config")
	// flag.StringVar(&configPath, "c", defaultConfigDir, "config path")
	flag.StringVar(&configPath, "c", "", "config path")
	flag.Parse()

	if configPath == "" {
		_, _ = fmt.Fprintln(os.Stderr, "config path is empty")
		os.Exit(1)
		return
	}

	fmt.Printf(" Config Path: %s\n", configPath)

	share, apiConfig, err := initConfig(configPath)
	if err != nil {
		program.ExitWithError(err)
		return
	}

	ApiAddress = fmt.Sprintf("http://%s:%s", "127.0.0.1", fmt.Sprint(apiConfig.Api.Ports[0]))

	ctx, cancel := context.WithCancel(context.Background())
	// ch := make(chan struct{})

	st := &StressTest{
		Conf: &conf{
			Share: *share,
			Api:   *apiConfig,
		},
		AdminUserID: share.IMAdminUserID[0],
		Ctx:         ctx,
		Cancel:      cancel,
		HttpClient: &http.Client{
			Timeout: 50 * time.Second,
		},
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nReceived stop signal, stopping...")

		go func() {
			// time.Sleep(5 * time.Second)
			fmt.Println("Force exit")
			os.Exit(0)
		}()

		st.Cancel()
	}()

	token, err := st.GetAdminToken(st.Ctx)
	if err != nil {
		log.ZError(ctx, "Get Admin Token failed.", err, "AdminUserID", st.AdminUserID)
	}

	st.AdminToken = token
	fmt.Println("Admin Token:", st.AdminToken)
	fmt.Println("ApiAddress:", ApiAddress)

	// 最后一位 后面在拼接
	// userID := fmt.Sprintf("v2_Stresstest_%d", st.UserCounter)

	for i := range MaxUser {
		userID := fmt.Sprintf("v2_StressTest_User_%d", i)
		st.CreatedUsers = append(st.CreatedUsers, userID)
		st.CreateUserCounter++
	}

	// err = st.CreateUserBatch(st.Ctx, st.CreatedUsers)
	// if err != nil {
	// 	log.ZError(ctx, "Create user failed.", err)
	// }

	const batchSize = 1000
	totalUsers := len(st.CreatedUsers)
	successCount := 0

	if st.DefaultUserID == "" && len(st.CreatedUsers) > 0 {
		st.DefaultUserID = st.CreatedUsers[0]
	}

	for i := 0; i < totalUsers; i += batchSize {
		end := min(i+batchSize, totalUsers)

		userBatch := st.CreatedUsers[i:end]
		log.ZInfo(st.Ctx, "创建用户批次", "批次", i/batchSize+1, "数量", len(userBatch))

		err = st.CreateUserBatch(st.Ctx, userBatch)
		if err != nil {
			log.ZError(st.Ctx, "批量创建用户失败", err, "批次", i/batchSize+1)
		} else {
			successCount += len(userBatch)
			log.ZInfo(st.Ctx, "批量创建用户成功", "批次", i/batchSize+1,
				"进度", fmt.Sprintf("%d/%d", successCount, totalUsers))
		}

		// Execute create 100k group
		st.Wg.Add(1)
		go func() {
			create100kGroupTicker := time.NewTicker(Create100KGroupTicker)
			defer create100kGroupTicker.Stop()

			for i := range Max100KGroup {
				select {
				case <-st.Ctx.Done():
					log.ZInfo(st.Ctx, "Stop Create 100K Group")
					return

				case <-create100kGroupTicker.C:
					// Create 100K groups
					st.Wg.Add(1)
					go func(idx int) {
						defer st.Wg.Done()
						defer func() {
							st.Create100kGroupCounter++
						}()

						groupID := fmt.Sprintf("v2_StressTest_Group_100K_%d", idx)

						if _, err = st.CreateGroup(st.Ctx, groupID, st.DefaultUserID, TestTargetUserList); err != nil {
							log.ZError(st.Ctx, "Create group failed.", err)
							// continue
						}

						for i := 0; i < MaxUser/MaxInviteUserLimit; i++ {
							InviteUserIDs := make([]string, 0)

							startIdx := max(i*MaxInviteUserLimit, 1)
							endIdx := min((i+1)*MaxInviteUserLimit, MaxUser)

							for j := startIdx; j < endIdx; j++ {
								userCreatedID := fmt.Sprintf("v2_StressTest_User_%d", j)
								InviteUserIDs = append(InviteUserIDs, userCreatedID)
							}

							if len(InviteUserIDs) == 0 {
								log.ZWarn(st.Ctx, "InviteUserIDs is empty", nil, "groupID", groupID)
								continue
							}

							// Invite To Group
							if err = st.InviteToGroup(st.Ctx, groupID, InviteUserIDs); err != nil {
								log.ZError(st.Ctx, "Invite To Group failed.", err, "UserID", InviteUserIDs)
								continue
								// os.Exit(1)
								// return
							}
						}
					}(i)
				}
			}
		}()

		// create 999 groups
		st.Wg.Add(1)
		go func() {
			create999GroupTicker := time.NewTicker(Create999GroupTicker)
			defer create999GroupTicker.Stop()

			for i := range Max999Group {
				select {
				case <-st.Ctx.Done():
					log.ZInfo(st.Ctx, "Stop Create 999 Group")
					return

				case <-create999GroupTicker.C:
					// Create 999 groups
					st.Wg.Add(1)
					go func(idx int) {
						defer st.Wg.Done()
						defer func() {
							st.Create999GroupCounter++
						}()

						groupID := fmt.Sprintf("v2_StressTest_Group_1K_%d", idx)

						if groupID, err = st.CreateGroup(st.Ctx, groupID, st.DefaultUserID, TestTargetUserList); err != nil {
							log.ZError(st.Ctx, "Create group failed.", err)
							// continue
						}
						for i := 0; i < MaxUser/MaxInviteUserLimit; i++ {
							InviteUserIDs := make([]string, 0)

							startIdx := i * MaxInviteUserLimit
							endIdx := min((i+1)*MaxInviteUserLimit, MaxUser)

							for j := startIdx; j < endIdx; j++ {
								userCreatedID := fmt.Sprintf("v2_StressTest_User_%d", j)
								InviteUserIDs = append(InviteUserIDs, userCreatedID)
							}

							if len(InviteUserIDs) == 0 {
								log.ZWarn(st.Ctx, "InviteUserIDs is empty", nil, "groupID", groupID)
								continue
							}

							// Invite To Group
							if err = st.InviteToGroup(st.Ctx, groupID, InviteUserIDs); err != nil {
								log.ZError(st.Ctx, "Invite To Group failed.", err, "UserID", InviteUserIDs)
								continue
								// os.Exit(1)
								// return
							}
						}
					}(i)
				}
			}
		}()

		// Send message to 100K groups
		st.Wg.Wait()
		fmt.Println("All groups created successfully, starting to send messages...")

		var groups100K []string
		var groups999 []string

		for i := range Max100KGroup {
			groupID := fmt.Sprintf("v2_StressTest_Group_100K_%d", i)
			groups100K = append(groups100K, groupID)
		}

		for i := range Max999Group {
			groupID := fmt.Sprintf("v2_StressTest_Group_1K_%d", i)
			groups999 = append(groups999, groupID)
		}

		send100kGroupLimiter := make(chan struct{}, 20)
		send999GroupLimiter := make(chan struct{}, 100)

		// execute Send message to 100K groups
		go func() {
			ticker := time.NewTicker(SendMsgTo100KGroupTicker)
			defer ticker.Stop()

			for {
				select {
				case <-st.Ctx.Done():
					log.ZInfo(st.Ctx, "Stop Send Message to 100K Group")
					return

				case <-ticker.C:
					// Send message to 100K groups
					for _, groupID := range groups100K {
						send100kGroupLimiter <- struct{}{}
						go func(groupID string) {
							defer func() { <-send100kGroupLimiter }()
							if err := st.SendMsg(st.Ctx, st.AdminUserID, groupID); err != nil {
								log.ZError(st.Ctx, "Send message to 100K group failed.", err)
							}
						}(groupID)
					}
					log.ZInfo(st.Ctx, "Send message to 100K groups successfully.")
				}
			}
		}()

		// execute Send message to 999 groups
		go func() {
			ticker := time.NewTicker(SendMsgTo999GroupTicker)
			defer ticker.Stop()

			for {
				select {
				case <-st.Ctx.Done():
					log.ZInfo(st.Ctx, "Stop Send Message to 999 Group")
					return

				case <-ticker.C:
					// Send message to 999 groups
					for _, groupID := range groups999 {
						send999GroupLimiter <- struct{}{}
						go func(groupID string) {
							defer func() { <-send999GroupLimiter }()

							if err := st.SendMsg(st.Ctx, st.AdminUserID, groupID); err != nil {
								log.ZError(st.Ctx, "Send message to 999 group failed.", err)
							}
						}(groupID)
					}
					log.ZInfo(st.Ctx, "Send message to 999 groups successfully.")
				}
			}
		}()

	}
	<-st.Ctx.Done()
	fmt.Println("Received signal to exit, shutting down...")
}
