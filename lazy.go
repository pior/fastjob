package fastjob

import (
	"context"

	"cloud.google.com/go/pubsub"
)

// type JobErrorHandler func(*JobRequest, error)
// type JobSuccessHandler func(*JobRequest, *Job)

// Job is the interface that jobs must implement
type Job interface {
	Name() string
	Perform(context.Context) error
}

// NewWorker creates a new PubSub worker to execute jobs enqueued in a PubSub Topic.
func NewWorker(subscription *pubsub.Subscription, jobRegistry *JobRegistry, executor *Executor, logger Logger) *Worker {
	if executor == nil {
		executor = &Executor{}
	}
	if logger == nil {
		logger = &standardLogger{}
	}
	return &Worker{subscription, jobRegistry, executor, logger}
}
