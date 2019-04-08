package fastjob_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/pior/fastjob"
)

var sentinel = &Sentinel{}
var helper *pubsubHelper

func init() {
	helper = &pubsubHelper{projectID: "fake-id"}
}

func NewMockJob() *MockJob {
	return &MockJob{}
}

type MockJob struct {
	err error
}

func (m *MockJob) Name() string {
	return "MockJob"
}

func (m *MockJob) Perform(ctx context.Context) error {
	sentinel.Touch()
	return m.err
}

func TestLocalRunner(t *testing.T) {
	ctx := context.Background()

	registry := fastjob.NewRegistry().WithJobs(&MockJob{})
	config := fastjob.NewConfig(registry)
	runner := fastjob.NewLocalRunner(config)

	sentinel.Reset()

	job := &MockJob{}
	err := runner.Enqueue(ctx, job)
	require.NoError(t, err)

	require.NoError(t, sentinel.Wait())
}

func TestPubSubRunner(t *testing.T) {
	ctx := context.Background()

	helper.WithRandomTopic().CreateResources()

	sentinel.Reset()

	err := helper.Runner().Enqueue(ctx, &MockJob{})
	require.NoError(t, err)

	wctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		err := helper.Worker().Run(wctx)
		require.NoError(t, err)
		wg.Done()
	}()

	require.NoError(t, sentinel.Wait())
	wg.Wait()
}

func TestPubSubWorkerStop(t *testing.T) {
	helper.WithRandomTopic().CreateResources()

	ctx := context.Background()
	wctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	err := helper.Worker().Run(wctx)
	require.NoError(t, err)
}

func TestPubSubWorkerError(t *testing.T) {
	helper.WithTopic("nopenope")

	ctx := context.Background()

	wctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	err := helper.Worker().Run(wctx)
	require.EqualError(t, err, "rpc error: code = NotFound desc = Subscription does not exist (resource=sub-nopenope)")
}
