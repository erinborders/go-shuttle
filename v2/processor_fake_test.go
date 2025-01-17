package shuttle_test

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"

	v2 "github.com/Azure/go-shuttle/v2"
)

var _ v2.MessageSettler = &fakeSettler{}

type fakeSettler struct {
	AbandonCalled    int
	CompleteCalled   int
	DeadLetterCalled int
	DeferCalled      int
	RenewCalled      int
}

func (f *fakeSettler) AbandonMessage(ctx context.Context, message *azservicebus.ReceivedMessage, options *azservicebus.AbandonMessageOptions) error {
	f.AbandonCalled++
	return nil
}

func (f *fakeSettler) CompleteMessage(ctx context.Context, message *azservicebus.ReceivedMessage, options *azservicebus.CompleteMessageOptions) error {
	f.CompleteCalled++
	return nil
}

func (f *fakeSettler) DeadLetterMessage(ctx context.Context, message *azservicebus.ReceivedMessage, options *azservicebus.DeadLetterOptions) error {
	f.DeadLetterCalled++
	return nil
}

func (f *fakeSettler) DeferMessage(ctx context.Context, message *azservicebus.ReceivedMessage, options *azservicebus.DeferMessageOptions) error {
	f.DeferCalled++
	return nil
}

func (f *fakeSettler) RenewMessageLock(ctx context.Context, message *azservicebus.ReceivedMessage, options *azservicebus.RenewMessageLockOptions) error {
	f.RenewCalled++
	return nil
}

type fakeReceiver struct {
	// outcomes to verify
	ReceiveCalls []int // array of maxMessage value passed to receive calls in the lifetime of the fake receiver

	// configure fake
	SetupReceiveError     error
	SetupReceivedMessages chan *azservicebus.ReceivedMessage
	*fakeSettler
	SetupMaxReceiveCalls int
}

func (f *fakeReceiver) ReceiveMessages(_ context.Context, maxMessages int, _ *azservicebus.ReceiveMessagesOptions) ([]*azservicebus.ReceivedMessage, error) {
	f.ReceiveCalls = append(f.ReceiveCalls, maxMessages)
	if maxMessages == 0 && len(f.SetupReceivedMessages) > 0 {
		return nil, nil
	}
	var result []*azservicebus.ReceivedMessage
	for msg := range f.SetupReceivedMessages {
		result = append(result, msg)
		if len(result) == maxMessages || len(f.SetupReceivedMessages) == 0 {
			break
		}
	}

	// return an error if we request more messages than there are available.
	if len(f.ReceiveCalls) >= f.SetupMaxReceiveCalls {
		return result, fmt.Errorf("max receive calls exceeded")
	}

	return result, f.SetupReceiveError
}
