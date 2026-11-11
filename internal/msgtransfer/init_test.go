package msgtransfer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/openimsdk/tools/mq"
)

type retryTestMessage struct {
	marked int
}

func (m *retryTestMessage) Context() context.Context { return context.Background() }
func (m *retryTestMessage) Key() string              { return "retry-test" }
func (m *retryTestMessage) Value() []byte            { return nil }
func (m *retryTestMessage) Mark()                    { m.marked++ }
func (m *retryTestMessage) Commit()                  {}

func TestConsumeMongoMessageRetriesBeforeMarking(t *testing.T) {
	msg := &retryTestMessage{}
	attempts := 0
	err := consumeMongoMessage(context.Background(), msg, func(mq.Message) error {
		attempts++
		if attempts < 3 {
			return errors.New("MongoDB unavailable")
		}
		return nil
	}, 0, 0)
	if err != nil {
		t.Fatalf("consumeMongoMessage() error = %v", err)
	}
	if attempts != 3 {
		t.Fatalf("attempts = %d, want 3", attempts)
	}
	if msg.marked != 1 {
		t.Fatalf("Mark calls = %d, want 1 after persistence succeeds", msg.marked)
	}
}

func TestConsumeMongoMessageDoesNotMarkWhenContextEnds(t *testing.T) {
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(errors.New("test shutdown"))
	msg := &retryTestMessage{}
	err := consumeMongoMessage(ctx, msg, func(mq.Message) error {
		return errors.New("MongoDB unavailable")
	}, time.Hour, time.Hour)
	if err == nil {
		t.Fatal("consumeMongoMessage() error = nil, want cancellation cause")
	}
	if msg.marked != 0 {
		t.Fatalf("Mark calls = %d, want 0 when persistence never succeeds", msg.marked)
	}
}
