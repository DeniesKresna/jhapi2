package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
)

type ConfigV1 struct {
	Application *ApplicationConfig   `yaml:"application"`
	Sentry      *SentryConfig        `yaml:"sentry"`
	NewRelic    *NewRelicConfig      `yaml:"newrelic"`
	Database    *DatabaseConfig      `yaml:"database"`
	Cache       *CacheConfig         `yaml:"cache"`
	AWS         *AWSConfig           `yaml:"aws"`
	Constant    *ApplicationConstant `yaml:"constant"`
}

type ApplicationConfig struct {
	Environment *string `yaml:"environment"`
	Name        *string `yaml:"name"`
	URL         *string `yaml:"url"`
	Debug       *bool   `yaml:"debug"`
	Cloud       *string `yaml:"cloud"`
}

type NewRelicConfig struct {
	Enable  *string `yaml:"enable"`
	License *string `yaml:"license"`
	LogAPI  *string `yaml:"log_api"`
}

type SentryConfig struct {
	DSN *string `json:"dsn"`
}

type DatabaseConfig struct {
	DSN *string `yaml:"dsn"`
}

type CacheConfig struct {
	Cache *Redis `yml:"redis"`
}

type Redis struct {
	Password string `yaml:"password" env:"password"`
	Host     string `yaml:"host" env:"host"`
}

type AWSConfig struct {
	PrivateKey *string `yaml:"private_key"`
	PublicKey  *string `yaml:"public_key"`
	Bucket     *string `yaml:"bucket"`
	Prefix     *string `yaml:"prefix"`
	Region     *string `yaml:"region"`
}

type ApplicationConstant struct {
	RequestTimeout time.Duration `yaml:"request_timeout"`
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
