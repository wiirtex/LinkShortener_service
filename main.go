package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"ozonLinkShortener/internal/handlers"
	"ozonLinkShortener/internal/memory"
	"ozonLinkShortener/internal/memory/db"
	"ozonLinkShortener/internal/memory/inMemory"
	"ozonLinkShortener/pkg/middlewares"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	if len(os.Args) < 2 || os.Args[1] == "--use-memory" {
		handlers.Memory = memory.Memory{
			Storage: inMemory.GetMemoryInstance(),
		}
		logger.Info("Using in-memory storage")
	} else if os.Args[1] == "--use-db" {
		handlers.Memory = memory.Memory{
			Storage: db.GetMemoryInstance(),
		}
		logger.Info("Using db storage")
	} else {
		logger.Fatal("Unknown argument. Valid options: --use-db or --use-memory.")
	}

	r := chi.NewRouter()
	r.Use(middlewares.AddLogger)
	r.Route("/", func(r chi.Router) {
		r.Use(middlewares.ParseInputJson)
		r.Use(middlewares.AddClientInfo)

		r.Post("/", handlers.ManageLongLink)
		r.Get("/", handlers.ManageShortLink)
	})
	r.Get("/{shortId}", handlers.Redirect)

	ctx := context.Background()

	srv := &http.Server{
		Addr:    ":15001",
		Handler: r,
	}
	defer srv.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.WithError(err).Fatal("start server")
		}
	}()

	<-stop
	logrus.Info("caught stop signal")

	ctx, cancelFunc := context.WithTimeout(ctx, 5*time.Second)

	if err := srv.Shutdown(ctx); err != nil {
		logrus.WithError(err).Fatal("Server Shutdown Failed")
	}

	cancelFunc()

	logrus.Info("Server Exited Properly")
}
