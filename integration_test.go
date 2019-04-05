package lazy_test

import (
	"context"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/stretchr/testify/require"

	"github.com/pior/lazy"
)

func NewMockJob() lazy.Job {
	return &MockJob{}
}

type MockJob struct {
	jobFunc func(context.Context) error
	calls   int
}

func (m *MockJob) Name() string {
	return "MockJob"
}

func (m *MockJob) Perform(ctx context.Context) error {
	m.calls++
	return nil
}

func TestLocalRunner(t *testing.T) {
	ctx := context.Background()
	runner := lazy.LocalRunner{}
	job := &MockJob{}

	runner.Enqueue(ctx, job)
	time.Sleep(time.Millisecond * 50)
	require.Equal(t, 1, job.calls)
}

func TestPubSubRunner(t *testing.T) {
	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, "gcloud-emulator")
	require.NoError(t, err)

	runner := lazy.NewPubSubRunner(client, "test-topic")

	job := &MockJob{}

	runner.Enqueue(ctx, job)

	registry := &lazy.JobRegistry{}
	registry.Register(NewMockJob)

	executor := &lazy.Executor{}

	// TODO
	// subscription :=

	worker := lazy.NewWorker(subscription, registry, executor)

	wctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	worker.Run(wctx)

	time.Sleep(time.Millisecond * 50)
	require.Equal(t, 1, job.calls)
}
