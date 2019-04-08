package fastjob_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/pior/fastjob"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewTestJob() *TestJob {
	return &TestJob{}
}

type TestJob struct {
	ObjectID  int
	Operation string
}

func (m *TestJob) Name() string {
	return "TestJob"
}

func (m *TestJob) Perform(ctx context.Context) error {
	return nil
}

func TestNewJobRequest(t *testing.T) {
	job := NewTestJob()
	job.ObjectID = 42
	job.Operation = "work"

	req, err := fastjob.NewJobRequest(job)
	require.NoError(t, err)
	assert.Equal(t, "TestJob", req.JobName)
	assert.Equal(t, "{\"ObjectID\":42,\"Operation\":\"work\"}", string(req.JobData))
	assert.Contains(t, req.String(), "<TestJob-")

	newJob := NewTestJob()
	err = json.Unmarshal(req.JobData, newJob)
	require.NoError(t, err)
	require.Equal(t, job, newJob)
}
