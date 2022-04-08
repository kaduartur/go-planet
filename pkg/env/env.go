package env

import (
	"context"
	"time"

	"github.com/kaduartur/go-planet/pkg/log"
	"github.com/spf13/viper"
)

const fileName = "config.toml"

type App struct {
	HttpAddr string `mapstructure:"http_addr"`
	Env      string `mapstructure:"env"`
}

type Database struct {
	HostName string `mapstructure:"hostname"`
	Port     string `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Swapi struct {
	HostName string        `mapstructure:"hostname"`
	Timeout  time.Duration `mapstructure:"timeout"`
}

type Config struct {
	App   App      `mapstructure:"app"`
	DB    Database `mapstructure:"database"`
	Swapi Swapi    `mapstructure:"swapi"`
}

func New() *Config {
	var cfg Config
	v := viper.New()
	v.SetConfigFile(fileName)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		log.Error(context.TODO(), "no env file, using environment variables.", err)
	}

	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatal(context.TODO(), "error trying to unmarshal configuration", err)
	}

	return &cfg
}
