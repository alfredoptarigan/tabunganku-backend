package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	fiberCors "github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"

	"alfredo/tabunganku/pkg/log"
	"alfredo/tabunganku/pkg/middleware"
)

type Application struct {
	App    *fiber.App
	Db     *gorm.DB
	Config ViperConfig
	Logger log.Logger
}

func NewApplication(db *gorm.DB) *Application {
	config := NewViperConfig()

	// Set the default timezone to UTC
	var loggingConfig log.LoggingConfig
	if err := config.UnmarshalKey("logging", &loggingConfig); err != nil {
		panic("failed to unmarshal logging config")
	}

	// Initialize the logger with the logging configuration
	logger := log.NewMultiLogger(loggingConfig)

	return &Application{
		App:    fiber.New(FiberConfig()),
		Db:     db,
		Config: config,
		Logger: logger,
	}
}

func FiberConfig() fiber.Config {
	return fiber.Config{
		EnableTrustedProxyCheck: true,
		ProxyHeader:             "X-Forwarded-*",
		Prefork:                 false,
	}
}

// RegisterMiddlewares registers the middlewares for the Fiber application
func (a *Application) RegisterMiddlewares() {
	// Register compression middleware
	a.App.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	// Register logger middleware
	a.App.Use(middleware.Logger(a.Logger))

	// Register idempotency middleware
	a.App.Use(idempotency.New(idempotency.Config{
		Lock: idempotency.NewMemoryLock(),
	}))

	a.App.Use(a.CorsMiddleware())

	// Recover Middleware
	a.App.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, err interface{}) {
			a.Logger.Error("Panic recovered: ", err)
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"data":   "Internal Server Error",
			})
		},
	}))
}

// CloseConnectionDatabase closes the database connection
func (a *Application) CloseConnectionDatabase() {
	if sqlDB, err := a.Db.DB(); err == nil {
		sqlDB.Close()
	}

	a.Logger.Info("Database connection closed")
}

func (a *Application) CorsMiddleware() fiber.Handler {
	return fiberCors.New(fiberCors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, Lang, lang, Accept-Encoding",
	})
}
