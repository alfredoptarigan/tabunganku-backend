package repositories

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"alfredo/tabunganku/pkg/dtos"
	"alfredo/tabunganku/pkg/helpers"
	"alfredo/tabunganku/pkg/models"
)

type UserRepository interface {
	Login(email string, password string) (user *models.User, err error)
	Register(req *dtos.RegisterRequest) error
	FindUserByUuid(uuid string) (*models.User, error)
	FindUserByEmail(email string) (*models.User, error)
}

type userRepositoryImpl struct {
	db *gorm.DB
}

// FindUserByUuid implements UserRepository.
func (u *userRepositoryImpl) FindUserByUuid(uuid string) (*models.User, error) {
	var user models.User
	if err := u.db.Where("uuid = ? AND deleted_at IS NULL", uuid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		} else {
			return nil, fmt.Errorf("please try again later")
		}
	}

	return &user, nil
}

// FindUserByEmail implements UserRepository.
func (u *userRepositoryImpl) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := u.db.Where("email = ? AND deleted_at IS NULL", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		} else {
			return nil, fmt.Errorf("please try again later")
		}
	}

	return &user, nil
}

// Login implements UserRepository.
func (u *userRepositoryImpl) Login(email string, password string) (user *models.User, err error) {
	if err := u.db.Where("email = ? AND password = ? AND deleted_at IS NULL", email, password).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, fmt.Errorf("user not found")
		} else {
			// Handle other errors
			panic(err)
		}
	}

	return user, nil
}

// Register implements UserRepository.
func (u *userRepositoryImpl) Register(req *dtos.RegisterRequest) error {
	var existingUser models.User
	if err := u.db.Debug().Unscoped().Where("email = ? ", req.Email).First(&existingUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Email tidak ditemukan, boleh lanjut register
			// Lanjut ke proses create user di bawah
		} else {
			return fmt.Errorf("please try again later")
		}
	} else {
		// Email ditemukan, berarti sudah ada
		return fmt.Errorf("email already exists")
	}

	// Hash password with argon2
	hashedPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	return u.db.Transaction(func(tx *gorm.DB) error {
		user := models.User{
			Name:        req.Name,
			Email:       req.Email,
			Password:    hashedPassword,
			PhoneNumber: &req.PhoneNumber,
			Photo:       &req.Image,
		}

		if err := tx.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		return nil
	})
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}
