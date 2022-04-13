package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
	"os"
)

type ConfigV1 struct {
	Application *ApplicationConfig `yaml:"application"`
}

type ApplicationConfig struct {
	Environment *string `yaml:"environment"`
	Name        *string `yaml:"name"`
	URL         *string `yaml:"url"`
	Debug       *bool   `yaml:"debug"`
	Cloud       *string `yaml:"cloud"`
}

var cfg ConfigV1

func Get() ConfigV1 {
	return cfg
}

func Provide() {
	env := os.Getenv("ENV")
	err := cleanenv.ReadConfig("files/etc/"+env+".config.yml", &cfg)
	if err != nil {
		log.Info().Msg("read config from local")
		err := cleanenv.ReadConfig("/etc/"+env+".config.yml", &cfg)
		if err != nil {
			log.Panic().Err(err).Msg("fail to create config from file")
		}
	}
	log.Info().Msg("config initializing from file done")
}
