package mysql

import (
	"fmt"
	"strings"
	"time"

	"github.com/kaduartur/go-planet/core/port"
	"github.com/kaduartur/go-planet/pkg/env"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Planet struct {
	ID        string    `gorm:"primary_key" json:"id"`
	Name      string    `gorm:"type:string;not null" json:"name"`
	Climates  string    `gorm:"type:string;not null" json:"climates"`
	Terrains  string    `gorm:"type:string;not null" json:"terrains"`
	Films     Films     `gorm:"foreignKey:PlanetID" json:"films"`
	CreatedAt time.Time `gorm:"type:time;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:time" json:"updated_at"`
}

type Film struct {
	ID          int       `gorm:"primary_key" json:"id"`
	PlanetID    string    `gorm:"type:string;not null;index" json:"planet_id"`
	Title       string    `gorm:"type:string;not null" json:"title"`
	ReleaseDate time.Time `gorm:"type:time;not null" json:"release_date"`
}

type Films []Film

func (ff Films) ToFilmsPort() port.FilmsDTO {
	films := make(port.FilmsDTO, len(ff))
	for i, f := range ff {
		films[i] = port.FilmDTO{
			Title:       f.Title,
			ReleaseDate: f.ReleaseDate,
		}
	}

	return films
}

type Planets []Planet

func (pp Planets) ToPlanetsPort() port.Planets {
	planets := make(port.Planets, len(pp))
	for i, p := range pp {
		planets[i] = port.Planet{
			Id:        p.ID,
			Name:      p.Name,
			Climates:  p.ClimatesToArr(),
			Terrains:  p.TerrainsToArr(),
			Films:     p.Films.ToFilmsPort(),
			CreatedAt: p.CreatedAt,
		}
	}

	return planets
}

func (p Planet) ClimatesToArr() []string {
	return strings.Split(p.Climates, ",")
}

func (p Planet) TerrainsToArr() []string {
	return strings.Split(p.Terrains, ",")
}

func NewConnection(cfg env.Database) (*gorm.DB, error) {
	const format = "%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True"
	credentials := fmt.Sprintf(
		format,
		cfg.Username,
		cfg.Password,
		cfg.HostName,
		cfg.Port,
		cfg.Name,
	)
	db, err := gorm.Open(mysql.Open(credentials))
	if err != nil {
		return nil, err
	}

	sqlDB, _ := db.DB()

	sqlDB.SetConnMaxLifetime(time.Second * 10)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	doMigration(db)

	return db, nil
}

func doMigration(db *gorm.DB) {
	if !db.Migrator().HasTable(&Planet{}) {
		err := db.Migrator().CreateTable(&Planet{})
		if err != nil {
			panic(err)
		}
	}

	if !db.Migrator().HasTable(&Film{}) {
		err := db.Migrator().CreateTable(&Film{})
		if err != nil {
			panic(err)
		}
	}
}
