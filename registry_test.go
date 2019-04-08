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
	reg := fastjob.NewRegistry().WithJobs(&MockJob{})

	jobType, err := reg.Get("MockJob")
	require.NoError(t, err)
	require.Equal(t, NewMockJob().Name(), jobType().Name())
}

func TestJobRegistryInvalidJob(t *testing.T) {
	require.Panics(t, func() {
		fastjob.NewRegistry().WithJobs(nil)
	})

	require.Panics(t, func() {
		fastjob.NewRegistry().WithJobs(&NoNameJob{})
	})
}

func TestJobRegistryNotFound(t *testing.T) {
	reg := fastjob.NewRegistry()

	jobType, err := reg.Get("NotAJob")
	require.Error(t, err)
	require.Nil(t, jobType)
}
