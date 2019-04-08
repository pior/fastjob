package fastjob_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/pior/fastjob"
)

type Sentinel struct {
	value int
}

func (s *Sentinel) Reset() {
	s.value = 0
}

func (s *Sentinel) Touch() {
	s.value++
}

func (s *Sentinel) WaitFor(timeout time.Duration, count int) error {
	deadline := time.Now().Add(timeout)
	sleepTime := time.Millisecond * 100
	for {
		if s.value >= count {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("sentinel was not incremented after %s", timeout)
		}
		time.Sleep(sleepTime)
	}
}

func (s *Sentinel) Wait() error {
	return s.WaitFor(time.Second*5, 1)
}

type pubsubHelper struct {
	projectID string
	client    *pubsub.Client

	topicName string
	topic     *pubsub.Topic

	subscriptionName string
	subscription     *pubsub.Subscription
}

func (h *pubsubHelper) GetClient() *pubsub.Client {
	if h.client == nil {
		var err error
		h.client, err = pubsub.NewClient(context.Background(), h.projectID)
		if err != nil {
			log.Panicf("failed to create pubsub client: %s", err)
		}
	}
	return h.client
}

func (h *pubsubHelper) WithRandomTopic() *pubsubHelper {
	return h.WithTopic(fmt.Sprintf("topic-%d", time.Now().UnixNano()))
}

func (h *pubsubHelper) WithTopic(topicName string) *pubsubHelper {
	h.topicName = topicName
	h.topic = h.GetClient().Topic(topicName)
	h.subscriptionName = "sub-" + topicName
	h.subscription = h.GetClient().Subscription(h.subscriptionName)
	return h
}

func (h *pubsubHelper) CreateResources() *pubsubHelper {
	var err error
	ctx := context.Background()
	h.topic, err = h.GetClient().CreateTopic(ctx, h.topicName)
	if err != nil {
		log.Panicf("failed to create a topic: %s", err)
	}

	subConfig := pubsub.SubscriptionConfig{Topic: h.topic}
	h.subscription, err = h.GetClient().CreateSubscription(ctx, h.subscriptionName, subConfig)
	if err != nil {
		log.Panicf("failed to create a subscription: %s", err)
	}

	return h
}

func (h *pubsubHelper) Runner() fastjob.Runner {
	return fastjob.NewPubSubRunner(h.GetClient(), h.topicName)
}

func (h *pubsubHelper) Worker() *fastjob.PubsubWorker {
	registry := fastjob.NewRegistry().WithJobs(&MockJob{})
	config := fastjob.NewConfig(registry)
	return fastjob.NewPubsubWorker(config, h.subscription)
}
