package service

import (
	"context"
	"net/http"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/kaduartur/go-planet/core/domain"
	"github.com/kaduartur/go-planet/core/port"
	"github.com/kaduartur/go-planet/pkg/log"
	"github.com/kaduartur/go-planet/pkg/web"
)

type planetService struct {
	wg    *sync.WaitGroup
	mu    *sync.Mutex
	repo  port.PlanetRepository
	swapi port.SwapiClient
}

func NewPlanetService(repo port.PlanetRepository, client port.SwapiClient) port.PlanetService {
	return &planetService{
		wg:    &sync.WaitGroup{},
		mu:    &sync.Mutex{},
		repo:  repo,
		swapi: client,
	}
}

func (ps *planetService) Create(ctx context.Context, createPlanet port.CreatePlanet) (*port.Planet, error) {
	planet := createPlanet.ToDomain()

	p, err := ps.swapi.FindPlanetByName(ctx, planet.Name())
	if err != nil {
		return nil, err
	}

	for _, url := range p.Films {
		ps.wg.Add(1)
		go ps.findFilm(ctx, url, &planet)
	}
	ps.wg.Wait()

	climates := strings.Split(p.Climate, ", ")
	planet.SetClimates(climates)

	terrains := strings.Split(p.Terrain, ", ")
	planet.SetTerrains(terrains)

	if err := ps.repo.Create(ctx, planet); err != nil {
		log.Error(ctx, "error to create a new planet", err)
		return nil, web.NewError(http.StatusInternalServerError, "error to create a new planet")
	}

	return port.ToPlanetCreated(&planet), nil
}

func (ps *planetService) List(ctx context.Context, page int, perPage int) (port.Planets, int, error) {
	planets, total, err := ps.repo.List(ctx, page, perPage)
	if err != nil {
		log.Error(ctx, "error to list planets", err)
		return nil, 0, web.NewError(http.StatusInternalServerError, "error to list planets")
	}

	return planets, total, nil
}

func (ps *planetService) findFilm(ctx context.Context, url string, planet *domain.Planet) {
	base := path.Base(url)
	film, err := ps.swapi.FindFilmByID(ctx, base)
	if err != nil {
		return
	}
	defer ps.wg.Done()

	date, _ := time.Parse("2006-01-02", film.ReleaseData)
	newFilm := *domain.NewFilm(film.Title, date)

	ps.mu.Lock()
	planet.SetFilms(append(planet.Films(), newFilm))
	ps.mu.Unlock()
}
