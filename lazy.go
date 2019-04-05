package lazy

import "context"

// type JobErrorHandler func(*JobRequest, error)
// type JobSuccessHandler func(*JobRequest, *Job)

type Job interface {
	Name() string
	Perform(context.Context) error
}
