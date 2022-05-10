package config

import (
	"errors"
	"time"

	"github.com/DeniesKresna/jhapi2/types"
	"github.com/go-redis/redis"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// SQLDB is the concrete implementation of sqlGdbc by using *sql.DB
type GormClient struct {
	DB *gorm.DB
}

type RedisClient struct {
	Cache *redis.Client
}

var gormClient *GormClient
var redisClient *RedisClient

func ProvideDB() (*GormClient, error) {
	dsn := *Get().Database.DSN

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic().Err(err).Msg("fail to initialize database")
		return nil, err
	}

	db.AutoMigrate(&types.User{})
	//seems this line code cannot be executed
	//raise error Error 1227: Access denied; you need (at least one of) the SUPER privilege(s) for this operation
	//db.Exec("SET GLOBAL sql_mode=(SELECT REPLACE(@@sql_mode,'ONLY_FULL_GROUP_BY',''));")

	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Info().Msg("database connected")
	gormClient = &GormClient{DB: db}
	if *Get().Application.Debug {
		log.Debug().Msg("enable debug db")
		gormClient.DB = gormClient.DB.Debug()
	}
	return gormClient, nil
}

func GetGormClient() *GormClient {
	return gormClient
}

func ProvideCache() (*redis.Client, error) {
	redisDB := *Get().Cache.Cache

	client := redis.NewClient(&redis.Options{
		Addr:     redisDB.Host,
		Password: redisDB.Password,
		DB:       0,
	})
	if client == nil {
		err := errors.New("redis db not found")
		log.Error().Err(err).Msg("redis is nil")
		return nil, err
	}

	redisClient = &RedisClient{Cache: client}
	if err := ping(redisClient.Cache); err != nil {
		log.Error().Err(err).Msg("redis connection not ok")
		return nil, err
	}
	log.Info().Msg("redis connected")
	return client, nil
}

func GetRedisClient() *RedisClient {
	return redisClient
}

func ping(client *redis.Client) error {
	_, err := client.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}
