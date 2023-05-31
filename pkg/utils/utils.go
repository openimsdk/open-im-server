package utils

import (
	"hash/crc32"
	"math/rand"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/unrelation"
	sdkws "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

// copy a by b  b->a
func CopyStructFields(a interface{}, b interface{}, fields ...string) (err error) {
	return copier.Copy(a, b)
}

func Wrap1(err error) error {
	return errors.Wrap(err, "==> "+printCallerNameAndLine())
}

func Wrap2[T any](a T, err error) (T, error) {
	if err != nil {
		return a, errors.Wrap(err, "==> "+printCallerNameAndLine())
	}
	return a, nil
}

func Wrap3[T any, V any](a T, b V, err error) (T, V, error) {
	if err != nil {
		return a, b, errors.Wrap(err, "==> "+printCallerNameAndLine())
	}
	return a, b, nil
}

func Wrap(err error, message string) error {
	return errors.Wrap(err, "==> "+printCallerNameAndLine()+message)
}

func WithMessage(err error, message string) error {
	return errors.WithMessage(err, "==> "+printCallerNameAndLine()+message)
}

func printCallerNameAndLine() string {
	pc, _, line, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name() + "()@" + strconv.Itoa(line) + ": "
}

func GetSelfFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return cleanUpFuncName(runtime.FuncForPC(pc).Name())
}

func GetFuncName(skips ...int) string {
	skip := 1
	if len(skips) > 0 {
		skip = skips[0] + 1
	}
	pc, _, _, _ := runtime.Caller(skip)
	return cleanUpFuncName(runtime.FuncForPC(pc).Name())
}

func cleanUpFuncName(funcName string) string {
	end := strings.LastIndex(funcName, ".")
	if end == -1 {
		return ""
	}
	return funcName[end+1:]
}

// Get the intersection of two slices
func Intersect(slice1, slice2 []int64) []int64 {
	m := make(map[int64]bool)
	n := make([]int64, 0)
	for _, v := range slice1 {
		m[v] = true
	}
	for _, v := range slice2 {
		flag, _ := m[v]
		if flag {
			n = append(n, v)
		}
	}
	return n
}

// Get the diff of two slices
func Difference(slice1, slice2 []int64) []int64 {
	m := make(map[int64]bool)
	n := make([]int64, 0)
	inter := Intersect(slice1, slice2)
	for _, v := range inter {
		m[v] = true
	}
	for _, v := range slice1 {
		if !m[v] {
			n = append(n, v)
		}
	}

	for _, v := range slice2 {
		if !m[v] {
			n = append(n, v)
		}
	}
	return n
}

// Get the intersection of two slices
func IntersectString(slice1, slice2 []string) []string {
	m := make(map[string]bool)
	n := make([]string, 0)
	for _, v := range slice1 {
		m[v] = true
	}
	for _, v := range slice2 {
		flag, _ := m[v]
		if flag {
			n = append(n, v)
		}
	}
	return n
}

// Get the diff of two slices
func DifferenceString(slice1, slice2 []string) []string {
	m := make(map[string]bool)
	n := make([]string, 0)
	inter := IntersectString(slice1, slice2)
	for _, v := range inter {
		m[v] = true
	}
	for _, v := range slice1 {
		if !m[v] {
			n = append(n, v)
		}
	}

	for _, v := range slice2 {
		if !m[v] {
			n = append(n, v)
		}
	}
	return n
}

func OperationIDGenerator() string {
	return strconv.FormatInt(time.Now().UnixNano()+int64(rand.Uint32()), 10)
}

func Pb2String(pb proto.Message) (string, error) {
	s, err := proto.Marshal(pb)
	if err != nil {
		return "", err
	}
	return string(s), nil
}

func String2Pb(s string, pb proto.Message) error {
	return proto.Unmarshal([]byte(s), pb)
}

func GetHashCode(s string) uint32 {
	return crc32.ChecksumIEEE([]byte(s))
}

func GetNotificationConversationID(msg *sdkws.MsgData) string {
	switch msg.SessionType {
	case constant.SingleChatType:
		l := []string{msg.SendID, msg.RecvID}
		sort.Strings(l)
		return "n_" + strings.Join(l, "_")
	case constant.GroupChatType:
		return "n_" + msg.GroupID
	case constant.SuperGroupChatType:
		return "n_" + msg.GroupID
	case constant.NotificationChatType:
		return "n_" + msg.SendID + "_" + msg.RecvID
	}
	return ""
}

func GetChatConversationIDByMsg(msg *sdkws.MsgData) string {
	switch msg.SessionType {
	case constant.SingleChatType:
		l := []string{msg.SendID, msg.RecvID}
		sort.Strings(l)
		return "si_" + strings.Join(l, "_")
	case constant.GroupChatType:
		return "g_" + msg.GroupID
	case constant.SuperGroupChatType:
		return "sg_" + msg.GroupID
	case constant.NotificationChatType:
		return "sn_" + msg.SendID + "_" + msg.RecvID
	}
	return ""
}

func GenConversationIDForSingle(sendID, recvID string) string {
	l := []string{sendID, recvID}
	sort.Strings(l)
	return "si_" + strings.Join(l, "_")
}

