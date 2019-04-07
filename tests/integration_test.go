package tests

import (
	"context"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/stretchr/testify/require"

	"github.com/pior/fastjob"
)

var sentinel = &Sentinel{}

func NewMockJob() fastjob.Job {
	return &MockJob{}
}

type MockJob struct {
	err error
}

func (m *MockJob) Name() string {
	return "MockJob"
}

func (m *MockJob) Perform(ctx context.Context) error {
	sentinel.Touch()
	return m.err
}

func TestLocalRunner(t *testing.T) {
	ctx := context.Background()

	registry := fastjob.NewRegistry(NewMockJob)
	config := fastjob.NewConfig(registry)
	runner := fastjob.NewLocalRunner(config)

	sentinel.Reset()

	job := &MockJob{}
	err := runner.Enqueue(ctx, job)
	require.NoError(t, err)

	require.NoError(t, sentinel.Wait())
}

func TestPubSubRunner(t *testing.T) {
	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, "fake-id")
	require.NoError(t, err)
	defer client.Close()

	topicName, subscriptionName := prepareTopicSub(t, client)

	sentinel.Reset()

	// Enqueue job
	runner := fastjob.NewPubSubRunner(client, topicName)

	err = runner.Enqueue(ctx, &MockJob{})
	require.NoError(t, err)

	// Run the worker
	registry := fastjob.NewRegistry(NewMockJob)
	config := fastjob.NewConfig(registry)
	sub := client.Subscription(subscriptionName)
	worker := fastjob.NewWorker(config, sub)

	wctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	go func() {
		err := worker.Run(wctx)
		require.NoError(t, err)
	}()

	require.NoError(t, sentinel.Wait())
}
