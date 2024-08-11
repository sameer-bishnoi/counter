package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sameer-bishnoi/counter/app/infra/algo"
	"github.com/sameer-bishnoi/counter/app/infra/scheduler"
	"github.com/sameer-bishnoi/counter/app/infra/storage"
	"github.com/sameer-bishnoi/counter/app/presentation/http/handler"
	"github.com/sameer-bishnoi/counter/app/usecase"
)

const filePath = "./resources/tmp/storage.txt"

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	queue := algo.NewQueue()

	fileStorage := storage.NewFile()
	fileStorage.WithFilePath(filePath)
	fileStorage.WithQueue(queue)

	// Load file data in queue.
	if err := fileStorage.Load(); err != nil {
		log.Printf("unable to load data from file: %v", err)
		os.Exit(1)
	}

	healthCheckHandler := handler.NewHealthCheck()
	healthCheckHandler.WithMessage("OK")

	counterService := usecase.NewRequestCounterService()
	counterService.WithQueue(queue)

	counterHandler := handler.NewRequestCounter()
	counterHandler.WithRequestCounterService(counterService)

	// EvictionJob will delete all the expired requests counter from the queue.
	evictionJob := scheduler.NewJob()
	evictionJob.WithQueue(queue)
	evictionJob.WithInterval(100 * time.Millisecond)

	// BackGround Job to remove expired timestamps from the queue.
	done := make(chan struct{})
	go evictionJob.BackgroundJob(done)

	// Define the endpoints.
	mux := http.NewServeMux()
	mux.HandleFunc("/health-check", healthCheckHandler.CheckServerHealth)
	mux.HandleFunc("/counter", counterHandler.GetCounter)

	server := initHTTPServer(mux)

	errChan := make(chan error)
	go func() {
		runHttpServer(server, errChan)
	}()

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, syscall.SIGINT, syscall.SIGTERM)
	select {
	case err := <-errChan:
		log.Fatalf("unable to start the server: %v", err)
	case <-stopSignal:
		gracefulShutdown(ctx, server, fileStorage, done)
	}
	log.Println("Application Stopped. Bye Bye !")
}

func initHTTPServer(routes http.Handler) *http.Server {
	server := &http.Server{
		Addr:         ":8080",
		Handler:      routes,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	return server
}

func runHttpServer(server *http.Server, errChan chan<- error) {
	log.Println("Starting the server:", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}
}

func gracefulShutdown(ctx context.Context, server *http.Server, fileStorage *storage.File, done chan struct{}) {
	if err := fileStorage.Store(); err != nil {
		log.Printf("unable to save data in file: %v", err)
	}
	close(done)
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("cannot stop server gracefully: %v", err)
	}
	log.Println("server stopped gracefully")
}
