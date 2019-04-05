package lazy

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
)

type Worker struct {
	subscription *pubsub.Subscription
	registry     *JobRegistry
	executor     *Executor
}

func NewWorker(subscription *pubsub.Subscription, registry *JobRegistry, executor *Executor) *Worker {
	return &Worker{subscription, registry, executor}
}

// func NewWorker(ctx context.Context, projectID, subscriptionName string) (*Worker, error) {
// 	client, err := pubsub.NewClient(ctx, projectID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	sub := client.Subscription(subscriptionName)
// 	sub.ReceiveSettings.MaxOutstandingMessages = 50
// 	sub.ReceiveSettings.MaxExtension = 20 * time.Second

// 	w := &Worker{
// 		subscription: sub,
// 		registry:     registry,

// 	}
// 	return w, nil
// }

func (w *Worker) Run(ctx context.Context) error {
	return w.subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		job, err := w.prepareJob(ctx, msg)
		if err != nil {
			msg.Nack()
		}

		err = w.executor.Execute(ctx, job)
		if err != nil {
			msg.Nack()
		}

		msg.Ack()
	})
}

func (w *Worker) prepareJob(ctx context.Context, msg *pubsub.Message) (Job, error) {
	req := &JobRequest{}
	err := json.Unmarshal(msg.Data, req)
	if err != nil {
		return nil, err
	}

	jobType, err := w.registry.Get(req.JobName)
	if err != nil {
		return nil, err
	}

	job := jobType()
	err = json.Unmarshal(req.JobData, &job)
	if err != nil {
		return nil, err
	}

	return job, nil
}
