# fastjob

[![GoDoc](https://godoc.org/github.com/pior/fastjob?status.svg)](https://godoc.org/github.com/pior/fastjob)
[![Go Report Card](https://goreportcard.com/badge/github.com/pior/fastjob)](https://goreportcard.com/report/github.com/pior/fastjob)
[![CircleCI](https://circleci.com/gh/pior/fastjob.svg?style=svg)](https://circleci.com/gh/pior/fastjob)
[![codecov](https://codecov.io/gh/pior/fastjob/branch/master/graph/badge.svg)](https://codecov.io/gh/pior/fastjob)

Fastjob is a fast and robust job queue using GoogleCloud PubSub 🛰

**Work In Progress**

Design objectives:
- Robustness: never lose a job.
- Reliability: never let the main queue be blocked by failing jobs.

Strategies:
- Robustness: only one external dependencies: PubSub.
- Robustness: the durability is garanteed by PubSub.
- Reliability: the core features are mostly only the PubSub semantics and features (but extensible).
- Reliability: route failing jobs to a dead letter queue.

## Usage

#### Define a job:

```golang
type PingHTTP struct{
    Url string
}

func (m *PingHTTP) Name() string {
	return "PingHTTP"
}

func (m *PingHTTP) Perform(ctx context.Context) error {
    _, err := http.Post(m.Url)
	return err
}
```

Note: the job will be JSON encoded, only **public** fields should be used to define the job inputs.


#### Register the job:

```golang
registry := fastjob.NewRegistry().WithJob(&PingHTTP{})
```

#### Run the worker:

```golang
client, err := pubsub.NewClient(ctx, "my-gcp-project-id")
sub := client.Subscription("sub-test")
config := fastjob.NewConfig(registry)

worker := fastjob.NewPubsubWorker(config, sub)
worker.Run(ctx)
```

#### Enqueue a job:

```golang
runner := fastjob.NewPubSubRunner(client, topicName)

job := &PingHTTP{Url: "http://example.org/hello"}
err = runner.Enqueue(ctx, job)
```

#### Use a local runner for testing:

```golang
runner := fastjob.NewLocalRunner()

err = runner.Enqueue(ctx, job)
```

## License

[MIT](LICENSE)
