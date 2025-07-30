//go:build wireinject
// +build wireinject

package injectors

import (
	"github.com/google/wire"

	"alfredo/tabunganku/config"
	"alfredo/tabunganku/pkg/controllers"
	"alfredo/tabunganku/pkg/repositories"
	"alfredo/tabunganku/pkg/services"
	"alfredo/tabunganku/pkg/validator"
)

var initDBPostgresSet = wire.NewSet(
	config.InitDatabasePostgres,
)

var redisSet = wire.NewSet(
	config.InitRedis,
	repositories.NewRedisRepository,
	services.NewRedisService,
)

var jwtSet = wire.NewSet(
	services.NewJwtService,
)

var authSet = wire.NewSet(
	redisSet,
	initDBPostgresSet,
	services.NewUserService,
	repositories.NewUserRepository,
	validator.NewValidator,
)

func InitializeApplication() *config.Application {
	wire.Build(config.NewApplication, config.InitDatabasePostgres)
	return nil
}

func InitializeUserController() controllers.UserController {
	wire.Build(
		authSet,
		jwtSet,
		controllers.NewUserController,
	)

	return nil
}
