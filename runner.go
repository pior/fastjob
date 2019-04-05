package lazy

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
}

func (r *LocalRunner) Enqueue(ctx context.Context, job Job) error {
	go func() {
		r.executor.Execute(ctx, job)
	}()

	return nil
}

// PubSubRunner enqueues jobs through GCP pubsub
type PubSubRunner struct {
	client    *pubsub.Client
	topicName string
}

func NewPubSubRunner(client *pubsub.Client, topicName string) *PubSubRunner {
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
