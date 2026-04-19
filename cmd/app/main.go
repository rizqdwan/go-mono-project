package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/rizqdwan/go-mono-project/config"
	_ "github.com/rizqdwan/go-mono-project/docs"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	app, err := newApplication(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer app.close()

	port := strconv.Itoa(cfg.App.Port)
	log.Printf("Starting %s on port %s", cfg.App.Name, port)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: app.echo,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("Server exited cleanly")
}
