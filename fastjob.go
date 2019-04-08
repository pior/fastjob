package fastjob

import (
	"context"
)

// Job is the interface that jobs must implement
type Job interface {
	Name() string
	Perform(context.Context) error
}
