package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"

	"cloudevents-explorer/internal/types"
)

type PullParams struct {
	EmulatorHost   string `json:"emulatorHost"`
	ProjectID      string `json:"projectId"`
	SubscriptionID string `json:"subscriptionId"`
	MaxMessages    int    `json:"maxMessages"`
}

type PullResult struct {
	Messages []types.CloudEvent `json:"messages"`
	Count    int                `json:"count"`
}

func Pull(params PullParams) (*PullResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	os.Setenv("PUBSUB_EMULATOR_HOST", params.EmulatorHost)

	client, err := pubsub.NewClient(ctx, params.ProjectID,
		option.WithEndpoint(params.EmulatorHost),
		option.WithoutAuthentication(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

	subscription := client.Subscription(params.SubscriptionID)

	messages := []types.CloudEvent{}
	var msgMu sync.Mutex

	receiveCtx, receiveCancel := context.WithTimeout(ctx, 5*time.Second)
	defer receiveCancel()

	err = subscription.Receive(receiveCtx, func(ctx context.Context, msg *pubsub.Message) {
		event := types.CloudEvent{
			ID:        msg.ID,
			Type:      msg.Attributes["ce-type"],
			Subject:   msg.Attributes["ce-subject"],
			Source:    msg.Attributes["ce-source"],
			Schema:    msg.Attributes["ce-dataschema"],
			Published: msg.PublishTime.Format(time.RFC3339),
			Timestamp: msg.PublishTime.Unix(),
		}

		if len(msg.Data) > 0 {
			var data map[string]interface{}
			if err := json.Unmarshal(msg.Data, &data); err == nil {
				event.Data = data
			}
		}

		msgMu.Lock()
		messages = append(messages, event)
		msgMu.Unlock()

		msg.Ack()

		if len(messages) >= params.MaxMessages {
			receiveCancel()
		}
	})

	if err != nil && err != context.Canceled && err != context.DeadlineExceeded {
		return nil, fmt.Errorf("failed to receive: %w", err)
	}

	return &PullResult{
		Messages: messages,
		Count:    len(messages),
	}, nil
}
