package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx := gracefulShutdownOnSignal()
	factory := Factory{}

	redirector := factory.NewRedirector(ctx)
	log.Println("starting event sender")
	go redirector.Start(ctx)

	log.Println("event sender started")

	<-ctx.Done()
	log.Println("detected signal interruption, exiting...")
	time.Sleep(5 * time.Second)
}

func gracefulShutdownOnSignal() context.Context {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	return ctx
}
