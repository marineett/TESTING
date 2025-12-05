package service_logic

import (
	"data_base_project/data_base"
	"data_base_project/types"
)

type IAuthService interface {
	Authorize(auth_data types.ServiceAuthData) (types.ServiceAuthVerdict, error)
	CheckLogin(username string) (bool, error)
	AuthorizeByToken(token string, login string) (types.ServiceAuthVerdict, error)
	UpdateToken(login string, password string, token string) (string, error)
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
	authVerdict, err := s.AuthRepository.Authorize(*types.MapperAuthDataServiceToDB(&auth_data, auth_data.Token))
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

func (s *AuthService) AuthorizeByToken(token string, login string) (types.ServiceAuthVerdict, error) {
	authVerdict, err := s.AuthRepository.AuthorizeByToken(token, login)
	if err != nil {
		return types.ServiceAuthVerdict{}, err
	}
	return *types.MapperAuthVerdictDBToService(&authVerdict), nil
}

func (s *AuthService) UpdateToken(login string, password string, token string) (string, error) {
	return s.AuthRepository.UpdateToken(login, password, token)
}
