package server

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/kaduartur/go-planet/pkg/log"
)

type App struct {
	server *http.Server
}

func New() *App {
	s := &http.Server{
		Addr:        "0.0.0.0:8080",
		IdleTimeout: 5 * time.Second,
	}

	return &App{server: s}
}

func (a *App) Address(addr string) *App {
	a.server.Addr = addr
	return a
}

func (a *App) Routes(routes http.Handler) *App {
	a.server.Handler = routes
	return a
}

func (a *App) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	go func() {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error(context.TODO(), "http server shutdown", err)
		}
	}()
	log.Info(context.TODO(), "server HTTP started at: "+a.server.Addr)

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Info(context.TODO(), "shutting down http server")
	if err := a.server.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error(context.TODO(), "http server shutdown", err)
	}
}
