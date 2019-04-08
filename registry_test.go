package fastjob_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pior/fastjob"
)

type NoNameJob struct{}

func (m *NoNameJob) Name() string {
	return ""
}

func (m *NoNameJob) Perform(ctx context.Context) error {
	return nil
}

func TestJobRegistry(t *testing.T) {
	job := &MockJob{Value: 42}
	reg := fastjob.NewRegistry().WithJobs(job)

	jobType, err := reg.Get(job.Name())
	require.NoError(t, err)
	require.Equal(t, job.Name(), jobType().Name())
}

func TestJobRegistryInvalidJob(t *testing.T) {
	require.Panics(t, func() {
		fastjob.NewRegistry().WithJobs(nil)
	})

	require.PanicsWithValue(t, "The name of *fastjob_test.NoNameJob cannot be empty string", func() {
		fastjob.NewRegistry().WithJobs(&NoNameJob{})
	})
}

func TestJobRegistryNotFound(t *testing.T) {
	reg := fastjob.NewRegistry()

	jobType, err := reg.Get("NotAJob")
	require.Error(t, err)
	require.Nil(t, jobType)
}
