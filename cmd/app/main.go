package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/rizqdwan/go-mono-project/config"
)

func main() {
	if err := config.LoadConfig(); err != nil{
		log.Fatalf("Failed to load configuration: %v", err)
	}

	port := strconv.Itoa(config.Cfg.App.Port)
	log.Printf("Starting server on %s", port)

	e := echo.New()

	srv := &http.Server{
		Addr: ":" + port,
		Handler: e,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}