package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/kubemq-io/kubemq-go"
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

func (k *KubeMQ) StartListenToCommands(ctx context.Context, channel, group string, commandsCh chan *kubemq.CommandReceive, errCh chan error) error {
	commandCh, err := k.client.SubscribeToCommands(ctx, channel, group, errCh)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case command, more := <-commandCh:
				if !more {
					return
				}
				commandsCh <- command

			case <-ctx.Done():
				return
			}

		}
	}()
	return nil
}

func (k *KubeMQ) StartListenToQueries(ctx context.Context, channel, group string, queryCh chan *kubemq.QueryReceive, errCh chan error) error {
	queriesCh, err := k.client.SubscribeToQueries(ctx, channel, group, errCh)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case query, more := <-queriesCh:
				if !more {
					return
				}
				queryCh <- query

			case <-ctx.Done():
				return
			}

		}
	}()
	return nil
}

func (k *KubeMQ) SendResponse(ctx context.Context, response *kubemq.Response) error {
	return k.client.R().
		SetBody(response.Body).
		SetMetadata(response.Metadata).
		SetError(response.Err).
		SetExecutedAt(response.ExecutedAt).
		SetResponseTo(response.ResponseTo).
		SetRequestId(response.RequestId).
		Send(ctx)

}
func (k *KubeMQ) Close() {
	_ = k.client.Close()
}