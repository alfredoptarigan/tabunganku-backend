package services

import (
	"fmt"

	"github.com/go-playground/validator"

	"alfredo/tabunganku/pkg/dtos"
	"alfredo/tabunganku/pkg/helpers"
	"alfredo/tabunganku/pkg/models"
	"alfredo/tabunganku/pkg/repositories"
)

type UserService interface {
	FindUserByUuid(uuid string) (*models.User, error)
	Login(request *dtos.LoginRequest) (response dtos.LoginResponse, err error)
	Register(req *dtos.RegisterRequest) error
}

type userServiceImpl struct {
	repo       repositories.UserRepository
	jwtService JwtService
}

// FindUserByUuid implements UserService.
func (u *userServiceImpl) FindUserByUuid(uuid string) (*models.User, error) {
	return u.repo.FindUserByUuid(uuid)
}

// Login implements UserService.
func (u *userServiceImpl) Login(request *dtos.LoginRequest) (response dtos.LoginResponse, err error) {
	// Find user by email
	user, err := u.repo.FindUserByEmail(request.Email)
	if err != nil {
		return response, fmt.Errorf("user not found")
	}

	// Check if the password same with the hash password
	if ok, err := helpers.CheckPasswordHashWithArgon2(request.Password, user.Password); !ok || err != nil {
		return response, fmt.Errorf("invalid email or password")
	}

	generateToken := helpers.GenerateToken(32)
	token, err := u.jwtService.GenerateToken(user.UUID, generateToken)
	if err != nil {
		return response, fmt.Errorf("failed to generate token")
	}

	return dtos.LoginResponse{
		TokenType:    "Bearer",
		ExpiresIn:    int64(token.ExpiresIn),
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Email:        user.Email,
		UserUuid:     user.UUID,
		Name:         user.Name,
	}, nil
}

// Register implements UserService.
func (u *userServiceImpl) Register(req *dtos.RegisterRequest) error {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return fmt.Errorf("invalid request")
	}

	if err := u.repo.Register(req); err != nil {
		return err
	}

	return nil
}

func NewUserService(repo repositories.UserRepository, jwtsService JwtService) UserService {
	return &userServiceImpl{repo: repo, jwtService: jwtsService}
}
