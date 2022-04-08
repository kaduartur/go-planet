package presenter

import (
	"time"

	"github.com/kaduartur/go-planet/core/port"
)

type Metadata struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	PageCount  int `json:"page_count"`
	TotalCount int `json:"total_count"`
}

type Pagination struct {
	Meta    Metadata    `json:"_metadata"`
	Results interface{} `json:"results"`
}

type PlanetResponse struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Climates  []string  `json:"climates"`
	Terrains  []string  `json:"terrains"`
	Films     []Film    `json:"films"`
	CreatedAt time.Time `json:"created_at"`
}

type PlanetsResponse []PlanetResponse

type Film struct {
	Title       string    `json:"title"`
	ReleaseDate time.Time `json:"release_date"`
}

func ToPlanetPresenter(p *port.Planet) PlanetResponse {
	return PlanetResponse{
		Id:        p.Id,
		Name:      p.Name,
		Climates:  p.Climates,
		Terrains:  p.Terrains,
		Films:     ToFilmsPresenter(p.Films),
		CreatedAt: p.CreatedAt,
	}
}

func ToFilmsPresenter(f port.FilmsDTO) []Film {
	films := make([]Film, len(f))
	for i, dto := range f {
		films[i] = Film{
			Title:       dto.Title,
			ReleaseDate: dto.ReleaseDate,
		}
	}

	return films
}

func ToPlanetsPresenter(pp port.Planets) PlanetsResponse {
	planets := make(PlanetsResponse, len(pp))
	for i, p := range pp {
		planets[i] = ToPlanetPresenter(&p)
	}

	return planets
}
