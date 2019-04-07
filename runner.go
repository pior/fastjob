package fastjob

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
)

type Runner interface {
	Enqueue(context.Context, Job) error
}

// LocalRunner runs jobs synchronously
type LocalRunner struct {
	executor *Executor
	logger   Logger
}

func NewLocalRunner(config *config) Runner {
	return &LocalRunner{
		executor: NewExecutor(config),
		logger:   config.logger,
	}
}

func (r *LocalRunner) Enqueue(ctx context.Context, job Job) error {
	// LocalRunner is probably used during development and test
	// So let's re-use the same flow as the PubSub runner:
	// Job -> JobRequest -> (queue) -> lookup in JobRegistry -> Job.Perform()

	request, err := NewJobRequest(job)
	if err != nil {
		return err
	}

	err = r.executor.Execute(ctx, request)
	return err
}

// PubSubRunner enqueues jobs through GCP pubsub
type PubSubRunner struct {
	client    *pubsub.Client
	topicName string
}

func NewPubSubRunner(client *pubsub.Client, topicName string) Runner {
	return &PubSubRunner{client, topicName}
}

func (r *PubSubRunner) Enqueue(ctx context.Context, job Job) error {
	request, err := NewJobRequest(job)
	if err != nil {
		return err
	}

	data, err := json.Marshal(request)
	if err != nil {
		return err
	}

	msg := &pubsub.Message{Data: data}

	topic := r.client.Topic(r.topicName)
	defer topic.Stop()

	publishResult := topic.Publish(ctx, msg)
	_, err = publishResult.Get(ctx)
	return err
}
