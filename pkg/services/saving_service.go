package services

import (
	"alfredo/tabunganku/pkg/dtos"
	"alfredo/tabunganku/pkg/repositories"
)

type SavingService interface {
	CreateSaving(saving *dtos.SavingRequest) (response *dtos.SavingResponse, err error)
	GetSavings(userUuid string) (response []*dtos.SavingResponse, err error)
}

type savingServiceImpl struct {
	savingRepository repositories.SavingRepository
}

// CreateSaving implements SavingService.
func (s *savingServiceImpl) CreateSaving(saving *dtos.SavingRequest) (response *dtos.SavingResponse, err error) {
	return s.savingRepository.CreateSaving(saving)
}

// GetSavings implements SavingService.
func (s *savingServiceImpl) GetSavings(userUuid string) (response []*dtos.SavingResponse, err error) {
	return s.savingRepository.GetSavings(userUuid)
}

func NewSavingService(savingRepository repositories.SavingRepository) SavingService {
	return &savingServiceImpl{savingRepository: savingRepository}
}
