package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	orderhttp "github.com/kaduartur/go-planet/adapter/http"
	pmiddleware "github.com/kaduartur/go-planet/adapter/http/middleware"
	"github.com/kaduartur/go-planet/adapter/http/server"
	"github.com/kaduartur/go-planet/adapter/integration/swapi"
	"github.com/kaduartur/go-planet/adapter/repository/mysql"
	"github.com/kaduartur/go-planet/core/service"
	"github.com/kaduartur/go-planet/pkg/env"
	"github.com/kaduartur/go-planet/pkg/log"
)

func main() {
	cfg := env.New()
	if cfg.App.Env == "prod" {
	}

	router := chi.NewRouter()
	router.Use(pmiddleware.Logger(log.ZapLogger, true))
	router.Use(middleware.Heartbeat("/health"))
	router.Use(middleware.Recoverer)
	router.Use(pmiddleware.RequestID)
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	db, err := mysql.NewConnection(cfg.DB)
	if err != nil {
		log.Fatal(context.TODO(), "error to setup database", err)
	}

	// Clients
	swapiClient := swapi.NewClient(cfg.Swapi.HostName, cfg.Swapi.Timeout)

	// Repositories
	orderRepo := mysql.NewOrderRepo(db)

	// Services
	orderService := service.NewPlanetService(orderRepo, swapiClient)

	orderhttp.MakeHandler(router, orderService)

	server.New().
		Address(cfg.App.HttpAddr).
		Routes(router).
		Run()
}
