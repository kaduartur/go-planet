package request

import "github.com/kaduartur/go-planet/core/port"

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

type Page struct {
	Page    int
	PerPage int
}

type CreatePlanet struct {
	Name string `json:"name"`
}

func (c CreatePlanet) ToCommand() port.CreatePlanet {
	return port.CreatePlanet{
		Name: c.Name,
	}
}
