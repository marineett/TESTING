package service_logic

import (
	"data_base_project/data_base"
	"data_base_project/types"
)

type IAuthService interface {
	Authorize(auth_data types.ServiceAuthData) (types.ServiceAuthVerdict, error)
	CheckLogin(username string) (bool, error)
}

type AuthService struct {
	AuthRepository data_base.IAuthRepository
}

func CreateAuthService(authRepository data_base.IAuthRepository) *AuthService {
	return &AuthService{
		AuthRepository: authRepository,
	}
}

func (s *AuthService) Authorize(auth_data types.ServiceAuthData) (types.ServiceAuthVerdict, error) {
	authVerdict, err := s.AuthRepository.Authorize(*types.MapperAuthDataServiceToDB(&auth_data))
	if err != nil {
		return types.ServiceAuthVerdict{}, err
	}
	return *types.MapperAuthVerdictDBToService(&authVerdict), nil
}

func (s *AuthService) CheckLogin(username string) (bool, error) {
	loginExists, err := s.AuthRepository.CheckLogin(username)
	if err != nil {
		return false, err
	}
	return loginExists, nil
}
