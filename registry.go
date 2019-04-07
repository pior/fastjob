package fastjob

import "fmt"

type JobRegistry struct {
	jobTypes map[string]JobType
}

func NewRegistry(jobTypes ...JobType) *JobRegistry {
	reg := &JobRegistry{
		jobTypes: make(map[string]JobType),
	}
	for _, jobType := range jobTypes {
		reg.Register(jobType)
	}
	return reg
}

func (r *JobRegistry) Register(jobType JobType) {
	jobName := jobType().Name()
	if jobName == "" {
		panic(fmt.Sprintf("Job %T.Name() cannot be empty string", jobType))
	}
	r.jobTypes[jobType().Name()] = jobType
}

func (r *JobRegistry) Get(name string) (JobType, error) {
	jobType, ok := r.jobTypes[name]
	if !ok {
		return nil, fmt.Errorf("JobType %q does not exist", name)
	}
	return jobType, nil
}
