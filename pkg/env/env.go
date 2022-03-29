package env

import (
	"github.com/spf13/viper"
	"log"
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

type Prometheus struct {
	Name string `mapstructure:"name"`
	Path string `mapstructure:"path"`
}

type Config struct {
	App        App        `mapstructure:"app"`
	DB         Database   `mapstructure:"database"`
	Prometheus Prometheus `mapstructure:"prometheus"`
}

func New() *Config {
	var cfg Config
	v := viper.New()
	v.SetConfigFile(fileName)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		log.Println("No env file, using environment variables.", err)
	}

	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatal("Error trying to unmarshal configuration", err)
	}

	return &cfg
}
