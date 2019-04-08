package fastjob

import (
	"fmt"
	"reflect"
)

// JobType is the job factory. Used to allocate a job struct when receiving a job request to process.
type JobType func() Job

type JobRegistry struct {
	jobTypes map[string]JobType
}

func NewRegistry() *JobRegistry {
	return &JobRegistry{make(map[string]JobType)}
}

func (r *JobRegistry) WithJobs(jobs ...Job) *JobRegistry {
	for _, job := range jobs {
		r.WithJob(job)
	}
	return r
}

func (r *JobRegistry) WithJob(job Job) *JobRegistry {
	if job.Name() == "" {
		panic(fmt.Sprintf("The name of %T cannot be empty string", job))
	}

	jobType := reflect.TypeOf(job).Elem()

	r.jobTypes[job.Name()] = func() Job { return reflect.New(jobType).Interface().(Job) }
	return r
}

func (r *JobRegistry) Get(name string) (JobType, error) {
	jobType, ok := r.jobTypes[name]
	if !ok {
		return nil, fmt.Errorf("JobType %q does not exist", name)
	}
	return jobType, nil
}
