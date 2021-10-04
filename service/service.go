package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/kimbellG/packet-filter/service/internal/controller"
	"github.com/kimbellG/packet-filter/service/internal/handler"
	"github.com/kimbellG/packet-filter/service/internal/scanner"

	"github.com/joho/godotenv"
)

func Run(source string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer cancel()

		oscall := <-c
		log.Printf("system call: %v", oscall)
	}()

	StartService(ctx, source)
}

func StartService(ctx context.Context, source string) {
	if err := initConfig(); err != nil {
		log.Fatalf("Failed to init config: %v", err)
	}

	interf := os.Getenv("INTERFACE")
	module, err := scanner.InitXDP(interf, source)
	if err != nil {
		log.Fatalf("Failed to load xdp module: %v", err)
	}
	defer func() {
		if err := module.RemoveXDP(interf); err != nil {
			log.Printf("Failed to close XDP module: %v", err)
		}
	}()
	log.Printf("XDP listening on %s interface\n", interf)

	router := mux.NewRouter()
	cont := controller.NewController(scanner.NewXDPScanner(module))

	handler.RegisterCount(router, cont)

	srv := &http.Server{
		Addr:    os.Getenv("PORT"),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen http connection: %v", err)
		}
	}()

	fmt.Printf("Server start listenning on %v\n", srv.Addr)
	<-ctx.Done()
	fmt.Println("Server start gracefull shutdown")

	gracefullShutdown(srv, 5*time.Second)
}

func initConfig() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("godotenv load: %v", err)
	}

	return nil

}

func gracefullShutdown(srv *http.Server, timeout time.Duration) {
	ctxShutdown, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Gracefull shutdown is failed: %s", err)
	}

	log.Println("Server closed")
}
