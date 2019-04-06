package fastjob_test

import (
	"context"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/pior/fastjob"
)

func NewMockJob() fastjob.Job {
	return &MockJob{}
}

type MockJob struct{}

func (m *MockJob) Name() string {
	return "MockJob"
}

func (m *MockJob) Perform(ctx context.Context) error {
	return nil
}

func ExampleNewWorker() {
	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, "my-gcp-project-id")
	if err != nil {
		// TODO: Handle error.
	}

	sub := client.Subscription("subscription-test")

	// https://github.com/googleapis/google-cloud-go/wiki/Fine-Tuning-PubSub-Receive-Performance
	// Defaults:
	//   MaxExtension:           10 * time.Minute,
	//   MaxOutstandingMessages: 1000,
	//   MaxOutstandingBytes:    1e9,
	//   NumGoroutines:          1,
	sub.ReceiveSettings.MaxOutstandingMessages = 100
	sub.ReceiveSettings.MaxOutstandingBytes = 100e6
	sub.ReceiveSettings.MaxExtension = 1 * time.Minute
	sub.ReceiveSettings.NumGoroutines = 1

	jobRegistry := fastjob.NewRegistry()
	// jobRegistry.Register(NewMockJob)

	worker := fastjob.NewWorker(sub, jobRegistry, nil, nil)

	worker.Run(ctx)
}
