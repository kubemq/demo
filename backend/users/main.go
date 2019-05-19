package main

import (
	"bytes"

	"encoding/json"
	"time"

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
	_, err = NewPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}

	//err = pq.Register(context.Background(), &NewUser{
	//	Name:     "lior",
	//	Password: "password",
	//	Email:    "lior.nabat@gmail.com",
	//})
	//log.Println(err)
	//err = pq.Register(context.Background(), &NewUser{
	//	Name:     "lior",
	//	Password: "password",
	//	Email:    "lior.nabat@gmail.com",
	//})
	//log.Println(err)

	//
	//kube, err := NewKubeMQClient(cfg.KubeMQHost, cfg.KubeMQPort, cfg.LogsChannel)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//commandsCh := make(chan *kubemq.CommandReceive, 1)
	//queriesCh := make(chan *kubemq.QueryReceive, 1)
	//errCh := make(chan error, 1)
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	//log.Println("Wait for kubemq to be ready")
	//for {
	//	err = kube.StartListenToCommands(ctx, cfg.UsersChannel, cfg.Group, commandsCh, errCh)
	//	if err != nil {
	//		log.Printf("error connecting to kubemq, error: %s, retrying...\n", err.Error())
	//		time.Sleep(time.Second)
	//	} else {
	//
	//		break
	//
	//	}
	//}
	//for {
	//	err = kube.StartListenToQueries(ctx, cfg.UsersChannel, cfg.Group, queriesCh, errCh)
	//	if err != nil {
	//		log.Printf("error connecting to kubemq, error: %s, retrying...\n", err.Error())
	//		time.Sleep(time.Second)
	//	} else {
	//
	//		break
	//
	//	}
	//}
	//log.Println("waiting for commands / queries from KubeMQ")
	//for {
	//	select {
	//	case command := <-commandsCh:
	//		log.Println(command)
	//	case query := <-queriesCh:
	//		log.Println(query)
	//
	//	case err := <-errCh:
	//		log.Fatal(err)
	//	case <-gracefulShutdown:
	//		kube.Close()
	//		return
	//	}
	//}
}
