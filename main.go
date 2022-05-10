package main

import (
	"time"

	"github.com/DeniesKresna/jhapi2/config"
	"github.com/DeniesKresna/jhapi2/router"
	"github.com/DeniesKresna/jhapi2/utils"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func init() {
	setTimezone()
	err := godotenv.Load()
	if err != nil {
		log.Error().Err(err).Msg("unable to load env through env config")
	}
	config.Provide()
	utils.InitZeroLog()
	config.ProvideDB()
	config.ProvideCache()
	router.Provide()
}

func setTimezone() {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Error().Err(err).Msg("fail to set timezone to Asia/Jakarta")
	}
	time.Local = loc
}

func main() {

}
