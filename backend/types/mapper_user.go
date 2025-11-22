package types

func MapperUserDBToService(user *DBUserData) *ServiceUserData {
	if user == nil {
		return nil
	}
	return &ServiceUserData{
		ID:               user.ID,
		RegistrationDate: user.RegistrationDate,
		LastLoginDate:    user.LastLoginDate,
		PersonalDataID:   user.PersonalDataID,
	}
}

func MapperUserServiceToDB(user *ServiceUserData) *DBUserData {
	if user == nil {
		return nil
	}
	return &DBUserData{
		ID:               user.ID,
		RegistrationDate: user.RegistrationDate,
		LastLoginDate:    user.LastLoginDate,
		PersonalDataID:   user.PersonalDataID,
	}
}

func MapperAuthDBToService(auth *DBAuthData) *ServiceAuthData {
	if auth == nil {
		return nil
	}
	return &ServiceAuthData{
		Login:    auth.Login,
		Password: auth.Password,
	}
}

func MapperAuthServiceToDB(auth *ServiceAuthData) *DBAuthData {
	if auth == nil {
		return nil
	}
	return &DBAuthData{
		Login:    auth.Login,
		Password: auth.Password,
	}
}

func MapperAuthVerdictDBToService(verdict *DBAuthVerdict) *ServiceAuthVerdict {
	if verdict == nil {
		return nil
	}
	return &ServiceAuthVerdict{
		UserID:   verdict.UserID,
		UserType: verdict.UserType,
	}
}

func MapperAuthVerdictServiceToDB(verdict *ServiceAuthVerdict) *DBAuthVerdict {
	if verdict == nil {
		return nil
	}
	return &DBAuthVerdict{
		UserID:   verdict.UserID,
		UserType: verdict.UserType,
	}
}

func MapperPassportDBToService(passport *DBPassportData) *ServicePassportData {
	if passport == nil {
		return nil
	}
	return &ServicePassportData{
		PassportNumber:   passport.PassportNumber,
		PassportSeries:   passport.PassportSeries,
		PassportDate:     passport.PassportDate,
		PassportIssuedBy: passport.PassportIssuedBy,
	}
}

func MapperPassportServiceToDB(passport *ServicePassportData) *DBPassportData {
	if passport == nil {
		return nil
	}
	return &DBPassportData{
		PassportNumber:   passport.PassportNumber,
		PassportSeries:   passport.PassportSeries,
		PassportDate:     passport.PassportDate,
		PassportIssuedBy: passport.PassportIssuedBy,
	}
}

func MapperPersonalDBToService(personal *DBPersonalData) *ServicePersonalData {
	if personal == nil {
		return nil
	}
	return &ServicePersonalData{
		TelephoneNumber: personal.TelephoneNumber,
		Email:           personal.Email,
		ServicePassportData: ServicePassportData{
			PassportNumber:   personal.PassportNumber,
			PassportSeries:   personal.PassportSeries,
			PassportDate:     personal.PassportDate,
			PassportIssuedBy: personal.PassportIssuedBy,
		},
		FirstName:  personal.FirstName,
		LastName:   personal.LastName,
		MiddleName: personal.MiddleName,
	}
}

func MapperPersonalServiceToDB(personal *ServicePersonalData) *DBPersonalData {
	if personal == nil {
		return nil
	}
	return &DBPersonalData{
		TelephoneNumber: personal.TelephoneNumber,
		Email:           personal.Email,
		DBPassportData: DBPassportData{
			PassportNumber:   personal.PassportNumber,
			PassportSeries:   personal.PassportSeries,
			PassportDate:     personal.PassportDate,
			PassportIssuedBy: personal.PassportIssuedBy,
		},
		FirstName:  personal.FirstName,
		LastName:   personal.LastName,
		MiddleName: personal.MiddleName,
	}
}

func MapperUserServiceToServerInit(user *ServiceInitUserData) *ServerInitUserData {
	if user == nil {
		return nil
	}
	return &ServerInitUserData{
		ServerPersonalData: ServerPersonalData{
			TelephoneNumber: user.TelephoneNumber,
			Email:           user.Email,
			ServerPassportData: ServerPassportData{
				PassportNumber:   user.PassportNumber,
				PassportSeries:   user.PassportSeries,
				PassportDate:     user.PassportDate,
				PassportIssuedBy: user.PassportIssuedBy,
			},
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			MiddleName: user.MiddleName,
		},
		ServerAuthData: ServerAuthData{
			Login:    user.Login,
			Password: user.Password,
		},
	}
}

