package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nkhamm-spb/red_soft_test/config"
	"github.com/nkhamm-spb/red_soft_test/httpserver"
	"github.com/nkhamm-spb/red_soft_test/storage"
)

func main() {
	config, err := config.LoadConfig("config.yaml")

	if err != nil {
		log.Fatalf("Error occur on read config: %v", err)
	}

	storage, err := storage.New(context.Background(), &config.Storage)

	if err != nil {
		log.Fatalf("Error occur on init storage: %v", err)
	}

	server, err := httpserver.New(context.Background(), storage, &config.Server)
	
	if err != nil {
		log.Fatalf("Error occur on create server: %v", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.Run(); err != nil {
			log.Fatalf("failed to serve server")
		}
	}()

	log.Println("server started")

	<-done
	log.Println("stopping server")

	if err := server.Shutdown(); err != nil {
		log.Fatalf("failed to stop server: %v", err)
		return
	}

	log.Println("server stopped")
}
