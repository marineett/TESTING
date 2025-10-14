package service_logic

import (
	"data_base_project/data_base"
	"data_base_project/types"
	"errors"
)

type IAdminService interface {
	CreateAdmin(initData types.ServiceInitAdminData) error
	GetAdminData(userID int64) (*types.ServiceAdminData, error)
	GetAdminProfile(userID int64) (*types.ServiceAdminProfile, error)
	UpdateAdminPersonalData(userID int64, personalData types.ServicePersonalData) error
	UpdateAdminPassword(userID int64, authData types.ServiceAuthData, newPassword string) error
	UpdateAdminDepartment(userID int64, departmentID int64) error
	UpdateAdminSalary(userID int64, salary int64) error
}

type AdminService struct {
	adminRepository        data_base.IAdminRepository
	userRepository         data_base.IUserRepository
	personalDataRepository data_base.IPersonalDataRepository
}

func CreateAdminService(
	adminRepository data_base.IAdminRepository,
	userRepository data_base.IUserRepository,
	personalDataRepository data_base.IPersonalDataRepository,
) IAdminService {
	return &AdminService{
		adminRepository:        adminRepository,
		userRepository:         userRepository,
		personalDataRepository: personalDataRepository,
	}
}

func (s *AdminService) transormInitAdminData(initData types.ServiceInitAdminData) (*types.ServiceAdminData, *types.ServicePersonalData, *types.ServiceAuthData) {
	adminData := types.ServiceAdminData{
		Salary:       initData.Salary,
		DepartmentID: 0,
	}
	personalData := types.ServicePersonalData{
		FirstName:       initData.FirstName,
		LastName:        initData.LastName,
		MiddleName:      initData.MiddleName,
		TelephoneNumber: initData.TelephoneNumber,
		Email:           initData.Email,
	}
	authData := types.ServiceAuthData{
		Login:    initData.Login,
		Password: initData.Password,
	}
	return &adminData, &personalData, &authData
}
func (s *AdminService) CreateAdmin(initData types.ServiceInitAdminData) error {
	adminData, personalData, authData := s.transormInitAdminData(initData)
	if adminData.Salary <= 0 {
		return errors.New("salary must be greater than 0")
	}
	_, err := s.adminRepository.InsertAdmin(
		*types.MapperAdminDataServiceToDB(adminData),
		*types.MapperPersonalDataServiceToDB(personalData),
		*types.MapperAuthDataServiceToDB(authData),
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *AdminService) GetAdminData(userID int64) (*types.ServiceAdminData, error) {
	adminData, err := s.adminRepository.GetAdmin(userID)
	if err != nil {
		return nil, err
	}
	return types.MapperAdminDataDBToService(adminData), nil
}

func (s *AdminService) UpdateAdminPersonalData(userID int64, personalData types.ServicePersonalData) error {
	return s.adminRepository.UpdateAdminPersonalData(userID, types.DBPersonalData{
		FirstName:       personalData.FirstName,
		LastName:        personalData.LastName,
		MiddleName:      personalData.MiddleName,
		TelephoneNumber: personalData.TelephoneNumber,
		Email:           personalData.Email,
	})
}

func (s *AdminService) UpdateAdminPassword(userID int64, authData types.ServiceAuthData, newPassword string) error {
	return s.adminRepository.UpdateAdminPassword(userID, *types.MapperAuthDataServiceToDB(&authData), newPassword)
}

func (s *AdminService) UpdateAdminDepartment(userID int64, departmentID int64) error {
	return s.adminRepository.UpdateAdminDepartment(userID, departmentID)
}

func (s *AdminService) UpdateAdminSalary(userID int64, salary int64) error {
	if salary <= 0 {
		return errors.New("salary must be greater than 0")
	}
	return s.adminRepository.UpdateAdminSalary(userID, salary)
}

func (s *AdminService) GetAdminProfile(userID int64) (*types.ServiceAdminProfile, error) {
	userData, err := s.userRepository.GetUser(userID)
	if err != nil {
		return nil, err
	}
	personalData, err := s.personalDataRepository.GetPersonalData(userData.PersonalDataID)
	if err != nil {
		return nil, err
	}
	adminData, err := s.adminRepository.GetAdmin(userID)
	if err != nil {
		return nil, err
	}
	return &types.ServiceAdminProfile{
		FirstName:       personalData.FirstName,
		LastName:        personalData.LastName,
		MiddleName:      personalData.MiddleName,
		TelephoneNumber: personalData.TelephoneNumber,
		Email:           personalData.Email,
		Salary:          adminData.Salary,
	}, nil
}
