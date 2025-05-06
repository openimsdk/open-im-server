package controller

import (
	"context"

	"github.com/golang-jwt/jwt/v4"
	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/tokenverify"
)

type AuthDatabase interface {
	// If the result is empty, no error is returned.
	GetTokensWithoutError(ctx context.Context, userID string, platformID int) (map[string]int, error)

	GetTemporaryTokensWithoutError(ctx context.Context, userID string, platformID int, token string) error
	// Create token
	CreateToken(ctx context.Context, userID string, platformID int) (string, error)

	BatchSetTokenMapByUidPid(ctx context.Context, tokens []string) error

	SetTokenMapByUidPid(ctx context.Context, userID string, platformID int, m map[string]int) error
}

type multiLoginConfig struct {
	Policy       int
	MaxNumOneEnd int
}

type authDatabase struct {
	cache        cache.TokenModel
	accessSecret string
	accessExpire int64
	multiLogin   multiLoginConfig
	adminUserIDs []string
}

func NewAuthDatabase(cache cache.TokenModel, accessSecret string, accessExpire int64, multiLogin config.MultiLogin, adminUserIDs []string) AuthDatabase {
	return &authDatabase{cache: cache, accessSecret: accessSecret, accessExpire: accessExpire, multiLogin: multiLoginConfig{
		Policy:       multiLogin.Policy,
		MaxNumOneEnd: multiLogin.MaxNumOneEnd,
	}, adminUserIDs: adminUserIDs,
	}
}

// If the result is empty.
func (a *authDatabase) GetTokensWithoutError(ctx context.Context, userID string, platformID int) (map[string]int, error) {
	return a.cache.GetTokensWithoutError(ctx, userID, platformID)
}

func (a *authDatabase) GetTemporaryTokensWithoutError(ctx context.Context, userID string, platformID int, token string) error {
	return a.cache.HasTemporaryToken(ctx, userID, platformID, token)
}

func (a *authDatabase) SetTokenMapByUidPid(ctx context.Context, userID string, platformID int, m map[string]int) error {
	return a.cache.SetTokenMapByUidPid(ctx, userID, platformID, m)
}

func (a *authDatabase) BatchSetTokenMapByUidPid(ctx context.Context, tokens []string) error {
	setMap := make(map[string]map[string]any)
	for _, token := range tokens {
		claims, err := tokenverify.GetClaimFromToken(token, authverify.Secret(a.accessSecret))
		if err != nil {
			continue
		}
		key := cachekey.GetTokenKey(claims.UserID, claims.PlatformID)
		if v, ok := setMap[key]; ok {
			v[token] = constant.KickedToken
		} else {
			setMap[key] = map[string]any{
				token: constant.KickedToken,
			}
		}
	}
	if err := a.cache.BatchSetTokenMapByUidPid(ctx, setMap); err != nil {
		return err
	}
	return nil
}

// Create Token.
func (a *authDatabase) CreateToken(ctx context.Context, userID string, platformID int) (string, error) {
	isAdmin := authverify.IsManagerUserID(userID, a.adminUserIDs)
	if !isAdmin {
		tokens, err := a.cache.GetAllTokensWithoutError(ctx, userID)
		if err != nil {
			return "", err
		}

	deleteTokenKey, kickedTokenKey, adminTokens, err := a.checkToken(ctx, tokens, platformID)
	if err != nil {
		return "", err
	}
	if len(deleteTokenKey) != 0 {
		err = a.cache.DeleteTokenByTokenMap(ctx, userID, deleteTokenKey)
		if err != nil {
			return "", err
		}
	}
	if len(kickedTokenKey) != 0 {
		for plt, ks := range kickedTokenKey {
			for _, k := range ks {
				err := a.cache.SetTokenFlagEx(ctx, userID, plt, k, constant.KickedToken)
				if err != nil {
					return "", err
				}
				log.ZDebug(ctx, "kicked token in create token", "token", k)
			}
		}
	}
	if len(adminTokens) != 0 {
		if err = a.cache.DeleteAndSetTemporary(ctx, userID, constant.AdminPlatformID, adminTokens); err != nil {
			return "", err
		}
	}

	claims := tokenverify.BuildClaims(userID, platformID, a.accessExpire)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.accessSecret))
	if err != nil {
		return "", errs.WrapMsg(err, "token.SignedString")
	}

	if !isAdmin {
		if err = a.cache.SetTokenFlagEx(ctx, userID, platformID, tokenString, constant.NormalToken); err != nil {
			return "", err
		}
	}

	return tokenString, nil
}

