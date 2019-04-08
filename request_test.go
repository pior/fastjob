package fastjob_test

import (
	"encoding/json"
	"testing"

	"github.com/pior/fastjob"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJobRequest(t *testing.T) {
	job := &MockJob{}
	job.Value = 42

	req, err := fastjob.NewJobRequest(job)
	require.NoError(t, err)
	assert.Equal(t, "MockJob", req.JobName)
	assert.Equal(t, "{\"Value\":42}", string(req.JobData))
	assert.Contains(t, req.String(), "<MockJob-")

	newJob := &MockJob{}
	err = json.Unmarshal(req.JobData, newJob)
	require.NoError(t, err)
	require.Equal(t, job.Value, newJob.Value)
}
