package fastjob

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
)

// Worker
type Worker struct {
	subscription *pubsub.Subscription
	jobRegistry  *JobRegistry
	executor     *Executor
	logger       Logger
}

func (w *Worker) Run(ctx context.Context) error {
	w.logger.Infof(ctx, "Connecting to PubSub")

	return w.subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		w.logger.Debugf(ctx, "preparing job for msg: %s", msg.Data)
		job, err := w.prepareJob(ctx, msg)
		if err != nil {
			w.logger.Errorf(ctx, err, "failed to prepare job")
			msg.Nack()
		}

		w.logger.Debugf(ctx, "executing job %s", job.Name())
		err = w.executor.Execute(ctx, job)
		if err != nil {
			w.logger.Errorf(ctx, err, "failed to execute job %s", job.Name())
			msg.Nack()
		}

		w.logger.Debugf(ctx, "succeed to execute job: %s", job.Name())
		msg.Ack()
	})
}

func (w *Worker) prepareJob(ctx context.Context, msg *pubsub.Message) (Job, error) {
	req := &JobRequest{}
	err := json.Unmarshal(msg.Data, req)
	if err != nil {
		return nil, err
	}

	jobType, err := w.jobRegistry.Get(req.JobName)
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
