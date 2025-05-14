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
	"github.com/openimsdk/protocol/relation"
	"github.com/openimsdk/protocol/sdkws"
	pbuser "github.com/openimsdk/protocol/user"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/system/program"
)

/*
 1. Create one user every minute
 2. Import target users as friends
 3. Add users to the default group
 4. Send a message to the default group every second, containing index and current timestamp
 5. Create a new group every minute and invite target users to join
*/

// !!! ATTENTION: This variable is must be added!
var (
	//  Use default userIDs List for testing, need to be created.
	TestTargetUserList = []string{
		"<need-update-it>",
	}
	DefaultGroupID = "<need-update-it>" // Use default group ID for testing, need to be created.
)

var (
	ApiAddress string

	// API method
	GetAdminToken = "/auth/get_admin_token"
	CreateUser    = "/user/user_register"
	ImportFriend  = "/friend/import_friend"
	InviteToGroup = "/group/invite_user_to_group"
	SendMsg       = "/msg/send_msg"
	CreateGroup   = "/group/create_group"
	GetUserToken  = "/auth/user_token"
)

const (
	MaxUser  = 10000
	MaxGroup = 1000

	CreateUserTicker  = 1 * time.Minute // Ticker is 1min in create user
	SendMessageTicker = 1 * time.Second // Ticker is 1s in send message
	CreateGroupTicker = 1 * time.Minute
)

type BaseResp struct {
	ErrCode int             `json:"errCode"`
	ErrMsg  string          `json:"errMsg"`
	Data    json.RawMessage `json:"data"`
}

type StressTest struct {
	Conf              *conf
	AdminUserID       string
	AdminToken        string
	DefaultGroupID    string
	DefaultSendUserID string
	UserCounter       int
	GroupCounter      int
	MsgCounter        int
	CreatedUsers      []string
	CreatedGroups     []string
	Mutex             sync.Mutex
	Ctx               context.Context
	Cancel            context.CancelFunc
	HttpClient        *http.Client
	Wg                sync.WaitGroup
	Once              sync.Once
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

func (st *StressTest) ImportFriend(ctx context.Context, userID string) error {
	req := relation.ImportFriendReq{
		OwnerUserID:   userID,
		FriendUserIDs: TestTargetUserList,
	}

	_, err := st.PostRequest(ctx, ApiAddress+ImportFriend, &req)
	if err != nil {
		return err
	}

	return nil
}

func (st *StressTest) InviteToGroup(ctx context.Context, userID string) error {
	req := group.InviteUserToGroupReq{
		GroupID:        st.DefaultGroupID,
		InvitedUserIDs: []string{userID},
	}
	_, err := st.PostRequest(ctx, ApiAddress+InviteToGroup, &req)
	if err != nil {
		return err
	}

	return nil
}

func (st *StressTest) SendMsg(ctx context.Context, userID string) error {
	contentObj := map[string]any{
		"content": fmt.Sprintf("index %d. The current time is %s", st.MsgCounter, time.Now().Format("2006-01-02 15:04:05.000")),
	}

	req := map[string]any{
		"sendID":      userID,
		"groupID":     st.DefaultGroupID,
		"contentType": constant.Text,
		"sessionType": constant.ReadGroupChatType,
		"content":     contentObj,
	}

	_, err := st.PostRequest(ctx, ApiAddress+SendMsg, &req)
	if err != nil {
		log.ZError(ctx, "Failed to send message", err, "userID", userID, "req", &req)
		return err
	}

	st.MsgCounter++

	return nil
}

func (st *StressTest) CreateGroup(ctx context.Context, userID string) (string, error) {
	groupID := fmt.Sprintf("StressTestGroup_%d_%s", st.GroupCounter, time.Now().Format("20060102150405"))

	req := map[string]any{
		"memberUserIDs": TestTargetUserList,
		"ownerUserID":   userID,
		"groupInfo": map[string]any{
			"groupID":   groupID,
			"groupName": groupID,
			"groupType": constant.WorkingGroup,
		},
	}
	resp := group.CreateGroupResp{}

	response, err := st.PostRequest(ctx, ApiAddress+CreateGroup, &req)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(response, &resp); err != nil {
		return "", err
	}

	st.GroupCounter++

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
	ch := make(chan struct{})

	defer cancel()

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

		select {
		case <-ch:
		default:
			close(ch)
		}

		st.Cancel()
	}()

	token, err := st.GetAdminToken(st.Ctx)
	if err != nil {
		log.ZError(ctx, "Get Admin Token failed.", err, "AdminUserID", st.AdminUserID)
	}

	st.AdminToken = token
	fmt.Println("Admin Token:", st.AdminToken)
	fmt.Println("ApiAddress:", ApiAddress)

	st.DefaultGroupID = DefaultGroupID

	st.Wg.Add(1)
	go func() {
		defer st.Wg.Done()

		ticker := time.NewTicker(CreateUserTicker)
		defer ticker.Stop()

		for st.UserCounter < MaxUser {
			select {
			case <-st.Ctx.Done():
				log.ZInfo(st.Ctx, "Stop Create user", "reason", "context done")
				return

			case <-ticker.C:
				// Create User
				userID := fmt.Sprintf("%d_Stresstest_%s", st.UserCounter, time.Now().Format("0102150405"))

				userCreatedID, err := st.CreateUser(st.Ctx, userID)
				if err != nil {
					log.ZError(st.Ctx, "Create User failed.", err, "UserID", userID)
					os.Exit(1)
					return
				}
				// fmt.Println("User Created ID:", userCreatedID)

				// Import Friend
				if err = st.ImportFriend(st.Ctx, userCreatedID); err != nil {
					log.ZError(st.Ctx, "Import Friend failed.", err, "UserID", userCreatedID)
					os.Exit(1)
					return
				}

				// Invite To Group
				if err = st.InviteToGroup(st.Ctx, userCreatedID); err != nil {
					log.ZError(st.Ctx, "Invite To Group failed.", err, "UserID", userCreatedID)
					os.Exit(1)
					return
				}

				st.Once.Do(func() {
					st.DefaultSendUserID = userCreatedID
					fmt.Println("Default Send User Created ID:", userCreatedID)
					close(ch)
				})
			}
		}
	}()

	st.Wg.Add(1)
	go func() {
		defer st.Wg.Done()

		ticker := time.NewTicker(SendMessageTicker)
		defer ticker.Stop()
		<-ch

		for {
			select {
			case <-st.Ctx.Done():
				log.ZInfo(st.Ctx, "Stop Send message", "reason", "context done")
				return

			case <-ticker.C:
				// Send Message
				if err = st.SendMsg(st.Ctx, st.DefaultSendUserID); err != nil {
					log.ZError(st.Ctx, "Send Message failed.", err, "UserID", st.DefaultSendUserID)
					continue
				}
			}
		}
	}()

	st.Wg.Add(1)
	go func() {
		defer st.Wg.Done()

		ticker := time.NewTicker(CreateGroupTicker)
		defer ticker.Stop()
		<-ch

		for st.GroupCounter < MaxGroup {

			select {
			case <-st.Ctx.Done():
				log.ZInfo(st.Ctx, "Stop Create Group", "reason", "context done")
				return

			case <-ticker.C:

				// Create Group
				_, err := st.CreateGroup(st.Ctx, st.DefaultSendUserID)
				if err != nil {
					log.ZError(st.Ctx, "Create Group failed.", err, "UserID", st.DefaultSendUserID)
					os.Exit(1)
					return
				}

				// fmt.Println("Group Created ID:", groupID)
			}
		}
	}()

	st.Wg.Wait()
}
