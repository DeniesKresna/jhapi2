package router

import (
	"github.com/DeniesKresna/jhapi2/config"
	"github.com/DeniesKresna/jhapi2/repository/cache/redis"
	userSql "github.com/DeniesKresna/jhapi2/repository/sql/user"
	userService "github.com/DeniesKresna/jhapi2/service/user"
	"github.com/DeniesKresna/jhapi2/utils"
)

type Object struct {
}

func Provide() {
	utility := utils.Provide()

	// repositories
	cacheRepo := redis.NewRepository(config.GetRedisClient().Cache)
	userRepo := userSql.Provide(config.GetGormClient().DB, &cacheRepo, utility)
	userSvc := userService.Provide(&userRepo, utility)
}
