package fastjob

import (
	"context"
	"encoding/json"
	"time"

	"cloud.google.com/go/pubsub"
)

// PubsubWorker
type PubsubWorker struct {
	subscription *pubsub.Subscription
	executor     *Executor
	logger       Logger
}

// NewPubsubWorker creates a new PubSub worker to execute jobs enqueued in a PubSub Topic.
func NewPubsubWorker(config *config, subscription *pubsub.Subscription) *PubsubWorker {
	return &PubsubWorker{subscription, NewExecutor(config), config.logger}
}

// Run runs the Pubsub worker until an error occurs or the context is cancelled
func (w *PubsubWorker) Run(ctx context.Context) error {
	w.logger.Infof(ctx, "Connecting to PubSub (subscription: %s)", w.subscription)

	err := w.subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		w.logger.Debugf(ctx, "processing message %s (enqueued at %s)", msg.ID, msg.PublishTime)

		request := &JobRequest{}
		err := json.Unmarshal(msg.Data, request)
		if err != nil {
			w.logger.Errorf(ctx, err, "failed to unmarshal the job request from message %s", msg.ID)
			w.handleInvalidRequest(ctx, msg)
		}

		err = w.executor.Execute(ctx, request)
		if err != nil {
			w.logger.Errorf(ctx, err, "failed to execute job %s", request)
			w.handleJobFailure(ctx, msg, request)
		}

		msg.Ack()
	})

	if err != nil {
		w.logger.Errorf(ctx, err, "Subscription receive stopped")
	} else {
		w.logger.Infof(ctx, "Connection to PubSub was closed")
	}
	return err
}

// TODO: make the error handler a changeable component

func (w *PubsubWorker) handleInvalidRequest(ctx context.Context, msg *pubsub.Message) {
	// By default, do nothing, not even Nack the message
	// It will get dispatch again after the default lease period
}

func (w *PubsubWorker) handleJobFailure(ctx context.Context, msg *pubsub.Message, request *JobRequest) {
	// Block the message for some time, then Nack it to get it redelivered
	// This consumes one goroutine (4k on stack) per sleeping message
	time.Sleep(time.Second * 10)
	msg.Nack()
}
