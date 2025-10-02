package types

func MapperDepartmentDBToService(department *DBDepartment) *ServiceDepartment {
	if department == nil {
		return nil
	}
	return &ServiceDepartment{
		ID:     department.ID,
		Name:   department.Name,
		HeadID: department.HeadID,
	}
}

func MapperDepartmentServiceToDB(department *ServiceDepartment) *DBDepartment {
	if department == nil {
		return nil
	}
	return &DBDepartment{
		ID:     department.ID,
		Name:   department.Name,
		HeadID: department.HeadID,
	}
}

func MapperDepartmentInitDataServiceToServer(department *ServiceDepartmentInitData) *ServerDepartmentInitData {
	if department == nil {
		return nil
	}
	return &ServerDepartmentInitData{
		Name:   department.Name,
		HeadID: department.HeadID,
	}
}

func MapperDepartmentInitDataServerToService(department *ServerDepartmentInitData) *ServiceDepartmentInitData {
	if department == nil {
		return nil
	}
	return &ServiceDepartmentInitData{
		Name:   department.Name,
		HeadID: department.HeadID,
	}
}

func MapperDepartmentServiceToServer(department *ServiceDepartment) *ServerDepartment {
	if department == nil {
		return nil
	}
	return &ServerDepartment{
		Name:   department.Name,
		HeadID: department.HeadID,
	}
}

func MapperDepartmentServerToService(department *ServerDepartment) *ServiceDepartment {
	if department == nil {
		return nil
	}
	return &ServiceDepartment{
		Name:   department.Name,
		HeadID: department.HeadID,
	}
}
