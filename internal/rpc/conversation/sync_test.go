package conversation

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/open-im-server/v3/pkg/util/hashutil"
	pbconversation "github.com/openimsdk/protocol/conversation"
	pbmsg "github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/db/pagination"
	"github.com/openimsdk/tools/mcontext"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func testCtx(userID string) context.Context {
	return mcontext.WithOpUserIDContext(context.Background(), userID)
}

func TestGetFullOwnerConversationIDsFreshDeviceWithoutPaginationKeepsLegacyFullIDs(t *testing.T) {
	withReadInactiveConversationFilterEnabled(t, false)

	const conversationCount = 50000
	ids := newConversationIDs(conversationCount)

	srv := &conversationServer{
		conversationDatabase: &fakeConversationDatabase{conversationIDs: ids},
	}

	resp, err := srv.GetFullOwnerConversationIDs(testCtx("customer-service"), &pbconversation.GetFullOwnerConversationIDsReq{
		UserID: "customer-service",
		IdHash: 0,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Equal {
		t.Fatal("fresh device idHash=0 should not be equal to the server-side conversation ID hash")
	}
	if got := len(resp.ConversationIDs); got != conversationCount {
		t.Fatalf("expected fresh device sync to return all %d conversation IDs, got %d", conversationCount, got)
	}
	if resp.Total != conversationCount {
		t.Fatalf("expected total %d, got %d", conversationCount, resp.Total)
	}
}

func TestGetFullOwnerConversationIDsFreshDeviceWithPaginationReturnsPage(t *testing.T) {
	withReadInactiveConversationFilterEnabled(t, false)

	const conversationCount = 50000
	ids := newConversationIDs(conversationCount)

	srv := &conversationServer{
		conversationDatabase: &fakeConversationDatabase{conversationIDs: ids},
	}

	resp, err := srv.GetFullOwnerConversationIDs(testCtx("customer-service"), &pbconversation.GetFullOwnerConversationIDsReq{
		UserID: "customer-service",
		IdHash: 0,
		Pagination: &sdkws.RequestPagination{
			PageNumber: 2,
			ShowNumber: 100,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Equal {
		t.Fatal("fresh device idHash=0 should not be equal to the server-side conversation ID hash")
	}
	if got := len(resp.ConversationIDs); got != 100 {
		t.Fatalf("expected paged sync to return 100 conversation IDs, got %d", got)
	}
	if resp.Total != conversationCount {
		t.Fatalf("expected total %d, got %d", conversationCount, resp.Total)
	}
	if resp.ConversationIDs[0] != ids[100] {
		t.Fatalf("expected first ID on page 2 to be %q, got %q", ids[100], resp.ConversationIDs[0])
	}
}

func TestGetFullOwnerConversationIDsInvalidPaginationKeepsLegacyFullIDs(t *testing.T) {
	withReadInactiveConversationFilterEnabled(t, false)

	const conversationCount = 50000
	ids := newConversationIDs(conversationCount)

	srv := &conversationServer{
		conversationDatabase: &fakeConversationDatabase{conversationIDs: ids},
	}

	resp, err := srv.GetFullOwnerConversationIDs(testCtx("customer-service"), &pbconversation.GetFullOwnerConversationIDsReq{
		UserID: "customer-service",
		IdHash: 0,
		Pagination: &sdkws.RequestPagination{
			PageNumber: 0,
			ShowNumber: 100,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := len(resp.ConversationIDs); got != conversationCount {
		t.Fatalf("expected invalid pagination to keep legacy full IDs, got %d", got)
	}
}

func TestGetFullOwnerConversationIDsMatchingHashReturnsNoIDs(t *testing.T) {
	ids := []string{"si_user_1", "si_user_2"}
	srv := &conversationServer{
		conversationDatabase: &fakeConversationDatabase{conversationIDs: ids},
	}

	resp, err := srv.GetFullOwnerConversationIDs(testCtx("customer-service"), &pbconversation.GetFullOwnerConversationIDsReq{
		UserID: "customer-service",
		IdHash: hashutil.IdHash(ids),
	})
	if err != nil {
		t.Fatal(err)
	}
	if !resp.Equal {
		t.Fatal("matching idHash should be reported as equal")
	}
	if len(resp.ConversationIDs) != 0 {
		t.Fatalf("expected matching idHash to omit conversation IDs, got %d", len(resp.ConversationIDs))
	}
	if resp.Total != int64(len(ids)) {
		t.Fatalf("expected total %d, got %d", len(ids), resp.Total)
	}
}

func TestGetFullOwnerConversationIDsMatchingHashWithPaginationReturnsNoIDs(t *testing.T) {
	ids := []string{"si_user_1", "si_user_2"}
	srv := &conversationServer{
		conversationDatabase: &fakeConversationDatabase{conversationIDs: ids},
	}

	resp, err := srv.GetFullOwnerConversationIDs(testCtx("customer-service"), &pbconversation.GetFullOwnerConversationIDsReq{
		UserID: "customer-service",
		IdHash: hashutil.IdHash(ids),
		Pagination: &sdkws.RequestPagination{
			PageNumber: 1,
			ShowNumber: 1,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !resp.Equal {
		t.Fatal("matching idHash should be reported as equal")
	}
	if len(resp.ConversationIDs) != 0 {
		t.Fatalf("expected matching idHash to omit conversation IDs, got %d", len(resp.ConversationIDs))
	}
	if resp.Total != int64(len(ids)) {
		t.Fatalf("expected total %d, got %d", len(ids), resp.Total)
	}
}

func TestGetFullOwnerConversationIDsReadInactiveKeepsLegacyWhenFlagOff(t *testing.T) {
	withReadInactiveConversationFilterEnabled(t, false)

	ids := []string{"si_user_1", "si_user_2", "si_user_3"}
	srv := &conversationServer{
		conversationDatabase: &fakeConversationDatabase{conversationIDs: ids},
		msgClient: fakeMessageClient{seqs: map[string]*pbmsg.FullSyncSeqs{
			"si_user_2": readInactiveSeq(time.Now().Add(-2 * time.Hour)),
		}},
	}

	resp, err := srv.GetFullOwnerConversationIDs(testCtx("customer-service"), &pbconversation.GetFullOwnerConversationIDsReq{
		UserID: "customer-service",
		IdHash: 0,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := len(resp.ConversationIDs); got != len(ids) {
		t.Fatalf("expected legacy sync to return all %d conversation IDs, got %d", len(ids), got)
	}
	if resp.Total != int64(len(ids)) {
		t.Fatalf("expected total %d, got %d", len(ids), resp.Total)
	}
}

func TestGetFullOwnerConversationIDsReadInactiveFiltersReadInactiveConversations(t *testing.T) {
	withReadInactiveConversationCountThreshold(t, 0)
	withReadInactiveConversationDuration(t, int64(time.Hour/time.Millisecond))

	ids := []string{"si_user_1", "si_user_2", "si_user_3", "si_user_4", "si_user_5"}
	filteredIDs := []string{"si_user_1", "si_user_3"}
	now := time.Now()
	srv := &conversationServer{
		conversationDatabase: &fakeConversationDatabase{conversationIDs: ids},
		msgClient: fakeMessageClient{seqs: map[string]*pbmsg.FullSyncSeqs{
			"si_user_1": unreadSeq(now.Add(-2 * time.Hour)),
			"si_user_2": readInactiveSeq(now.Add(-2 * time.Hour)),
			"si_user_3": readActiveSeq(now.Add(-30 * time.Minute)),
			"si_user_4": readInactiveSeq(now.Add(-3 * time.Hour)),
			"si_user_5": clearedSeq(now.Add(-2 * time.Hour)),
		}},
	}

	resp, err := srv.GetFullOwnerConversationIDs(testCtx("customer-service"), &pbconversation.GetFullOwnerConversationIDsReq{
		UserID: "customer-service",
		IdHash: 0,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !sameStrings(resp.ConversationIDs, filteredIDs) {
		t.Fatalf("expected filtered conversation IDs %v, got %v", filteredIDs, resp.ConversationIDs)
	}
	if resp.Total != int64(len(filteredIDs)) {
		t.Fatalf("expected filtered total %d, got %d", len(filteredIDs), resp.Total)
	}
}

func TestGetFullOwnerConversationIDsReadInactiveFiltersEmptyConversations(t *testing.T) {
	withReadInactiveConversationCountThreshold(t, 0)
	withReadInactiveConversationDuration(t, int64(time.Hour/time.Millisecond))

	ids := []string{"si_user_1", "si_user_2", "si_user_3"}
	now := time.Now()
	srv := &conversationServer{
		conversationDatabase: &fakeConversationDatabase{conversationIDs: ids},
		msgClient: fakeMessageClient{seqs: map[string]*pbmsg.FullSyncSeqs{
			"si_user_1": {HasReadSeq: 0, MaxSeq: 0, MaxSeqTime: 0},
			"si_user_2": readInactiveSeq(now.Add(-2 * time.Hour)),
			"si_user_3": nil,
		}},
	}

	resp, err := srv.GetFullOwnerConversationIDs(testCtx("customer-service"), &pbconversation.GetFullOwnerConversationIDsReq{
		UserID: "customer-service",
		IdHash: 0,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !sameStrings(resp.ConversationIDs, []string{"si_user_3"}) {
		t.Fatalf("expected empty conversations to be filtered and missing seqs to be kept, got %v", resp.ConversationIDs)
	}
	if resp.Total != 1 {
		t.Fatalf("expected filtered total 1, got %d", resp.Total)
	}
}

func TestGetFullOwnerConversationIDsReadInactiveMatchingFilteredHashReturnsNoIDs(t *testing.T) {
	withReadInactiveConversationCountThreshold(t, 0)
	withReadInactiveConversationDuration(t, int64(time.Hour/time.Millisecond))

	ids := []string{"si_user_1", "si_user_2", "si_user_3"}
	filteredIDs := []string{"si_user_1", "si_user_3"}
	now := time.Now()
	srv := &conversationServer{
		conversationDatabase: &fakeConversationDatabase{conversationIDs: ids},
		msgClient: fakeMessageClient{seqs: map[string]*pbmsg.FullSyncSeqs{
			"si_user_1": unreadSeq(now.Add(-2 * time.Hour)),
			"si_user_2": readInactiveSeq(now.Add(-2 * time.Hour)),
			"si_user_3": readActiveSeq(now.Add(-30 * time.Minute)),
		}},
	}

	resp, err := srv.GetFullOwnerConversationIDs(testCtx("customer-service"), &pbconversation.GetFullOwnerConversationIDsReq{
		UserID: "customer-service",
		IdHash: hashutil.IdHash(filteredIDs),
	})
	if err != nil {
		t.Fatal(err)
	}
	if !resp.Equal {
		t.Fatal("matching filtered idHash should be reported as equal")
	}
	if len(resp.ConversationIDs) != 0 {
		t.Fatalf("expected matching filtered idHash to omit conversation IDs, got %d", len(resp.ConversationIDs))
	}
	if resp.Total != int64(len(filteredIDs)) {
		t.Fatalf("expected filtered total %d, got %d", len(filteredIDs), resp.Total)
	}
}

func TestGetFullOwnerConversationIDsReadInactiveKeepsConversationsWithUnknownMaxSeqTime(t *testing.T) {
	withReadInactiveConversationCountThreshold(t, 0)
	withReadInactiveConversationDuration(t, int64(time.Hour/time.Millisecond))

	ids := []string{"si_user_1", "si_user_2"}
	srv := &conversationServer{
		conversationDatabase: &fakeConversationDatabase{conversationIDs: ids},
		msgClient: fakeMessageClient{seqs: map[string]*pbmsg.FullSyncSeqs{
			"si_user_1": {HasReadSeq: 10, MaxSeq: 10, MaxSeqTime: 0},
			"si_user_2": readInactiveSeq(time.Now().Add(-2 * time.Hour)),
		}},
	}

	resp, err := srv.GetFullOwnerConversationIDs(testCtx("customer-service"), &pbconversation.GetFullOwnerConversationIDsReq{
		UserID: "customer-service",
		IdHash: 0,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !sameStrings(resp.ConversationIDs, []string{"si_user_1"}) {
		t.Fatalf("expected unknown maxSeqTime conversation to be kept, got %v", resp.ConversationIDs)
	}
}

func TestGetFullOwnerConversationIDsReadInactiveKeepsPinnedConversations(t *testing.T) {
	withReadInactiveConversationCountThreshold(t, 0)
	withReadInactiveConversationDuration(t, int64(time.Hour/time.Millisecond))

	ids := []string{"si_user_1", "si_user_2", "si_user_3"}
	srv := &conversationServer{
		conversationDatabase: &fakeConversationDatabase{
			conversationIDs:       ids,
			pinnedConversationIDs: []string{"si_user_1", "si_user_2"},
		},
		msgClient: fakeMessageClient{seqs: map[string]*pbmsg.FullSyncSeqs{
			"si_user_1": {HasReadSeq: 0, MaxSeq: 0, MaxSeqTime: 0},
			"si_user_2": readInactiveSeq(time.Now().Add(-2 * time.Hour)),
			"si_user_3": readInactiveSeq(time.Now().Add(-2 * time.Hour)),
		}},
	}

	resp, err := srv.GetFullOwnerConversationIDs(testCtx("customer-service"), &pbconversation.GetFullOwnerConversationIDsReq{
		UserID: "customer-service",
		IdHash: 0,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !sameStrings(resp.ConversationIDs, []string{"si_user_1", "si_user_2"}) {
		t.Fatalf("expected pinned conversations to be kept, got %v", resp.ConversationIDs)
	}
}

func TestGetFullOwnerConversationIDsReadInactivePaginatesFilteredIDs(t *testing.T) {
	withReadInactiveConversationCountThreshold(t, 0)
	withReadInactiveConversationDuration(t, int64(time.Hour/time.Millisecond))

	ids := []string{"si_user_1", "si_user_2", "si_user_3", "si_user_4", "si_user_5"}
	filteredIDs := []string{"si_user_1", "si_user_3", "si_user_5"}
	now := time.Now()
	srv := &conversationServer{
		conversationDatabase: &fakeConversationDatabase{conversationIDs: ids},
		msgClient: fakeMessageClient{seqs: map[string]*pbmsg.FullSyncSeqs{
			"si_user_1": unreadSeq(now.Add(-2 * time.Hour)),
			"si_user_2": readInactiveSeq(now.Add(-2 * time.Hour)),
			"si_user_3": readActiveSeq(now.Add(-30 * time.Minute)),
			"si_user_4": readInactiveSeq(now.Add(-3 * time.Hour)),
			"si_user_5": unreadSeq(now.Add(-3 * time.Hour)),
		}},
	}

	resp, err := srv.GetFullOwnerConversationIDs(testCtx("customer-service"), &pbconversation.GetFullOwnerConversationIDsReq{
		UserID: "customer-service",
		IdHash: 0,
		Pagination: &sdkws.RequestPagination{
			PageNumber: 2,
			ShowNumber: 1,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !sameStrings(resp.ConversationIDs, []string{filteredIDs[1]}) {
		t.Fatalf("expected second filtered page %v, got %v", []string{filteredIDs[1]}, resp.ConversationIDs)
	}
	if resp.Total != int64(len(filteredIDs)) {
		t.Fatalf("expected filtered total %d, got %d", len(filteredIDs), resp.Total)
	}
}

func BenchmarkGetFullOwnerConversationIDsLegacyFullIDs(b *testing.B) {
	readInactiveConversationFilterEnabled = false
	b.Cleanup(func() {
		readInactiveConversationFilterEnabled = true
	})

	ids := newConversationIDs(50000)
	srv := &conversationServer{
		conversationDatabase: &fakeConversationDatabase{conversationIDs: ids},
	}
	req := &pbconversation.GetFullOwnerConversationIDsReq{
		UserID: "customer-service",
		IdHash: 0,
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := srv.GetFullOwnerConversationIDs(testCtx("customer-service"), req)
		if err != nil {
			b.Fatal(err)
		}
		b.ReportMetric(float64(len(resp.ConversationIDs)), "ids/op")
		b.ReportMetric(float64(proto.Size(resp)), "proto_bytes/op")
	}
}

func BenchmarkGetFullOwnerConversationIDsReadInactiveFilteredIDs(b *testing.B) {
	readInactiveConversationDuration = int64(time.Hour / time.Millisecond)
	b.Cleanup(func() {
		readInactiveConversationDuration = int64((30 * 24 * time.Hour) / time.Millisecond)
	})

	ids := newConversationIDs(50000)
	seqs := make(map[string]*pbmsg.FullSyncSeqs, len(ids))
	now := time.Now()
	for i, conversationID := range ids {
		if i%10 == 0 {
			seqs[conversationID] = unreadSeq(now.Add(-2 * time.Hour))
		} else {
			seqs[conversationID] = readInactiveSeq(now.Add(-2 * time.Hour))
		}
	}
	srv := &conversationServer{
		conversationDatabase: &fakeConversationDatabase{conversationIDs: ids},
		msgClient:            fakeMessageClient{seqs: seqs},
	}
	req := &pbconversation.GetFullOwnerConversationIDsReq{
		UserID: "customer-service",
		IdHash: 0,
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := srv.GetFullOwnerConversationIDs(testCtx("customer-service"), req)
		if err != nil {
			b.Fatal(err)
		}
		b.ReportMetric(float64(len(resp.ConversationIDs)), "ids/op")
		b.ReportMetric(float64(proto.Size(resp)), "proto_bytes/op")
	}
}

func TestGetFullOwnerConversationIDsReadInactiveSkipsSmallConversationSetsByDefault(t *testing.T) {
	ids := []string{"si_user_1", "si_user_2", "si_user_3"}
	srv := &conversationServer{
		conversationDatabase: &fakeConversationDatabase{conversationIDs: ids},
		msgClient: fakeMessageClient{seqs: map[string]*pbmsg.FullSyncSeqs{
			"si_user_1": readInactiveSeq(time.Now().Add(-2 * time.Hour)),
			"si_user_2": readInactiveSeq(time.Now().Add(-2 * time.Hour)),
			"si_user_3": readInactiveSeq(time.Now().Add(-2 * time.Hour)),
		}},
	}

	resp, err := srv.GetFullOwnerConversationIDs(testCtx("customer-service"), &pbconversation.GetFullOwnerConversationIDsReq{
		UserID: "customer-service",
		IdHash: 0,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !sameStrings(resp.ConversationIDs, ids) {
		t.Fatalf("expected small conversation sets to keep legacy IDs, got %v", resp.ConversationIDs)
	}
}

func readInactiveSeq(maxSeqTime time.Time) *pbmsg.FullSyncSeqs {
	return &pbmsg.FullSyncSeqs{
		HasReadSeq: 10,
		MaxSeq:     10,
		MaxSeqTime: maxSeqTime.UnixMilli(),
	}
}

func readActiveSeq(maxSeqTime time.Time) *pbmsg.FullSyncSeqs {
	return &pbmsg.FullSyncSeqs{
		HasReadSeq: 10,
		MaxSeq:     10,
		MaxSeqTime: maxSeqTime.UnixMilli(),
	}
}

func clearedSeq(maxSeqTime time.Time) *pbmsg.FullSyncSeqs {
	return &pbmsg.FullSyncSeqs{
		HasReadSeq: 10,
		MaxSeq:     10,
		MaxSeqTime: maxSeqTime.UnixMilli(),
		UserMinSeq: 11,
	}
}

func unreadSeq(maxSeqTime time.Time) *pbmsg.FullSyncSeqs {
	return &pbmsg.FullSyncSeqs{
		HasReadSeq: 9,
		MaxSeq:     10,
		MaxSeqTime: maxSeqTime.UnixMilli(),
	}
}

func newConversationIDs(count int) []string {
	ids := make([]string, 0, count)
	for i := 0; i < count; i++ {
		ids = append(ids, fmt.Sprintf("si_user_%05d", i))
	}
	return ids
}

func sameStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func withReadInactiveConversationCountThreshold(t *testing.T, threshold int) {
	t.Helper()
	old := readInactiveConversationCountThreshold
	readInactiveConversationCountThreshold = threshold
	t.Cleanup(func() {
		readInactiveConversationCountThreshold = old
	})
}

func withReadInactiveConversationFilterEnabled(t *testing.T, enabled bool) {
	t.Helper()
	old := readInactiveConversationFilterEnabled
	readInactiveConversationFilterEnabled = enabled
	t.Cleanup(func() {
		readInactiveConversationFilterEnabled = old
	})
}

func withReadInactiveConversationDuration(t *testing.T, duration int64) {
	t.Helper()
	old := readInactiveConversationDuration
	readInactiveConversationDuration = duration
	t.Cleanup(func() {
		readInactiveConversationDuration = old
	})
}

type fakeConversationDatabase struct {
	conversationIDs       []string
	pinnedConversationIDs []string
}

func (f *fakeConversationDatabase) UpdateUsersConversationField(context.Context, []string, string, map[string]any) error {
	return nil
}

func (f *fakeConversationDatabase) CreateConversation(context.Context, []*model.Conversation) error {
	return nil
}

func (f *fakeConversationDatabase) SyncPeerUserPrivateConversationTx(context.Context, []*model.Conversation) error {
	return nil
}

func (f *fakeConversationDatabase) FindConversations(_ context.Context, _ string, conversationIDs []string) ([]*model.Conversation, error) {
	conversations := make([]*model.Conversation, 0, len(conversationIDs))
	for _, conversationID := range conversationIDs {
		conversations = append(conversations, &model.Conversation{ConversationID: conversationID})
	}
	return conversations, nil
}

func (f *fakeConversationDatabase) GetUserAllConversation(context.Context, string) ([]*model.Conversation, error) {
	return nil, nil
}

func (f *fakeConversationDatabase) SetUserConversations(context.Context, string, []*model.Conversation) error {
	return nil
}

func (f *fakeConversationDatabase) SetUsersConversationFieldTx(context.Context, []string, *model.Conversation, map[string]any) error {
	return nil
}

func (f *fakeConversationDatabase) UpdateUserConversations(context.Context, string, map[string]any) error {
	return nil
}

func (f *fakeConversationDatabase) CreateGroupChatConversation(context.Context, string, []string, *model.Conversation) error {
	return nil
}

func (f *fakeConversationDatabase) GetConversationIDs(context.Context, string) ([]string, error) {
	return f.conversationIDs, nil
}

func (f *fakeConversationDatabase) GetUserConversationIDsHash(context.Context, string) (uint64, error) {
	return hashutil.IdHash(f.conversationIDs), nil
}

func (f *fakeConversationDatabase) GetAllConversationIDs(context.Context) ([]string, error) {
	return nil, nil
}

func (f *fakeConversationDatabase) GetAllConversationIDsNumber(context.Context) (int64, error) {
	return 0, nil
}

func (f *fakeConversationDatabase) PageConversationIDs(context.Context, pagination.Pagination) ([]string, error) {
	return nil, nil
}

func (f *fakeConversationDatabase) GetConversationsByConversationID(context.Context, []string) ([]*model.Conversation, error) {
	return nil, nil
}

func (f *fakeConversationDatabase) GetConversationIDsNeedDestruct(context.Context) ([]*model.Conversation, error) {
	return nil, nil
}

func (f *fakeConversationDatabase) GetConversationNotReceiveMessageUserIDs(context.Context, string) ([]string, error) {
	return nil, nil
}

func (f *fakeConversationDatabase) FindConversationUserVersion(context.Context, string, uint, int) (*model.VersionLog, error) {
	return nil, nil
}

func (f *fakeConversationDatabase) FindMaxConversationUserVersionCache(context.Context, string) (*model.VersionLog, error) {
	return &model.VersionLog{ID: primitive.NewObjectID(), LastUpdate: time.Now()}, nil
}

func (f *fakeConversationDatabase) GetOwnerConversation(context.Context, string, pagination.Pagination) (int64, []*model.Conversation, error) {
	return int64(len(f.conversationIDs)), nil, nil
}

func (f *fakeConversationDatabase) GetNotNotifyConversationIDs(context.Context, string) ([]string, error) {
	return nil, nil
}

func (f *fakeConversationDatabase) GetPinnedConversationIDs(context.Context, string) ([]string, error) {
	return f.pinnedConversationIDs, nil
}

func (f *fakeConversationDatabase) FindRandConversation(context.Context, int64, int) ([]*model.Conversation, error) {
	return nil, nil
}

func (f *fakeConversationDatabase) DeleteUsersConversations(context.Context, string, []string) error {
	return nil
}

var _ controller.ConversationDatabase = (*fakeConversationDatabase)(nil)

type fakeMessageClient struct {
	seqs map[string]*pbmsg.FullSyncSeqs
}

func (f fakeMessageClient) GetMaxSeqs(context.Context, []string) (map[string]int64, error) {
	return nil, nil
}

func (f fakeMessageClient) GetMsgByConversationIDs(context.Context, []string, map[string]int64) (map[string]*sdkws.MsgData, error) {
	return nil, nil
}

func (f fakeMessageClient) GetHasReadSeqs(context.Context, []string, string) (map[string]int64, error) {
	return nil, nil
}

func (f fakeMessageClient) GetConversationsFullSyncSeqs(_ context.Context, req *pbmsg.GetConversationsFullSyncSeqsReq) (*pbmsg.GetConversationsFullSyncSeqsResp, error) {
	seqs := make(map[string]*pbmsg.FullSyncSeqs, len(req.ConversationIDs))
	for _, conversationID := range req.ConversationIDs {
		if seq, ok := f.seqs[conversationID]; ok {
			seqs[conversationID] = seq
		}
	}
	return &pbmsg.GetConversationsFullSyncSeqsResp{Seqs: seqs}, nil
}

func (f fakeMessageClient) SetUserConversationMaxSeq(context.Context, string, []string, int64) error {
	return nil
}

func (f fakeMessageClient) SetUserConversationMin(context.Context, string, []string, int64) error {
	return nil
}

func (f fakeMessageClient) GetLastMessageSeqByTime(context.Context, string, int64) (int64, error) {
	return 0, nil
}

func (f fakeMessageClient) GetLastMessage(context.Context, *pbmsg.GetLastMessageReq, ...grpc.CallOption) (*pbmsg.GetLastMessageResp, error) {
	return &pbmsg.GetLastMessageResp{}, nil
}
