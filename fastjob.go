package fastjob

import (
	"context"
)

// Job is the interface that jobs must implement
type Job interface {
	Name() string
	Perform(context.Context) error
}

// JobType is the job factory. Used to allocate a job struct when receiving a job request to process.
type JobType func() Job