func GenConversationUniqueKeyForGroup(groupID string) string {
	return groupID
}

func GenGroupConversationID(groupID string) string {
	return "sg_" + groupID
}

func GenConversationUniqueKeyForSingle(sendID, recvID string) string {
	l := []string{sendID, recvID}
	sort.Strings(l)
	return strings.Join(l, "_")
}

func GenConversationUniqueKey(msg *sdkws.MsgData) string {
	switch msg.SessionType {
	case constant.SingleChatType, constant.NotificationChatType:
		l := []string{msg.SendID, msg.RecvID}
		sort.Strings(l)
		return strings.Join(l, "_")
	case constant.SuperGroupChatType:
		return msg.GroupID
	}
	return ""
}

func GetConversationIDByMsgModel(msg *unrelation.MsgDataModel) string {
	options := Options(msg.Options)
	switch msg.SessionType {
	case constant.SingleChatType:
		l := []string{msg.SendID, msg.RecvID}
		sort.Strings(l)
		if !options.IsNotNotification() {
			return "n_" + strings.Join(l, "_")
		}
		return "si_" + strings.Join(l, "_") // single chat
	case constant.GroupChatType:
		if !options.IsNotNotification() {
			return "n_" + msg.GroupID // group chat
		}
		return "g_" + msg.GroupID // group chat
	case constant.SuperGroupChatType:
		if !options.IsNotNotification() {
			return "n_" + msg.GroupID // super group chat
		}
		return "sg_" + msg.GroupID // super group chat
	case constant.NotificationChatType:
		if !options.IsNotNotification() {
			return "n_" + msg.SendID + "_" + msg.RecvID // super group chat
		}
		return "sn_" + msg.SendID + "_" + msg.RecvID // server notification chat
	}
	return ""
}

func GetConversationIDByMsg(msg *sdkws.MsgData) string {
	options := Options(msg.Options)
	switch msg.SessionType {
	case constant.SingleChatType:
		l := []string{msg.SendID, msg.RecvID}
		sort.Strings(l)
		if !options.IsNotNotification() {
			return "n_" + strings.Join(l, "_")
		}
		return "si_" + strings.Join(l, "_") // single chat
	case constant.GroupChatType:
		if !options.IsNotNotification() {
			return "n_" + msg.GroupID // group chat
		}
		return "g_" + msg.GroupID // group chat
	case constant.SuperGroupChatType:
		if !options.IsNotNotification() {
			return "n_" + msg.GroupID // super group chat
		}
		return "sg_" + msg.GroupID // super group chat
	case constant.NotificationChatType:
		if !options.IsNotNotification() {
			return "n_" + msg.SendID + "_" + msg.RecvID // super group chat
		}
		return "sn_" + msg.SendID + "_" + msg.RecvID // server notification chat
	}
	return ""
}

func GetConversationIDBySessionType(sessionType int, ids ...string) string {
	sort.Strings(ids)
	if len(ids) > 2 || len(ids) < 1 {
		return ""
	}
	switch sessionType {
	case constant.SingleChatType:
		return "si_" + strings.Join(ids, "_") // single chat
	case constant.GroupChatType:
		return "g_" + ids[0] // group chat
	case constant.SuperGroupChatType:
		return "sg_" + ids[0] // super group chat
	case constant.NotificationChatType:
		return "sn_" + ids[0] // server notification chat
	}
	return ""
}

func IsNotification(conversationID string) bool {
	return strings.HasPrefix(conversationID, "n_")
}

func IsNotificationByMsg(msg *sdkws.MsgData) bool {
	return !Options(msg.Options).IsNotNotification()
}

func ParseConversationID(msg *sdkws.MsgData) (isNotification bool, conversationID string) {
	options := Options(msg.Options)
	switch msg.SessionType {
	case constant.SingleChatType:
		l := []string{msg.SendID, msg.RecvID}
		sort.Strings(l)
		if !options.IsNotNotification() {
			return true, "n_" + strings.Join(l, "_")
		}
		return false, "si_" + strings.Join(l, "_") // single chat
	case constant.SuperGroupChatType:
		if !options.IsNotNotification() {
			return true, "n_" + msg.GroupID // super group chat
		}
		return false, "sg_" + msg.GroupID // super group chat
	case constant.NotificationChatType:
		if !options.IsNotNotification() {
			return true, "n_" + msg.SendID + "_" + msg.RecvID // super group chat
		}
		return false, "sn_" + msg.SendID + "_" + msg.RecvID // server notification chat
	}
	return false, ""
}

func GetNotificationConversationIDByConversationID(conversationID string) string {
	l := strings.Split(conversationID, "_")
	if len(l) > 1 {
		l[0] = "n"
		return strings.Join(l, "_")
	}
	return ""
}

func GetSeqsBeginEnd(seqs []int64) (int64, int64) {
	if len(seqs) == 0 {
		return 0, 0
	}
	return seqs[0], seqs[len(seqs)-1]
}

type MsgBySeq []*sdkws.MsgData

func (s MsgBySeq) Len() int {
	return len(s)
}

func (s MsgBySeq) Less(i, j int) bool {
	return s[i].Seq < s[j].Seq
}

func (s MsgBySeq) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
