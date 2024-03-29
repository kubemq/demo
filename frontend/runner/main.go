package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go run(ctx, cfg.ApiAddress, cfg.Interval)
	<-gracefulShutdown

}
