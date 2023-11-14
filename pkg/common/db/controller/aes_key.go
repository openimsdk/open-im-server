package controller

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/OpenIMSDK/tools/tx"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
	"sort"
	"strings"
)

type AesKeyDatabase interface {
	AcquireAesKey(ctx context.Context, conversationType int32, userId, friendId, groupId string) (key *relation.AesKeyModel, err error)
	AcquireAesKeys(ctx context.Context, userId string) (key []*relation.AesKeyModel, err error)
}

type aesKeyDatabase struct {
	key relation.AesKeyModelInterface
	tx  tx.Tx
}

func newAesKeyDatabase(key relation.AesKeyModelInterface, tx tx.Tx) *aesKeyDatabase {
	return &aesKeyDatabase{key: key, tx: tx}
}

func (a *aesKeyDatabase) AcquireAesKey(ctx context.Context, conversationType int32, userId, friendId, groupId string) (key *relation.AesKeyModel, err error) {
	//TODO implement me
	panic("implement me")
}

func (a *aesKeyDatabase) AcquireAesKeys(ctx context.Context, userId string) (key []*relation.AesKeyModel, err error) {
	//TODO implement me
	panic("implement me")
}

func (a *aesKeyDatabase) generateKeyConversationsID(args ...string) string {
	sort.Strings(args)
	combinedStr := strings.Join(args, "")
	md5Value := md5.Sum([]byte(combinedStr))
	md5Str := fmt.Sprintf("%x", md5Value)
	return md5Str[:16]
}
