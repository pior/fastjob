package tests

import (
	"context"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/stretchr/testify/require"

	"github.com/pior/fastjob"
)

const subscriptionName = "sub-test"
const topicName = "topic-test"

var sentinel = &Sentinel{}

func NewMockJob() fastjob.Job {
	return &MockJob{}
}

type MockJob struct{}

func (m *MockJob) Name() string {
	return "MockJob"
}

func (m *MockJob) Perform(ctx context.Context) error {
	sentinel.Touch()
	return nil
}

func TestLocalRunner(t *testing.T) {
	ctx := context.Background()
	runner := fastjob.LocalRunner{}
	job := &MockJob{}

	sentinel.Reset()

	runner.Enqueue(ctx, job)

	require.NoError(t, sentinel.Wait())
}

func TestPubSubRunner(t *testing.T) {
	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, "fake-projectid-emulator")
	require.NoError(t, err)
	defer client.Close()

	ensureResourceReady(t, client, topicName, subscriptionName)

	sentinel.Reset()

	// Enqueue job
	runner := fastjob.NewPubSubRunner(client, topicName)
	err = runner.Enqueue(ctx, &MockJob{})
	require.NoError(t, err)

	// Run the worker
	registry := fastjob.NewRegistry(NewMockJob)
	executor := &fastjob.Executor{}
	worker := fastjob.NewWorker(client.Subscription(subscriptionName), registry, executor)

	wctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	go worker.Run(wctx)

	require.NoError(t, sentinel.Wait())
}
