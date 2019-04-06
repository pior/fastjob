package fastjob

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type JobRequest struct {
	RequestID   string
	RequestTime int64

	JobName string
	JobData []byte
}

func NewJobRequest(job Job) (*JobRequest, error) {
	data, err := json.Marshal(job)
	if err != nil {
		return nil, err
	}

	req := &JobRequest{
		RequestID:   uuid.New().String(),
		RequestTime: time.Now().Unix(),

		JobName: job.Name(),
		JobData: data,
	}

	return req, nil
}
