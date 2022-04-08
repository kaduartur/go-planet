package port

import "context"

type SwapiClient interface {
	FindPlanetByName(ctx context.Context, name string) (*SwapiPlanet, error)
	FindFilmByID(ctx context.Context, id string) (*SwapiFilm, error)
}

type SwapiResponse struct {
	Planets SwapiPlanets `json:"results"`
}

type SwapiPlanets []SwapiPlanet

func (pp SwapiPlanets) First() *SwapiPlanet {
	if len(pp) <= 0 {
		return nil
	}
	return &pp[0]
}

type SwapiPlanet struct {
	Name    string   `json:"name"`
	Films   []string `json:"films"`
	Climate string   `json:"climate"`
	Terrain string   `json:"terrain"`
}

type SwapiFilm struct {
	Title       string `json:"title"`
	ReleaseData string `json:"release_date"`
}
