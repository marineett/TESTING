package service_logic

import (
	"data_base_project/data_base"
	"data_base_project/types"
	"errors"
)

type IModeratorService interface {
	CreateModerator(init_data types.ServiceInitModeratorData) error
	GetModeratorData(user_id int64) (*types.ServiceModeratorData, error)
	GetModeratorProfile(user_id int64) (*types.ServiceModeratorProfile, error)
	UpdateModeratorPersonalData(user_id int64, personal_data types.ServicePersonalData) error
	UpdateModeratorPassword(user_id int64, authData types.ServiceAuthData, newPassword string) error
	UpdateModeratorSalary(user_id int64, salary int64) error
	GetModerators() ([]*types.ServiceModeratorProfileWithID, error)
	GetModeratorProfileWithId(moderator_id int64) (*types.ServiceModeratorProfileWithID, error)
}

type ModeratorService struct {
	moderatorRepository    data_base.IModeratorRepository
	personalDataRepository data_base.IPersonalDataRepository
	userRepository         data_base.IUserRepository
	departmentRepository   data_base.IDepartmentRepository
}

func CreateModeratorService(
	moderatorRepository data_base.IModeratorRepository,
	personalDataRepository data_base.IPersonalDataRepository,
	userRepository data_base.IUserRepository,
	departmentRepository data_base.IDepartmentRepository,
) IModeratorService {
	return &ModeratorService{
		moderatorRepository:    moderatorRepository,
		personalDataRepository: personalDataRepository,
		userRepository:         userRepository,
		departmentRepository:   departmentRepository,
	}
}

func (s *ModeratorService) transormInitModeratorData(init_data types.ServiceInitModeratorData) (*types.ServiceModeratorData, *types.ServicePersonalData, *types.ServiceAuthData) {
	moderator := types.ServiceModeratorData{
		Salary: int64(init_data.Salary),
	}
	return &moderator, &init_data.ServicePersonalData, &types.ServiceAuthData{
		Login:    init_data.Login,
		Password: init_data.Password,
	}
}

func (s *ModeratorService) CreateModerator(init_data types.ServiceInitModeratorData) error {
	moderator_data, personal_data, auth_info := s.transormInitModeratorData(init_data)
	_, err := s.moderatorRepository.InsertModerator(types.DBModeratorData{
		Salary: moderator_data.Salary,
	}, types.DBPersonalData{
		FirstName:       personal_data.FirstName,
		LastName:        personal_data.LastName,
		MiddleName:      personal_data.MiddleName,
		TelephoneNumber: personal_data.TelephoneNumber,
		Email:           personal_data.Email,
	}, types.DBAuthData{
		Login:    auth_info.Login,
		Password: auth_info.Password,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *ModeratorService) GetModeratorData(user_id int64) (*types.ServiceModeratorData, error) {
	moderator, err := s.moderatorRepository.GetModerator(user_id)
	if err != nil {
		return nil, err
	}
	return types.MapperModeratorDataDBToService(moderator), nil
}

func (s *ModeratorService) UpdateModeratorPersonalData(user_id int64, personal_data types.ServicePersonalData) error {
	return s.moderatorRepository.UpdateModeratorPersonalData(user_id, *types.MapperPersonalDataServiceToDB(&personal_data))
}

func (s *ModeratorService) UpdateModeratorPassword(user_id int64, authData types.ServiceAuthData, newPassword string) error {
	return s.moderatorRepository.UpdateModeratorPassword(user_id, types.DBAuthData{
		Login:    authData.Login,
		Password: authData.Password,
	}, newPassword)
}

func (s *ModeratorService) GetModeratorProfile(user_id int64) (*types.ServiceModeratorProfile, error) {
	user_data, err := s.userRepository.GetUser(user_id)
	if err != nil {
		return nil, err
	}
	moderator_data, err := s.GetModeratorData(user_id)
	if err != nil {
		return nil, err
	}
	personal_data, err := s.personalDataRepository.GetPersonalData(user_data.PersonalDataID)
	if err != nil {
		return nil, err
	}
	departmentsIds, err := s.departmentRepository.GetUserDepartmentsIDs(user_id)
	if err != nil {
		return nil, err
	}
	departments := make([]string, len(departmentsIds))
	for _, departmentId := range departmentsIds {
		department, err := s.departmentRepository.GetDepartment(departmentId)
		departments = append(departments, department.Name)
		if err != nil {
			return nil, err
		}
	}
	return &types.ServiceModeratorProfile{
		FirstName:       personal_data.FirstName,
		LastName:        personal_data.LastName,
		MiddleName:      personal_data.MiddleName,
		TelephoneNumber: personal_data.TelephoneNumber,
		Email:           personal_data.Email,
		Salary:          moderator_data.Salary,
		Departments:     departments,
	}, nil
}

func (s *ModeratorService) GetModerators() ([]*types.ServiceModeratorProfileWithID, error) {
	moderatorsIDs, err := s.moderatorRepository.GetModerators()
	if err != nil {
		return nil, err
	}
	moderatorsWithID := make([]*types.ServiceModeratorProfileWithID, len(moderatorsIDs))
	for i, moderatorID := range moderatorsIDs {
		moderatorProfile, err := s.GetModeratorProfile(moderatorID)
		if err != nil {
			return nil, err
		}
		moderatorsWithID[i] = &types.ServiceModeratorProfileWithID{
			ID:                      moderatorID,
			ServiceModeratorProfile: *moderatorProfile,
		}
	}
	return moderatorsWithID, nil
}

func (s *ModeratorService) GetModeratorProfileWithId(moderator_id int64) (*types.ServiceModeratorProfileWithID, error) {
	moderatorProfile, err := s.GetModeratorProfile(moderator_id)
	if err != nil {
		return nil, err
	}
	moderatorProfile.Departments = nil
	return &types.ServiceModeratorProfileWithID{
		ID:                      moderator_id,
		ServiceModeratorProfile: *moderatorProfile,
	}, nil
}

func (s *ModeratorService) UpdateModeratorSalary(user_id int64, salary int64) error {
	if salary <= 0 {
		return errors.New("salary must be greater than 0")
	}
	return s.moderatorRepository.UpdateModeratorSalary(user_id, salary)
}
