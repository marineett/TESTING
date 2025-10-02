package service_logic

import (
	"data_base_project/data_base"
	"data_base_project/types"
)

type IPersonalDataService interface {
	GetPersonalData(id int64) (*types.ServicePersonalData, error)
}

type PersonalDataService struct {
	personalDataRepository data_base.IPersonalDataRepository
}

func CreatePersonalDataService(personalDataRepository data_base.IPersonalDataRepository) *PersonalDataService {
	return &PersonalDataService{
		personalDataRepository: personalDataRepository,
	}
}

func (s *PersonalDataService) GetPersonalData(id int64) (*types.ServicePersonalData, error) {
	personalData, err := s.personalDataRepository.GetPersonalData(id)
	if err != nil {
		return nil, err
	}
	return types.MapperPersonalDataDBToService(personalData), nil
}
