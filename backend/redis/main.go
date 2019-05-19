package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"time"

	"github.com/kubemq-io/kubemq-go"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type RedisLog struct {
	ID       string `json:"id"`
	CreateAt time.Time
	Command  string `json:"command"`
	Key      string `json:"key"`
	Value    string `json:"value"`
	Result   string `json:"result"`
}

func (rl *RedisLog) String() string {
	data, _ := json.Marshal(rl)
	var buff *bytes.Buffer
	_ = json.Indent(buff, data, "", "\t")

	return buff.String()
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
	var redis *Redis
	log.Println("Wait for redis to be ready")
	for {
		redis, err = NewRedisClient(cfg.RedisAddress)
		if err != nil {
			log.Printf("error connecting to redis, error: %s retrying...\n", err.Error())
		} else {
			break
		}
	}

	kube, err := NewKubeMQClient(cfg.KubeMQHost, cfg.KubeMQPort, cfg.LogsChannel)
	if err != nil {
		log.Fatal(err)
	}
	commandsCh := make(chan *kubemq.CommandReceive, 1)
	queriesCh := make(chan *kubemq.QueryReceive, 1)
	errCh := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log.Println("Wait for kubemq to be ready")
	for {
		err = kube.StartListenToCommands(ctx, cfg.Channel, cfg.Group, commandsCh, errCh)
		if err != nil {
			log.Printf("error connecting to kubemq, error: %s, retrying...\n", err.Error())
			time.Sleep(time.Second)
		} else {

			break

		}
	}
	for {
		err = kube.StartListenToQueries(ctx, cfg.Channel, cfg.Group, queriesCh, errCh)
		if err != nil {
			log.Printf("error connecting to kubemq, error: %s, retrying...\n", err.Error())
			time.Sleep(time.Second)
		} else {

			break

		}
	}
	log.Println("waiting for commands / queries from KubeMQ")
	for {
		select {
		case command := <-commandsCh:
			log.Println(fmt.Sprintf("redis Set command received - Key: %s, Value: %s", command.Metadata, string(command.Body)))
			err := redis.Set(command.Metadata, command.Body)
			resp := &kubemq.Response{
				RequestId:  command.Id,
				ResponseTo: command.ResponseTo,
				Metadata:   command.Metadata,
				Body:       nil,
			}
			if err != nil {
				log.Printf("error on sending command to redis: %s\n", err.Error())
				resp.Err = err

			} else {
				resp.ExecutedAt = time.Now()
			}
			err = kube.SendResponse(ctx, resp)
			if err != nil {
				log.Printf("error on sending response from redis: %s\n", err.Error())

			}
			rl := &RedisLog{
				ID:       uuid.New().String(),
				CreateAt: time.Now(),
				Command:  "Set",
				Key:      command.Metadata,
				Value:    string(command.Body),
				Result:   "",
			}
			if resp.Err != nil {
				rl.Result = resp.Err.Error()
			} else {
				rl.Result = "ok"
			}
			err = kube.SendLog(ctx, "redis", rl.String())
			log.Println("redis Set command completed")
		case query := <-queriesCh:
			log.Println(fmt.Sprintf("redis Get command received - Key: %s", query.Metadata))
			result, err := redis.Get(query.Metadata)
			resp := &kubemq.Response{
				RequestId:  query.Id,
				ResponseTo: query.ResponseTo,
				Metadata:   query.Metadata,
			}
			if err != nil {
				log.Printf("error on sending command to redis: %s\n", err.Error())
				resp.Err = err

			} else {
				resp.ExecutedAt = time.Now()
				resp.Body = result
			}
			err = kube.SendResponse(ctx, resp)
			if err != nil {
				log.Printf("error on sending response from redis: %s\n", err.Error())

			}
			rl := &RedisLog{
				ID:       uuid.New().String(),
				CreateAt: time.Now(),
				Command:  "Get",
				Key:      query.Metadata,
				Value:    string(resp.Body),
				Result:   "",
			}
			if resp.Err != nil {
				rl.Result = resp.Err.Error()
			} else {
				rl.Result = "ok"
			}
			log.Println(fmt.Sprintf("redis Get command completed with Value: %s", string(result)))
		case err := <-errCh:
			log.Fatal(err)
		case <-gracefulShutdown:
			kube.Close()
			return
		}
	}
}
