package main

import (
	"context"
	"fmt"

	kubemq "github.com/kubemq-io/kubemq-go"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Log struct {
	ID       string `json:"id"`
	ClientID string `json:"client_id"`
	Metadata string `json:"metadata"`
	Body     string `json:"body"`
}

func main() {
	var gracefulShutdown = make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGTERM)
	signal.Notify(gracefulShutdown, syscall.SIGINT)
	signal.Notify(gracefulShutdown, syscall.SIGKILL)
	signal.Notify(gracefulShutdown, syscall.SIGQUIT)
	cfg, err := LoadConfig()
	if err != nil {
		log.Println("error on loading config file:")
		log.Fatal(err)
	}
	el, err := NewElasticSearch(cfg.ElasticAddress)
	if err != nil {
		log.Fatal(err)
	}
	kube, err := NewKubeMQClient(cfg.KubeMQHost, cfg.KubeMQPort)
	if err != nil {
		log.Fatal(err)
	}
	eventsCh := make(chan *kubemq.Event, 1)
	errCh := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = kube.StartListen(ctx, cfg.Channel, cfg.Group, eventsCh, errCh)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("waiting for events from KubeMQ")
	for {
		select {
		case event := <-eventsCh:
			fmt.Println(event)
			err := el.Save(ctx, &Log{
				ID:       event.Id,
				ClientID: event.ClientId,
				Metadata: event.Metadata,
				Body:     string(event.Body),
			})
			if err != nil {
				log.Println(err)
			}

		case err := <-errCh:
			log.Fatal(err)
		case <-gracefulShutdown:
			kube.Close()
			return
		}
	}
}
