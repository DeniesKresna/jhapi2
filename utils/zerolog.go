package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	stdLog "log"

	"github.com/DeniesKresna/jhapi2/config"
	zlogSentry "github.com/archdx/zerolog-sentry"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitZeroLog() {
	cw := zerolog.ConsoleWriter{
		Out: os.Stdout,
	}

	cw.FormatCaller = func(i interface{}) string {
		var c string
		if cc, ok := i.(string); ok {
			c = cc
		}
		if len(c) > 0 {
			if cwd, err := os.Getwd(); err == nil {
				if rel, err := filepath.Rel(cwd, c); err == nil {
					c = rel
				}
			}
			c = fmt.Sprintf("\x1b[34m%v:\x1b[0m", c)
		}
		return c
	}

	file := &lumberjack.Logger{
		Filename:   "/var/log/bagiklanapis.log",
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     1,
		Compress:   true,
	}
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
	zerolog.CallerFieldName = "source"

	if os.Getenv("ENV") == "dev" {
		log.Logger = zerolog.New(cw).With().Caller().Timestamp().Logger()
	} else {
		sentryWriter, err := zlogSentry.New(*config.Get().Sentry.DSN, zlogSentry.WithLevels(zerolog.FatalLevel, zerolog.PanicLevel))
		if err != nil {
			stdLog.Panic("sentry.Init:", err)
		}
		defer sentryWriter.Close()

		multi := zerolog.MultiLevelWriter(os.Stdout, file, sentryWriter)
		log.Logger = zerolog.New(multi).With().IPPrefix("host", GetIP()).Timestamp().Caller().Logger().Level(zerolog.InfoLevel)
	}
}

func GetIP() net.IPNet {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return net.IPNet{}
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return *ipnet
			}
		}
	}
	return net.IPNet{}
}

type NRMessage struct {
	Message interface{} `json:"message"`
}

func SendLogNR(msg NRMessage) error {
	byteMsg, err := json.Marshal(msg)
	if err != nil {
		log.Err(err).Msg("fail to marshal body request newrelic")
		return err
	}
	req, err := http.NewRequest(http.MethodPost, *config.Get().NewRelic.LogAPI, bytes.NewBuffer(byteMsg))
	if err != nil {
		log.Err(err).Msg("fail to create request newrelic")
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Api-Key", *config.Get().NewRelic.License)

	client := &http.Client{Timeout: 30 * time.Second}
	_, err = client.Do(req)
	if err != nil {
		log.Err(err).Msg("fail to store log to newrelic")
		return err
	}

	return nil
}
