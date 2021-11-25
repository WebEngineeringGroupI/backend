package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	ctx := gracefulShutdownOnSignal()
	wg := sync.WaitGroup{}
	factory := newFactory()

	launchHTTPServer(ctx, factory, &wg)
	launchGRPCServer(ctx, factory, &wg)

	<-ctx.Done()
	log.Println("attempting graceful shutdown...")
	wg.Wait()
	log.Println("server exited properly")
}

func launchHTTPServer(ctx context.Context, factory *factory, wg *sync.WaitGroup) {
	server := &http.Server{Addr: ":8080", Handler: factory.NewHTTPRouter()}

	go func() {
		log.Println("starting listening for HTTP requests on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("error in http listener: %s", err.Error())
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		cancelCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := server.Shutdown(cancelCtx)
		if err != nil {
			log.Printf("error occurred while shutting down HTTP server: %s", err.Error())
		}
		log.Println("closed HTTP server")
	}()
}

func launchGRPCServer(ctx context.Context, factory *factory, wg *sync.WaitGroup) {
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("unable to listen for gRPC: %s", err.Error())
	}

	grpcServer := factory.NewGRPCServer()
	go func() {
		log.Println("starting listening for gRPC requests on :8081")
		if err = grpcServer.Serve(listener); err != nil {
			log.Fatalf("error in gRPC server: %s", err.Error())
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		grpcServer.GracefulStop()
		log.Println("closed gRPC server")
	}()
}

func gracefulShutdownOnSignal() context.Context {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	return ctx
}
