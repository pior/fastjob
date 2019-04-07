package fastjob

import (
	"context"
	"encoding/json"
)

type Executor struct {
	registry *JobRegistry
	logger   Logger
}

func NewExecutor(config *config) *Executor {
	return &Executor{config.registry, config.logger}
}

// Execute runs a job,  recover panics,
func (e *Executor) Execute(ctx context.Context, request *JobRequest) error {
	jobType, err := e.registry.Get(request.JobName)
	if err != nil {
		e.logger.Errorf(ctx, err, "no jobType exists in registry for %s", request)
		return err
	}

	job := jobType()
	err = json.Unmarshal(request.JobData, &job)
	if err != nil {
		e.logger.Errorf(ctx, err, "failed to unmarshal job for %s", request)
		return err
	}

	e.logger.Debugf(ctx, "executing %s", request)
	err = job.Perform(ctx)

	if err != nil {
		e.logger.Errorf(ctx, err, "failed to execute %s", request)
	} else {
		e.logger.Debugf(ctx, "completed %s successfully", request)
	}

	return err
}
