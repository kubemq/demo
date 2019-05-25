package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	kubemq "github.com/kubemq-io/kubemq-go"
)

func PrettyJson(data interface{}) string {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")

	err := encoder.Encode(data)
	if err != nil {
		return ""
	}
	return buffer.String()
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
	s := NewSlack(cfg.SlackToken)
	kube, err := NewKubeMQClient(cfg.KubeMQHost, cfg.KubeMQPort)
	if err != nil {

		log.Fatal(err)
	}
	eventsCh := make(chan *kubemq.EventStoreReceive, 1)
	errCh := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log.Println("Wait for kubemq to be ready")
	for {
		err = kube.StartListen(ctx, cfg.Channel, cfg.Group, eventsCh, errCh)
		if err != nil {
			log.Printf("error connecting to kubemq, error: %s, retrying...\n", err.Error())
			time.Sleep(time.Second)
		} else {

			break

		}
	}

	log.Println("waiting for events from KubeMQ")
	for {
		select {
		case event := <-eventsCh:
			_ = s.SendMessage(cfg.SlackChannel, event.Metadata, fmt.Sprintf("%s", event.Body))

		case err := <-errCh:
			log.Fatal(err)
		case <-gracefulShutdown:
			kube.Close()
			return
		}
	}
}
