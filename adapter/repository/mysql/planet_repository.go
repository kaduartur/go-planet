package mysql

import (
	"context"
	"strings"

	"gorm.io/gorm"

	"github.com/kaduartur/go-planet/core/domain"
	"github.com/kaduartur/go-planet/core/port"
)

type planetRepository struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) port.PlanetRepository {
	return &planetRepository{
		db: db,
	}
}

func (p *planetRepository) Create(ctx context.Context, order domain.Planet) error {
	planet := Planet{
		ID:        order.ID().String(),
		Name:      order.Name(),
		Climates:  strings.Join(order.Climates(), ","),
		Terrains:  strings.Join(order.Terrains(), ","),
		CreatedAt: order.CreatedAt(),
	}

	result := p.db.WithContext(ctx).
		Table("planets").
		Create(&planet)

	if result.Error != nil {
		return result.Error
	}

	films := make(Films, len(order.Films()))
	for i, f := range order.Films() {
		films[i] = Film{
			PlanetID:    order.ID().String(),
			Title:       f.Title(),
			ReleaseDate: f.ReleaseDate(),
		}
	}

	result = p.db.WithContext(ctx).
		Table("films").
		CreateInBatches(&films, len(films))

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (p *planetRepository) List(ctx context.Context, page, perPage int) (port.Planets, int, error) {
	var count int64 = 0
	offset := (page - 1) * perPage
	planets := make(Planets, 0)
	result := p.db.WithContext(ctx).
		Table("planets").
		Count(&count).
		Preload("Films").
		Offset(offset).
		Limit(perPage).
		Find(&planets)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return planets.ToPlanetsPort(), int(count), nil
}

func (p *planetRepository) Count(ctx context.Context) (int, error) {
	var count *int64
	p.db.WithContext(ctx).
		Table("planets").
		Count(count)
	return int(*count), nil
}
