package main

import (
	"context"
	"github.com/google/uuid"
	kubemq "github.com/kubemq-io/kubemq-go"
	"time"
)

type KubeMQ struct {
	client *kubemq.Client
}

func NewKubeMQClient(host string, port int) (*KubeMQ, error) {
	client, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress(host, port),
		kubemq.WithClientId(uuid.New().String()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	if err != nil {
		return nil, err
	}
	k := &KubeMQ{
		client: client,
	}
	return k, nil
}

func (k *KubeMQ) StartListen(ctx context.Context, channel, group string, msgCh chan *kubemq.EventStoreReceive, errCh chan error) error {
	eventsCh, err := k.client.SubscribeToEventsStore(ctx, channel, group, errCh, kubemq.StartFromTimeDelta(1*time.Minute))
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case event, more := <-eventsCh:
				if !more {
					return
				}
				msgCh <- event

			case <-ctx.Done():
				return
			}

		}
	}()
	return nil
}

func (k *KubeMQ) Close() {
	_ = k.client.Close()
}
