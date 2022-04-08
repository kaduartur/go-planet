package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/kaduartur/go-planet/adapter/http/request"
	"github.com/kaduartur/go-planet/core/port"
	"github.com/kaduartur/go-planet/pkg/web"
)

func MakeHandler(router *chi.Mux, planetService port.PlanetService) {
	router.Route("/api/v1", func(r chi.Router) {
		makePlanetHandler(r, planetService)
	})
}

func newError(w http.ResponseWriter, err error) {
	if err := web.ToPlanetError(err); err != nil {
		b, _ := err.JSON()
		w.WriteHeader(err.Status)
		w.Write(b)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("unknown error"))
}

const PageRequestErr = `The query params "page" or "per_page" cannot be empty.`

func Page(r *http.Request) (*request.Page, error) {
	defaultPage := &request.Page{
		Page:    1,
		PerPage: 30,
	}

	query := r.URL.Query()
	pageStr := query.Get("page")
	perPageStr := query.Get("per_page")

	if strings.TrimSpace(pageStr) == "" || strings.TrimSpace(perPageStr) == "" {
		return nil, web.NewError(http.StatusBadRequest, PageRequestErr)
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return nil, web.NewError(http.StatusBadRequest, PageRequestErr)
	}

	perPage, err := strconv.Atoi(perPageStr)
	if err != nil {
		return nil, web.NewError(http.StatusBadRequest, PageRequestErr)
	}

	if page <= 0 {
		page = 1
	}

	if perPage > 30 || perPage <= 0 {
		perPage = 30
	}

	defaultPage = &request.Page{
		Page:    page,
		PerPage: perPage,
	}

	return defaultPage, nil
}
