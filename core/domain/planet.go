package domain

import (
	"encoding/hex"
	"strings"
	"time"

	"golang.org/x/crypto/blake2b"
)

func newPlanetID(name string) PlanetID {
	hash, _ := blake2b.New256([]byte(strings.ToLower(name)))
	sum := hash.Sum(nil)
	return PlanetID{id: hex.EncodeToString(sum)}
}

type PlanetID struct {
	id string
}

func (p PlanetID) String() string {
	return p.id
}

type Planet struct {
	id        PlanetID
	name      string
	climates  Climates
	terrains  Terrains
	films     Films
	createdAt time.Time
	updatedAt *time.Time
}

type Climates []string

type Terrains []string

func (p Planet) ID() PlanetID {
	return p.id
}

func (p Planet) Name() string {
	return p.name
}

func (p Planet) Climates() []string {
	return p.climates
}

func (p *Planet) SetClimates(cc Climates) {
	p.climates = cc
}

func (p Planet) Terrains() []string {
	return p.terrains
}

func (p *Planet) SetTerrains(tt Terrains) {
	p.terrains = tt
}

func (p Planet) Films() Films {
	return p.films
}

func (p *Planet) SetFilms(ff Films) {
	p.films = ff
}

func (p *Planet) CreatedAt() time.Time {
	return p.createdAt
}

func (p *Planet) UpdateAt() *time.Time {
	return p.updatedAt
}

func NewPlanet(name string) Planet {
	return Planet{
		id:        newPlanetID(name),
		name:      name,
		climates:  make(Climates, 0),
		terrains:  make(Terrains, 0),
		films:     make(Films, 0),
		createdAt: time.Now(),
	}
}

type Films []Film

type Film struct {
	title       string
	releaseDate time.Time
}

func NewFilm(title string, releaseDate time.Time) *Film {
	return &Film{title: title, releaseDate: releaseDate}
}

func (f Film) Title() string {
	return f.title
}

func (f Film) ReleaseDate() time.Time {
	return f.releaseDate
}
