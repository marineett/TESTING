package types

func MapperAdminDBToService(admin *DBAdminData) *ServiceAdminData {
	if admin == nil {
		return nil
	}
	return &ServiceAdminData{
		ID:     admin.ID,
		Salary: admin.Salary,
	}
}

func MapperAdminServiceToDB(admin *ServiceAdminData) *DBAdminData {
	if admin == nil {
		return nil
	}
	return &DBAdminData{
		ID:     admin.ID,
		Salary: admin.Salary,
	}
}

func MapperAdminServiceToServerInit(admin *ServiceInitAdminData) *ServerInitAdminData {
	if admin == nil {
		return nil
	}
	serverInitUserData := MapperUserServiceToServerInit(&admin.ServiceInitUserData)
	return &ServerInitAdminData{
		ServerInitUserData: *serverInitUserData,
		Salary:             admin.Salary,
	}
}

func MapperAdminServerInitToService(admin *ServerInitAdminData) *ServiceInitAdminData {
	if admin == nil {
		return nil
	}
	serviceInitUserData := MapperUserServerInitToService(&admin.ServerInitUserData)
	return &ServiceInitAdminData{
		ServiceInitUserData: *serviceInitUserData,
		Salary:              admin.Salary,
	}
}

func MapperAdminProfileServiceToServer(profile *ServiceAdminProfile) *ServerAdminProfile {
	if profile == nil {
		return nil
	}
	return &ServerAdminProfile{
		FirstName:       profile.FirstName,
		LastName:        profile.LastName,
		MiddleName:      profile.MiddleName,
		TelephoneNumber: profile.TelephoneNumber,
		Email:           profile.Email,
		Salary:          profile.Salary,
	}
}

func MapperAdminProfileServerToService(profile *ServerAdminProfile) *ServiceAdminProfile {
	if profile == nil {
		return nil
	}
	return &ServiceAdminProfile{
		FirstName:       profile.FirstName,
		LastName:        profile.LastName,
		MiddleName:      profile.MiddleName,
		TelephoneNumber: profile.TelephoneNumber,
		Email:           profile.Email,
		Salary:          profile.Salary,
	}
}

func MapperInitAdminServerToService(data *ServerInitAdminData) *ServiceInitAdminData {
	if data == nil {
		return nil
	}
	return &ServiceInitAdminData{
		ServiceInitUserData: ServiceInitUserData{
			ServicePersonalData: ServicePersonalData{
				TelephoneNumber: data.TelephoneNumber,
				Email:           data.Email,
			},
			ServiceAuthData: ServiceAuthData{
				Login:    data.Login,
				Password: data.Password,
			},
		},
	}
}

func MapperAdminDataServiceToDB(data *ServiceAdminData) *DBAdminData {
	if data == nil {
		return nil
	}
	return &DBAdminData{
		ID:     data.ID,
		Salary: data.Salary,
	}
}

func MapperPersonalDataServiceToDB(data *ServicePersonalData) *DBPersonalData {
	if data == nil {
		return nil
	}
	return &DBPersonalData{
		FirstName:       data.FirstName,
		LastName:        data.LastName,
		MiddleName:      data.MiddleName,
		TelephoneNumber: data.TelephoneNumber,
		Email:           data.Email,
		DBPassportData: DBPassportData{
			PassportNumber:   data.PassportNumber,
			PassportSeries:   data.PassportSeries,
			PassportDate:     data.PassportDate,
			PassportIssuedBy: data.PassportIssuedBy,
		},
	}
}

func MapperAuthDataServiceToDB(data *ServiceAuthData) *DBAuthData {
	if data == nil {
		return nil
	}
	return &DBAuthData{
		Login:    data.Login,
		Password: data.Password,
	}
}

func MapperAdminDataDBToService(data *DBAdminData) *ServiceAdminData {
	if data == nil {
		return nil
	}
	return &ServiceAdminData{
		ID:           data.ID,
		Salary:       data.Salary,
		DepartmentID: data.DepartmentID,
	}
}

func MapperAdminProfileDBToService(profile *DBAdminData, personalData *DBPersonalData) *ServiceAdminProfile {
	if profile == nil {
		return nil
	}
	return &ServiceAdminProfile{
		Salary: profile.Salary,
	}
}
