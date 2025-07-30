package jwt

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"alfredo/tabunganku/pkg/dtos"
	"alfredo/tabunganku/pkg/services"
)

type handle struct {
	ctx         *fiber.Ctx
	claim       jwt.MapClaims
	jwtToken    string
	userService services.UserService
	Logger      log.Logger
}

func JwtMiddleware(userService services.UserService, redisService services.RedisService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		jwtToken := c.Get("Authorization")
		if jwtToken == "" || jwtToken == "Bearer " || jwtToken == "Bearer" {
			log.Println("Error while validating token", "error", "Unauthorized")
			return c.Status(fiber.StatusUnauthorized).JSON(dtos.ErrorResponseDTO{
				Message: "Unauthorized",
				Code:    fiber.StatusUnauthorized,
			})
		}

		jwtToken = jwtToken[7:]

		//validate token with jwt service
		jwtService := services.NewJwtService(redisService)

		//check if token is revoked
		isRevoked := jwtService.IsTokenRevoked(jwtToken)
		if isRevoked {
			log.Println("Error while validating token", "error", "Token Not Valid")
			return c.Status(fiber.StatusUnauthorized).JSON(dtos.ErrorResponseDTO{
				Message: "Token Not Valid",
				Code:    fiber.StatusUnauthorized,
			})
		}

		var token *jwt.Token

		// Check the token using key from golang
		token, err := jwtService.ValidateToken(jwtToken)
		if err != nil {
			log.Println("Error while validating token", "error", err)
		}

		isExpired := jwtService.IsTokenExpired(jwtToken) //check if token is expired

		if isExpired {
			// If the token is expired, then revoke the token
			return c.Status(fiber.StatusUnauthorized).JSON(dtos.ErrorResponseDTO{
				Message: "Token Expired",
				Code:    fiber.StatusUnauthorized,
			})
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(dtos.ErrorResponseDTO{
				Message: "Unauthorized",
				Code:    fiber.StatusUnauthorized,
			})
		}

		return handleToken(&handle{
			ctx:         c,
			claim:       claim,
			userService: userService,
			jwtToken:    jwtToken,
		})
	}
}

func handleToken(data *handle) error {
	// If it from golang then you must check the type access
	switch data.claim["type"].(string) {
	case "access":
		break
	case "refresh":
		return data.ctx.Status(fiber.StatusUnauthorized).JSON(dtos.ErrorResponseDTO{
			Message: "Unauthorized",
			Code:    fiber.StatusUnauthorized,
		})
	case "webview":
		break
	default:
		return data.ctx.Status(fiber.StatusUnauthorized).JSON(dtos.ErrorResponseDTO{
			Message: "Unauthorized",
			Code:    fiber.StatusUnauthorized,
		})
	}

	userId := data.claim["user_id"].(string)
	userData, err := data.userService.FindUserByUuid(userId)
	if err != nil {
		return data.ctx.Status(fiber.StatusUnauthorized).JSON(dtos.ErrorResponseDTO{
			Message: "Unauthorized",
			Code:    fiber.StatusUnauthorized,
		})
	}

	data.ctx.Locals("user", userData)
	data.ctx.Locals("email", userData.Email)
	data.ctx.Locals("user_uuid", userData.UUID)
	data.ctx.Locals("token", data.jwtToken)

	return data.ctx.Next()
}
