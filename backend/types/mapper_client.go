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

// --- V2 Client mappers ---
func MapperClientProfileServiceToServerV2(profile *ServiceClientProfile) *ServerClientProfileV2 {
	if profile == nil {
		return nil
	}
	return &ServerClientProfileV2{
		FirstName:       profile.FirstName,
		LastName:        profile.LastName,
		MiddleName:      profile.MiddleName,
		TelephoneNumber: profile.TelephoneNumber,
		Email:           profile.Email,
		Raiting:         profile.MeanRating,
	}
}

func MapperClientProfileV2ServerToService(profile *ServerClientProfileV2) *ServiceClientProfile {
	if profile == nil {
		return nil
	}
	return &ServiceClientProfile{
		FirstName:       profile.FirstName,
		LastName:        profile.LastName,
		MiddleName:      profile.MiddleName,
		TelephoneNumber: profile.TelephoneNumber,
		Email:           profile.Email,
		MeanRating:      profile.Raiting,
	}
}

func MapperClientProfileWithIDServiceToServerV2(id int64, profile *ServiceClientProfile) *ServerClientProfileWithIDV2 {
	if profile == nil {
		return &ServerClientProfileWithIDV2{ID: id}
	}
	return &ServerClientProfileWithIDV2{
		ID:     id,
		Client: *MapperClientProfileServiceToServerV2(profile),
	}
}

func MapperClientProfileWithIDServerV2ToService(wrapper *ServerClientProfileWithIDV2) (int64, *ServiceClientProfile) {
	if wrapper == nil {
		return 0, nil
	}
	return wrapper.ID, MapperClientProfileV2ServerToService(&wrapper.Client)
}