func MapperUserServerInitToService(user *ServerInitUserData) *ServiceInitUserData {
	if user == nil {
		return nil
	}
	return &ServiceInitUserData{
		ServicePersonalData: ServicePersonalData{
			TelephoneNumber: user.TelephoneNumber,
			Email:           user.Email,
			ServicePassportData: ServicePassportData{
				PassportNumber:   user.PassportNumber,
				PassportSeries:   user.PassportSeries,
				PassportDate:     user.PassportDate,
				PassportIssuedBy: user.PassportIssuedBy,
			},
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			MiddleName: user.MiddleName,
		},
		ServiceAuthData: ServiceAuthData{
			Login:    user.Login,
			Password: user.Password,
		},
	}
}

func MapperUserServiceToServer(user *ServiceUserData) *ServerUserData {
	if user == nil {
		return nil
	}
	return &ServerUserData{
		ID:               user.ID,
		RegistrationDate: user.RegistrationDate,
		LastLoginDate:    user.LastLoginDate,
		PersonalDataID:   user.PersonalDataID,
	}
}

func MapperUserServerToService(user *ServerUserData) *ServiceUserData {
	if user == nil {
		return nil
	}
	return &ServiceUserData{
		ID:               user.ID,
		RegistrationDate: user.RegistrationDate,
		LastLoginDate:    user.LastLoginDate,
		PersonalDataID:   user.PersonalDataID,
	}
}

func MapperPersonalServiceToServer(personal *ServicePersonalData) *ServerPersonalData {
	if personal == nil {
		return nil
	}
	return &ServerPersonalData{
		TelephoneNumber: personal.TelephoneNumber,
		Email:           personal.Email,
		ServerPassportData: ServerPassportData{
			PassportNumber:   personal.PassportNumber,
			PassportSeries:   personal.PassportSeries,
			PassportDate:     personal.PassportDate,
			PassportIssuedBy: personal.PassportIssuedBy,
		},
		FirstName:  personal.FirstName,
		LastName:   personal.LastName,
		MiddleName: personal.MiddleName,
	}
}

func MapperPersonalServerToService(personal *ServerPersonalData) *ServicePersonalData {
	if personal == nil {
		return nil
	}
	return &ServicePersonalData{
		TelephoneNumber: personal.TelephoneNumber,
		Email:           personal.Email,
		ServicePassportData: ServicePassportData{
			PassportNumber:   personal.PassportNumber,
			PassportSeries:   personal.PassportSeries,
			PassportDate:     personal.PassportDate,
			PassportIssuedBy: personal.PassportIssuedBy,
		},
		FirstName:  personal.FirstName,
		LastName:   personal.LastName,
		MiddleName: personal.MiddleName,
	}
}

func MapperAuthServiceToServer(auth *ServiceAuthData) *ServerAuthData {
	if auth == nil {
		return nil
	}
	return &ServerAuthData{
		Login:    auth.Login,
		Password: auth.Password,
	}
}

func MapperAuthServerToService(auth *ServerAuthData) *ServiceAuthData {
	if auth == nil {
		return nil
	}
	return &ServiceAuthData{
		Login:    auth.Login,
		Password: auth.Password,
	}
}

func MapperPassportServiceToServer(passport *ServicePassportData) *ServerPassportData {
	if passport == nil {
		return nil
	}
	return &ServerPassportData{
		PassportNumber:   passport.PassportNumber,
		PassportSeries:   passport.PassportSeries,
		PassportDate:     passport.PassportDate,
		PassportIssuedBy: passport.PassportIssuedBy,
	}
}

func MapperPassportServerToService(passport *ServerPassportData) *ServicePassportData {
	if passport == nil {
		return nil
	}
	return &ServicePassportData{
		PassportNumber:   passport.PassportNumber,
		PassportSeries:   passport.PassportSeries,
		PassportDate:     passport.PassportDate,
		PassportIssuedBy: passport.PassportIssuedBy,
	}
}

