package service_logic

import (
	"data_base_project/data_base"
	"data_base_project/types"
	"errors"
)

type IDepartmentService interface {
	CreateDepartment(department types.ServiceDepartmentInitData) error
	GetDepartmentsByHeadID(headID int64) ([]types.ServiceDepartment, error)
	GetDepartmentIdByName(name string) (int64, error)
	GetDepartment(id int64) (types.ServiceDepartment, error)
	AssignAdminToDepartment(adminID int64, departmentID int64) error
	FireAdminFromDepartment(adminID int64, departmentID int64) error
	AssignModeratorToDepartment(moderatorID int64, departmentID int64) error
	FireModeratorFromDepartment(moderatorID int64, departmentID int64) error
	GetDepartmentUsersIDs(departmentID int64) ([]int64, error)
	GetUserDepartmentsIDs(userID int64) ([]int64, error)
}

type DepartmentService struct {
	departmentRepository data_base.IDepartmentRepository
}

func CreateDepartmentService(
	departmentRepository data_base.IDepartmentRepository,
	moderatorRepository data_base.IModeratorRepository,
) IDepartmentService {
	return &DepartmentService{
		departmentRepository: departmentRepository,
	}
}

func (s *DepartmentService) CreateDepartment(department types.ServiceDepartmentInitData) error {
	serviceDepartment := types.ServiceDepartment{
		Name:   department.Name,
		HeadID: department.HeadID,
	}
	_, err := s.departmentRepository.InsertDepartment(*types.MapperDepartmentServiceToDB(&serviceDepartment))
	if err != nil {
		return err
	}
	return nil
}

func (s *DepartmentService) GetDepartment(id int64) (types.ServiceDepartment, error) {
	department, err := s.departmentRepository.GetDepartment(id)
	if err != nil {
		return types.ServiceDepartment{}, err
	}
	return *types.MapperDepartmentDBToService(department), nil
}

func (s *DepartmentService) GetDepartmentsByHeadID(headID int64) ([]types.ServiceDepartment, error) {
	departments, err := s.departmentRepository.GetDepartmentsByHeadID(headID)
	if err != nil {
		return nil, err
	}
	serviceDepartments := make([]types.ServiceDepartment, len(departments))
	for i, department := range departments {
		serviceDepartments[i] = *types.MapperDepartmentDBToService(&department)
	}
	return serviceDepartments, nil
}
func (s *DepartmentService) GetDepartmentIdByName(name string) (int64, error) {
	id, err := s.departmentRepository.GetDepartmentIdByName(name)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *DepartmentService) AssignAdminToDepartment(adminID int64, departmentID int64) error {
	department, err := s.departmentRepository.GetDepartment(departmentID)
	if err != nil {
		return err
	}
	if department.HeadID != 0 {
		return errors.New("department already has a head")
	}
	err = s.departmentRepository.ChangeDepartmentHead(departmentID, adminID)
	if err != nil {
		return err
	}
	return nil
}

func (s *DepartmentService) AssignModeratorToDepartment(moderatorID int64, departmentID int64) error {
	userDepartmentsIDs, err := s.departmentRepository.GetUserDepartmentsIDs(moderatorID)
	if err != nil {
		return err
	}
	for _, currentDepartmentID := range userDepartmentsIDs {
		if currentDepartmentID == departmentID {
			return errors.New("moderator already in this department")
		}
	}
	err = s.departmentRepository.HireInfoInsert(types.DBHireInfo{
		UserID:       moderatorID,
		DepartmentID: departmentID,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *DepartmentService) FireAdminFromDepartment(adminID int64, departmentID int64) error {
	department, err := s.departmentRepository.GetDepartment(departmentID)
	if err != nil {
		return err
	}
	if department.HeadID != adminID {
		return errors.New("admin not in this department")
	}
	err = s.departmentRepository.ChangeDepartmentHead(departmentID, 0)
	if err != nil {
		return err
	}
	return nil
}

func (s *DepartmentService) FireModeratorFromDepartment(moderatorID int64, departmentID int64) error {
	err := s.departmentRepository.HireInfoDelete(moderatorID, departmentID)
	if err != nil {
		return err
	}
	return nil
}

func (s *DepartmentService) GetDepartmentUsersIDs(departmentID int64) ([]int64, error) {
	_, err := s.departmentRepository.GetDepartment(departmentID)
	if err != nil {
		return nil, err
	}
	usersIDs, err := s.departmentRepository.GetDepartmentUsersIDs(departmentID)
	if err != nil {
		return nil, err
	}
	return usersIDs, nil
}

func (s *DepartmentService) GetUserDepartmentsIDs(userID int64) ([]int64, error) {

	departmentsIDs, err := s.departmentRepository.GetUserDepartmentsIDs(userID)
	if err != nil {
		return nil, err
	}
	return departmentsIDs, nil
}
