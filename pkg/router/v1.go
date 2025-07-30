package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	"alfredo/tabunganku/config"
	"alfredo/tabunganku/pkg/injectors"
)

func InitializeRouterV1(server *config.Application) {
	// Gunakan handler untuk path /swagger/* bukan middleware
	server.App.Get("/swagger/*", swagger.HandlerDefault) // Opsi sederhana

	// ATAU gunakan konfigurasi yang lebih spesifik
	server.App.Get("/swagger/*", swagger.New(swagger.Config{
		Title:        "Tabunganku API Documentation",
		URL:          "/swagger/doc.json", // Pastikan ini sesuai dengan path yang benar
		DeepLinking:  true,
		Layout:       "BaseLayout",
		DocExpansion: "none",
	}))
	api := server.App.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// Ping godoc
			// @Summary Health check endpoint
			// @Description Get server status
			// @Tags system
			// @Accept json
			// @Produce json
			// @Success 200 {object} map[string]string
			// @Router /ping [get]
			v1.Get("/ping", func(ctx *fiber.Ctx) error {
				return ctx.JSON(fiber.Map{
					"message": "pong",
				})
			})

			auth := v1.Group("/auth")
			{
				authController := injectors.InitializeUserController()
				authController.Router(auth)
			}

			saving := v1.Group("/savings")
			{
				savingController := injectors.InitializeSavingController()
				savingController.Router(saving)
			}

		}

	}
}