// checkToken will check token by tokenPolicy and return deleteToken,kickToken,deleteAdminToken
func (a *authDatabase) checkToken(ctx context.Context, tokens map[int]map[string]int, platformID int) (map[int][]string, map[int][]string, []string, error) {
	// todo: Asynchronous deletion of old data.
	var (
		loginTokenMap  = make(map[int][]string) // The length of the value of the map must be greater than 0
		deleteToken    = make(map[int][]string)
		kickToken      = make(map[int][]string)
		adminToken     = make([]string, 0)
		unkickTerminal = ""
	)

	for plfID, tks := range tokens {
		for k, v := range tks {
			_, err := tokenverify.GetClaimFromToken(k, authverify.Secret(a.accessSecret))
			if err != nil || v != constant.NormalToken {
				deleteToken[plfID] = append(deleteToken[plfID], k)
			} else {
				if plfID != constant.AdminPlatformID {
					loginTokenMap[plfID] = append(loginTokenMap[plfID], k)
				} else {
					adminToken = append(adminToken, k)
				}
			}
		}
	}

	switch a.multiLogin.Policy {
	case constant.DefalutNotKick:
		for plt, ts := range loginTokenMap {
			l := len(ts)
			if platformID == plt {
				l++
			}
			limit := a.multiLogin.MaxNumOneEnd
			if l > limit {
				kickToken[plt] = ts[:l-limit]
			}
		}
	case constant.AllLoginButSameTermKick:
		for plt, ts := range loginTokenMap {
			kickToken[plt] = ts[:len(ts)-1]

			if plt == platformID {
				kickToken[plt] = append(kickToken[plt], ts[len(ts)-1])
			}
		}
	case constant.PCAndOther:
		unkickTerminal = constant.TerminalPC
		if constant.PlatformIDToClass(platformID) != unkickTerminal {
			for plt, ts := range loginTokenMap {
				if constant.PlatformIDToClass(plt) != unkickTerminal {
					kickToken[plt] = ts
				}
			}
		} else {
			var (
				preKickToken string
				preKickPlt   int
				reserveToken = false
			)
			for plt, ts := range loginTokenMap {
				if constant.PlatformIDToClass(plt) != unkickTerminal {
					// Keep a token from another end
					if !reserveToken {
						reserveToken = true
						kickToken[plt] = ts[:len(ts)-1]
						preKickToken = ts[len(ts)-1]
						preKickPlt = plt
						continue
					} else {
						// Prioritize keeping Android
						if plt == constant.AndroidPlatformID {
							if preKickToken != "" {
								kickToken[preKickPlt] = append(kickToken[preKickPlt], preKickToken)
							}
							kickToken[plt] = ts[:len(ts)-1]
						} else {
							kickToken[plt] = ts
						}
					}
				}
			}
		}
	case constant.AllLoginButSameClassKick:
		var (
			reserved = make(map[string]struct{})
		)

		for plt, ts := range loginTokenMap {
			if constant.PlatformIDToClass(plt) == constant.PlatformIDToClass(platformID) {
				kickToken[plt] = ts
			} else {
				if _, ok := reserved[constant.PlatformIDToClass(plt)]; !ok {
					reserved[constant.PlatformIDToClass(plt)] = struct{}{}
					kickToken[plt] = ts[:len(ts)-1]
					continue
				} else {
					kickToken[plt] = ts
				}
			}
		}
	default:
		return nil, nil, nil, errs.New("unknown multiLogin policy").Wrap()
	}

	//var adminTokenMaxNum = a.multiLogin.MaxNumOneEnd
	//if a.multiLogin.Policy == constant.Customize {
	//	adminTokenMaxNum = a.multiLogin.CustomizeLoginNum[constant.AdminPlatformID]
	//}
	//l := len(adminToken)
	//if platformID == constant.AdminPlatformID {
	//	l++
	//}
	//if l > adminTokenMaxNum {
	//	kickToken = append(kickToken, adminToken[:l-adminTokenMaxNum]...)
	//}
	var deleteAdminToken []string
	if platformID == constant.AdminPlatformID {
		deleteAdminToken = adminToken
	}
	return deleteToken, kickToken, deleteAdminToken, nil
}
