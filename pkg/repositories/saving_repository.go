package repositories

import (
	"gorm.io/gorm"

	"alfredo/tabunganku/pkg/dtos"
	"alfredo/tabunganku/pkg/models"
)

type SavingRepository interface {
	CreateSaving(saving *dtos.SavingRequest) (response *dtos.SavingResponse, err error)
	GetSavings(userUuid string) (response []*dtos.SavingResponse, err error)
}

type savingRepositoryImpl struct {
	db *gorm.DB
}

// CreateSaving implements SavingRepository.
func (s *savingRepositoryImpl) CreateSaving(saving *dtos.SavingRequest) (response *dtos.SavingResponse, err error) {
	// Map data dari DTO ke model
	savingModel := models.Saving{
		UserUUID:       saving.UserUUID,
		Name:           saving.Name,
		TargetAmount:   saving.TargetAmount,
		CurrencyCode:   saving.CurrencyCode,
		Image:          saving.Image,
		FillingPlan:    saving.FillingPlan,
		FillingNominal: saving.FillingNominal,
	}

	// Create saving ke database
	err = s.db.Create(&savingModel).Error
	if err != nil {
		return nil, err
	}

	// Get user data
	var userModel models.User
	err = s.db.First(&userModel, "uuid = ?", saving.UserUUID).Error
	if err != nil {
		return nil, err
	}

	// Get currency data
	var currencyModel models.Currency
	err = s.db.First(&currencyModel, "currency_code = ?", saving.CurrencyCode).Error
	if err != nil {
		return nil, err
	}

	return &dtos.SavingResponse{
		UUID: savingModel.UUID,
		User: dtos.UserResponse{
			UUID:        userModel.UUID,
			Name:        userModel.Name,
			Email:       userModel.Email,
			PhoneNumber: *userModel.PhoneNumber,
			Image:       *userModel.Photo,
		},
		Name:           savingModel.Name,
		TargetAmount:   savingModel.TargetAmount,
		CurrencyCode:   savingModel.CurrencyCode,
		CurrencyFlag:   currencyModel.CountryFlag,
		Image:          savingModel.Image,
		FillingPlan:    savingModel.FillingPlan,
		FillingNominal: savingModel.FillingNominal,
		CreatedAt:      *savingModel.CreatedAt,
		UpdatedAt:      *savingModel.UpdatedAt,
	}, nil
}

// GetSavings implements SavingRepository.
func (s *savingRepositoryImpl) GetSavings(userUuid string) (response []*dtos.SavingResponse, err error) {
	var savingModels []models.Saving
	err = s.db.Where("user_uuid = ?", userUuid).Find(&savingModels).Error
	if err != nil {
		return nil, err
	}

	var currencyModels []models.Currency
	err = s.db.Find(&currencyModels).Error
	if err != nil {
		return nil, err
	}

	var userModel models.User
	err = s.db.First(&userModel, "uuid = ?", userUuid).Error
	if err != nil {
		return nil, err
	}

	for _, savingModel := range savingModels {
		var currencyModel models.Currency
		for _, currency := range currencyModels {
			if currency.CurrencyCode == savingModel.CurrencyCode {
				currencyModel = currency
			}
		}

		response = append(response, &dtos.SavingResponse{
			UUID: savingModel.UUID,
			User: dtos.UserResponse{
				UUID:        userModel.UUID,
				Name:        userModel.Name,
				Email:       userModel.Email,
				PhoneNumber: *userModel.PhoneNumber,
				Image:       *userModel.Photo,
			},
			Name:           savingModel.Name,
			TargetAmount:   savingModel.TargetAmount,
			CurrencyCode:   savingModel.CurrencyCode,
			CurrencyFlag:   currencyModel.CountryFlag,
			Image:          savingModel.Image,
			FillingPlan:    savingModel.FillingPlan,
			FillingNominal: savingModel.FillingNominal,
			CreatedAt:      *savingModel.CreatedAt,
			UpdatedAt:      *savingModel.UpdatedAt,
		})
	}

	return response, nil
}

func NewSavingRepository(db *gorm.DB) SavingRepository {
	return &savingRepositoryImpl{db: db}
}
