package lazy

import "context"

type Executor struct {
}

// Execute runs a job,  recover panics,
func (e *Executor) Execute(ctx context.Context, job Job) error {

	err := job.Perform(ctx)

	return err
}