func MapperVerdictServiceToServer(verdict *ServiceAuthVerdict) *ServerVerdict {
	if verdict == nil {
		return nil
	}
	return &ServerVerdict{
		UserID:   verdict.UserID,
		UserType: verdict.UserType,
	}
}

func MapperVerdictServerToService(verdict *ServerVerdict) *ServiceAuthVerdict {
	if verdict == nil {
		return nil
	}
	return &ServiceAuthVerdict{
		UserID:   verdict.UserID,
		UserType: verdict.UserType,
	}
}

func MapperPersonalDataDBToService(data *DBPersonalData) *ServicePersonalData {
	if data == nil {
		return nil
	}
	return &ServicePersonalData{
		FirstName:       data.FirstName,
		LastName:        data.LastName,
		MiddleName:      data.MiddleName,
		TelephoneNumber: data.TelephoneNumber,
		Email:           data.Email,
		ServicePassportData: ServicePassportData{
			PassportNumber:   data.PassportNumber,
			PassportSeries:   data.PassportSeries,
			PassportDate:     data.PassportDate,
			PassportIssuedBy: data.PassportIssuedBy,
		},
	}
}

func MapperRegistrationServerToClientRegistration(registration *ServerRegistrationDataV2) *ServerInitClientData {
	if registration == nil {
		return nil
	}
	return &ServerInitClientData{}
}

func MapperRegistrationV2ToServiceInitClient(data *ServerRegistrationDataV2) *ServiceInitClientData {
	if data == nil {
		return nil
	}
	return &ServiceInitClientData{
		ServiceInitUserData: ServiceInitUserData{
			ServicePersonalData: ServicePersonalData{
				TelephoneNumber:     data.TelephoneNumber,
				Email:               data.Email,
				ServicePassportData: ServicePassportData{},
				FirstName:           data.FirstName,
				LastName:            data.LastName,
				MiddleName:          data.MiddleName,
			},
			ServiceAuthData: ServiceAuthData{
				Login:    data.Login,
				Password: data.Password,
			},
		},
	}
}

func MapperRegistrationV2ToServiceInitRepetitor(data *ServerRegistrationDataV2) *ServiceInitRepetitorData {
	if data == nil {
		return nil
	}
	return &ServiceInitRepetitorData{
		ServiceInitUserData: ServiceInitUserData{
			ServicePersonalData: ServicePersonalData{
				TelephoneNumber:     data.TelephoneNumber,
				Email:               data.Email,
				ServicePassportData: ServicePassportData{},
				FirstName:           data.FirstName,
				LastName:            data.LastName,
				MiddleName:          data.MiddleName,
			},
			ServiceAuthData: ServiceAuthData{
				Login:    data.Login,
				Password: data.Password,
			},
		},
	}
}

func MapperRegistrationV2ToServiceInitModerator(data *ServerRegistrationDataV2) *ServiceInitModeratorData {
	if data == nil {
		return nil
	}
	return &ServiceInitModeratorData{
		ServiceInitUserData: ServiceInitUserData{
			ServicePersonalData: ServicePersonalData{
				TelephoneNumber:     data.TelephoneNumber,
				Email:               data.Email,
				ServicePassportData: ServicePassportData{},
				FirstName:           data.FirstName,
				LastName:            data.LastName,
				MiddleName:          data.MiddleName,
			},
			ServiceAuthData: ServiceAuthData{
				Login:    data.Login,
				Password: data.Password,
			},
		},
		Salary: data.Salary,
	}
}

func MapperRegistrationV2ToServiceInitAdmin(data *ServerRegistrationDataV2) *ServiceInitAdminData {
	if data == nil {
		return nil
	}
	return &ServiceInitAdminData{
		ServiceInitUserData: ServiceInitUserData{
			ServicePersonalData: ServicePersonalData{
				TelephoneNumber:     data.TelephoneNumber,
				Email:               data.Email,
				ServicePassportData: ServicePassportData{},
				FirstName:           data.FirstName,
				LastName:            data.LastName,
				MiddleName:          data.MiddleName,
			},
			ServiceAuthData: ServiceAuthData{
				Login:    data.Login,
				Password: data.Password,
			},
		},
		Salary: int64(data.Salary),
	}
}

// V2 client/repetitor mappers are defined in dedicated files:
// - mapper_client.go
// - mapper_repetitor.go
