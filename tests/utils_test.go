package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/stretchr/testify/require"
)

type Sentinel struct {
	value int
}

func (s *Sentinel) Reset() {
	s.value = 0
}

func (s *Sentinel) Touch() {
	s.value++
}

func (s *Sentinel) WaitFor(timeout time.Duration, count int) error {
	deadline := time.Now().Add(timeout)
	sleepTime := time.Millisecond * 100
	for {
		if s.value >= count {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("sentinel was not incremented after %s", timeout)
		}
		time.Sleep(sleepTime)
	}
}

func (s *Sentinel) Wait() error {
	return s.WaitFor(time.Second*5, 1)
}

func ensureResourceReady(t *testing.T, client *pubsub.Client, topicName, subscriptionName string) {
	ctx := context.Background()

	topic := client.Topic(topicName)

	exists, err := topic.Exists(ctx)
	require.NoError(t, err)

	if !exists {
		topic, err = client.CreateTopic(ctx, topicName)
		require.NoError(t, err)
	}

	exists, err = client.Subscription(subscriptionName).Exists(ctx)
	require.NoError(t, err)

	if !exists {
		subConfig := pubsub.SubscriptionConfig{Topic: topic}
		_, err = client.CreateSubscription(ctx, subscriptionName, subConfig)
		require.NoError(t, err)
	}
}
