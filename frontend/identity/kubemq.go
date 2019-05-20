package main

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/kubemq-io/kubemq-go"
	"time"
)

type KubeMQ struct {
	client      *kubemq.Client
	logsChannel string
}

func NewKubeMQClient(host string, port int, logsChannel string) (*KubeMQ, error) {
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
func (k *KubeMQ) SendLog(ctx context.Context, source, entry string) error {
	return k.client.E().
		SetId(uuid.New().String()).
		SetMetadata(source).
		SetBody([]byte(entry)).
		SetChannel(k.logsChannel).
		Send(ctx)
}

func (k *KubeMQ) SendCommand(ctx context.Context, channel, metadata string, data interface{}) (*kubemq.CommandResponse, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return k.client.C().
		SetId(uuid.New().String()).
		SetMetadata(metadata).
		SetBody(body).
		SetChannel(channel).
		SetTimeout(30 * time.Second).
		Send(ctx)
}

func (k *KubeMQ) SendQuery(ctx context.Context, channel, metadata string, data interface{}) (*kubemq.QueryResponse, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return k.client.Q().
		SetId(uuid.New().String()).
		SetMetadata(metadata).
		SetBody([]byte(body)).
		SetChannel(channel).
		SetTimeout(30 * time.Second).
		Send(ctx)
}

func (k *KubeMQ) Close() {
	_ = k.client.Close()
}
