package controllers

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"alfredo/tabunganku/pkg/dtos"
	"alfredo/tabunganku/pkg/middleware/jwt"
	"alfredo/tabunganku/pkg/services"
)

type SavingController interface {
	Router(router fiber.Router)
	CreateSaving(c *fiber.Ctx) error
	GetSavings(c *fiber.Ctx) error
}

type savingController struct {
	savingService services.SavingService
	redisService  services.RedisService
	userService   services.UserService
}

// CreateSaving godoc
// @Summary Create a new saving
// @Description Create a new saving record with optional image upload
// @Tags savings
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param name formData string true "Saving name" minlength(3) maxlength(50)
// @Param target_amount formData number true "Target amount" minimum(0.01)
// @Param currency_code formData string true "Currency code (3 characters)" minlength(3) maxlength(3)
// @Param filling_plan formData string true "Filling plan" Enums(Daily, Weekly, Monthly)
// @Param filling_nominal formData number true "Filling nominal amount" minimum(0.01)
// @Param image formData file true "Image file"
// @Success 200 {object} dtos.SuccessResponse
// @Failure 400 {object} dtos.ErrorResponseDTO
// @Failure 500 {object} dtos.ErrorResponseDTO
// @Router /savings [post]
func (s *savingController) CreateSaving(c *fiber.Ctx) error {
	var savingRequest dtos.SavingRequest
	err := c.BodyParser(&savingRequest)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dtos.ErrorResponseDTO{
			Success: false,
			Message: "Invalid request body",
			Code:    fiber.StatusBadRequest,
			Errors:  err.Error(),
		})
	}

	savingRequest.UserUUID = c.Locals("user_uuid").(string)

	// Handle file upload
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		// Create uploads directory if it doesn't exist
		if err = os.MkdirAll("./uploads", 0755); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dtos.ErrorResponseDTO{
				Success: false,
				Message: "Failed to create upload directory",
				Code:    fiber.StatusInternalServerError,
				Errors:  err.Error(),
			})
		}

		// Generate unique filename
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), uuid.New().String(), ext)
		filepath := fmt.Sprintf("./uploads/%s", filename)

		// Save the file
		if err := c.SaveFile(file, filepath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dtos.ErrorResponseDTO{
				Success: false,
				Message: "Failed to save file",
				Code:    fiber.StatusInternalServerError,
				Errors:  err.Error(),
			})
		}

		// Set the image path in the request
		savingRequest.Image = filepath
	}

	// Create saving
	savingResponse, err := s.savingService.CreateSaving(&savingRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.ErrorResponseDTO{
			Success: false,
			Message: "Failed to create saving",
			Code:    fiber.StatusInternalServerError,
			Errors:  err.Error(),
		})
	}

	return c.JSON(dtos.SuccessResponse{
		Success: true,
		Message: "Saving created successfully",
		Data:    savingResponse,
	})
}

// GetSavings godoc
// @Summary Get user savings
// @Description Get all savings records for the authenticated user
// @Tags savings
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} dtos.SuccessResponse
// @Failure 500 {object} dtos.ErrorResponseDTO
// @Router /savings [get]
func (s *savingController) GetSavings(c *fiber.Ctx) error {
	userUuid := c.Locals("user_uuid").(string)
	fmt.Println(userUuid)
	savings, err := s.savingService.GetSavings(userUuid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dtos.ErrorResponseDTO{
			Success: false,
			Message: "Failed to get savings",
			Code:    fiber.StatusInternalServerError,
			Errors:  err.Error(),
		})
	}

	return c.JSON(dtos.SuccessResponse{
		Success: true,
		Message: "Savings retrieved successfully",
		Data:    savings,
	})
}

// Router implements SavingController.
func (s *savingController) Router(router fiber.Router) {
	withMiddleware := router.Use(jwt.JwtMiddleware(s.userService, s.redisService))
	{
		withMiddleware.Post("/", s.CreateSaving)
		withMiddleware.Get("/", s.GetSavings)
	}
}

func NewSavingController(savingService services.SavingService, redisService services.RedisService, userService services.UserService) SavingController {
	return &savingController{savingService: savingService, redisService: redisService, userService: userService}
}
