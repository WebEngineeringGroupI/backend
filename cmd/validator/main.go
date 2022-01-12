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

	validator := factory.NewValidator(ctx)
	log.Println("starting url validator")
	go func() {
		err := validator.Start(ctx)
		if err != nil {
			log.Printf("unable to start validator: %s", err)
		}
	}()

	log.Println("url validator started")

	<-ctx.Done()
	log.Println("detected signal interruption, exiting...")
	time.Sleep(5 * time.Second)
}

func gracefulShutdownOnSignal() context.Context {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	return ctx
}
