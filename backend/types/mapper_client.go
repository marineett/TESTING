package types

func MapperClientDBToService(client *DBClientData) *ServiceClientProfile {
	if client == nil {
		return nil
	}
	return &ServiceClientProfile{
		MeanRating: float64(client.SummaryRating) / float64(client.ReviewsCount),
	}
}

func MapperClientServiceToDB(client *ServiceClientProfile) *DBClientData {
	if client == nil {
		return nil
	}
	return &DBClientData{
		SummaryRating: int64(client.MeanRating * float64(len(client.Reviews))),
		ReviewsCount:  int64(len(client.Reviews)),
	}
}

func MapperClientServiceToServerInit(client *ServiceInitClientData) *ServerInitClientData {
	if client == nil {
		return nil
	}
	serverInitUserData := MapperUserServiceToServerInit(&client.ServiceInitUserData)
	return &ServerInitClientData{
		ServerInitUserData: *serverInitUserData,
	}
}

func MapperClientServerInitToService(client *ServerInitClientData) *ServiceInitClientData {
	if client == nil {
		return nil
	}
	serviceInitUserData := MapperUserServerInitToService(&client.ServerInitUserData)
	return &ServiceInitClientData{
		ServiceInitUserData: *serviceInitUserData,
	}
}

func MapperClientProfileServiceToServer(profile *ServiceClientProfile) *ServerClientProfile {
	if profile == nil {
		return nil
	}
	return &ServerClientProfile{
		FirstName:       profile.FirstName,
		LastName:        profile.LastName,
		MiddleName:      profile.MiddleName,
		TelephoneNumber: profile.TelephoneNumber,
		Email:           profile.Email,
		MeanRating:      profile.MeanRating,
	}
}

func MapperClientProfileServerToService(profile *ServerClientProfile) *ServiceClientProfile {
	if profile == nil {
		return nil
	}
	return &ServiceClientProfile{
		FirstName:       profile.FirstName,
		LastName:        profile.LastName,
		MiddleName:      profile.MiddleName,
		TelephoneNumber: profile.TelephoneNumber,
		Email:           profile.Email,
		MeanRating:      profile.MeanRating,
	}
}

func MapperInitClientServerToService(data *ServerInitClientData) *ServiceInitClientData {
	if data == nil {
		return nil
	}
	return &ServiceInitClientData{
		ServiceInitUserData: ServiceInitUserData{
			ServicePersonalData: ServicePersonalData{
				TelephoneNumber: data.TelephoneNumber,
				Email:           data.Email,
				ServicePassportData: ServicePassportData{
					PassportNumber:   data.PassportNumber,
					PassportSeries:   data.PassportSeries,
					PassportDate:     data.PassportDate,
					PassportIssuedBy: data.PassportIssuedBy,
				},
				FirstName:  data.FirstName,
				LastName:   data.LastName,
				MiddleName: data.MiddleName,
			},
			ServiceAuthData: ServiceAuthData{
				Login:    data.Login,
				Password: data.Password,
			},
		},
	}
}
