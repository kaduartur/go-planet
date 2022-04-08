package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/kaduartur/go-planet/adapter/http/presenter"
	"github.com/kaduartur/go-planet/adapter/http/request"
	"github.com/kaduartur/go-planet/core/port"
	"github.com/kaduartur/go-planet/pkg/web"
)

type planetHandler struct {
	planet port.PlanetService
}

func (p *planetHandler) createPlanet(w http.ResponseWriter, r *http.Request) {
	var createPlanet request.CreatePlanet
	if err := json.NewDecoder(r.Body).Decode(&createPlanet); err != nil {
		newError(w, err)
		return
	}

	if strings.TrimSpace(createPlanet.Name) == "" {
		err := web.NewError(http.StatusBadRequest, "name cannot be empty")
		newError(w, err)
		return
	}

	cmd := createPlanet.ToCommand()
	planet, err := p.planet.Create(r.Context(), cmd)
	if err != nil {
		newError(w, err)
		return
	}

	res := presenter.ToPlanetPresenter(planet)
	bytes, err := json.Marshal(&res)
	if err != nil {
		newError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(bytes)
}

func (p *planetHandler) list(w http.ResponseWriter, r *http.Request) {
	page, err := Page(r)
	if err != nil {
		newError(w, err)
		return
	}

	planets, total, err := p.planet.List(r.Context(), page.Page, page.PerPage)
	if err != nil {
		newError(w, err)
		return
	}

	pCount := int(float64(total)/float64(page.PerPage) + 0.9)
	res := request.Pagination{
		Meta: request.Metadata{
			Page:       page.Page,
			PerPage:    page.PerPage,
			PageCount:  pCount,
			TotalCount: total,
		},
		Results: presenter.ToPlanetsPresenter(planets),
	}

	bytes, err := json.Marshal(&res)
	if err != nil {
		newError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(bytes)
}

func makePlanetHandler(router chi.Router, orderService port.PlanetService) {
	handler := &planetHandler{
		planet: orderService,
	}

	// Planet endpoints
	router.Route("/planets", func(r chi.Router) {
		r.Post("/", handler.createPlanet)
		r.Get("/", handler.list)
	})
}
