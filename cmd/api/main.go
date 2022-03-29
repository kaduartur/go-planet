package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	orderhttp "github.com/kaduartur/go-planet/adapter/http"
	"github.com/kaduartur/go-planet/adapter/http/middleware"
	"github.com/kaduartur/go-planet/adapter/repository/mysql"
	"github.com/kaduartur/go-planet/core/service"
	"github.com/kaduartur/go-planet/pkg/env"
	"github.com/kaduartur/go-planet/pkg/log"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	cfg := env.New()
	logger := log.ZapLogger(cfg)
	if cfg.App.Env == "prod" {
	}

	/*router := gin.New()

	// prometheus middleware
	prom := ginprom.New(
		ginprom.Engine(router),
		ginprom.Subsystem(cfg.Prometheus.Name),
		ginprom.Path(cfg.Prometheus.Path),
	)
	router.Use(prom.Instrument())

	// application handlers
	router.Use(middleware.GinHandler("/health"))

	// logger middleware
	router.Use(ginzap.Ginzap(log.LoggerZap, time.RFC3339, true))

	router.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	})*/

	router := chi.NewRouter()
	router.Use(middleware.Logger(log.LoggerZap, true))
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	db, err := mysql.NewConnection(cfg.DB)
	if err != nil {
		logger.Fatal(err)
	}

	// Repositories
	orderRepo := mysql.NewOrderRepo(logger, db)

	// Services
	orderService := service.NewOrderService(orderRepo)

	orderhttp.MakeHandler(router, orderService)

	server := &http.Server{
		Handler: router,
		Addr:    cfg.App.HttpAddr,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http server shutdown", err)
		}
	}()
	logger.Info("Server HTTP started at:", cfg.App.HttpAddr)

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Info("Shutting down http server")
	if err := server.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("http server shutdown", err)
	}
}
