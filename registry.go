package lazy

import "fmt"

type JobType func() Job

type JobRegistry struct {
	jobTypes map[string]func() Job
}

func (r *JobRegistry) Register(jobType func() Job) {
	jobName := jobType().Name()
	if jobName == "" {
		panic(fmt.Sprintf("Job %T.Name() cannot be empty string", jobType))
	}
	r.jobTypes[jobType().Name()] = jobType
}

func (r *JobRegistry) Get(name string) (func() Job, error) {
	jobType, ok := r.jobTypes[name]
	if !ok {
		return nil, fmt.Errorf("JobType %q does not exist", name)
	}
	return jobType, nil
}
