package port

import (
	"context"
	"time"

	"github.com/kaduartur/go-planet/core/domain"
)

type PlanetService interface {
	Create(ctx context.Context, create CreatePlanet) (*Planet, error)
	List(ctx context.Context, page int, perPage int) (Planets, int, error)
}

type PlanetRepository interface {
	Create(ctx context.Context, planet domain.Planet) error
	List(ctx context.Context, page, perPage int) (Planets, int, error)
	Count(ctx context.Context) (int, error)
}

type CreatePlanet struct {
	Name string
}

func (c CreatePlanet) ToDomain() domain.Planet {
	return domain.NewPlanet(c.Name)
}

type Planet struct {
	Id        string
	Name      string
	Climates  []string
	Terrains  []string
	Films     FilmsDTO
	CreatedAt time.Time
}

type Planets []Planet

type FilmDTO struct {
	Title       string
	ReleaseDate time.Time
}

type FilmsDTO []FilmDTO

func ToPlanetCreated(p *domain.Planet) *Planet {
	return &Planet{
		Id:        p.ID().String(),
		Name:      p.Name(),
		Climates:  p.Climates(),
		Terrains:  p.Terrains(),
		Films:     ToFilmsPresenter(p.Films()),
		CreatedAt: p.CreatedAt(),
	}
}

func ToFilmsPresenter(ff domain.Films) FilmsDTO {
	films := make(FilmsDTO, len(ff))
	for i, f := range ff {
		films[i] = FilmDTO{
			Title:       f.Title(),
			ReleaseDate: f.ReleaseDate(),
		}
	}

	return films
}
