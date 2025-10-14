package types

func MapperModeratorDBToService(moderator *DBModeratorData) *ServiceModeratorData {
	if moderator == nil {
		return nil
	}
	return &ServiceModeratorData{
		ID:     moderator.ID,
		Salary: moderator.Salary,
	}
}

func MapperModeratorServiceToDB(moderator *ServiceModeratorData) *DBModeratorData {
	if moderator == nil {
		return nil
	}
	return &DBModeratorData{
		ID:     moderator.ID,
		Salary: moderator.Salary,
	}
}

func MapperModeratorProfileServiceToServer(profile *ServiceModeratorProfile) *ServerModeratorProfile {
	if profile == nil {
		return nil
	}
	return &ServerModeratorProfile{
		FirstName:       profile.FirstName,
		LastName:        profile.LastName,
		MiddleName:      profile.MiddleName,
		TelephoneNumber: profile.TelephoneNumber,
		Email:           profile.Email,
	}
}

func MapperModeratorProfileServerToService(profile *ServerModeratorProfile) *ServiceModeratorProfile {
	if profile == nil {
		return nil
	}
	return &ServiceModeratorProfile{
		FirstName:       profile.FirstName,
		LastName:        profile.LastName,
		MiddleName:      profile.MiddleName,
		TelephoneNumber: profile.TelephoneNumber,
		Email:           profile.Email,
	}
}

func MapperModeratorProfileWithIDServiceToServer(profile *ServiceModeratorProfileWithID) *ServerModeratorProfileWithID {
	if profile == nil {
		return nil
	}
	return &ServerModeratorProfileWithID{
		ID: profile.ID,
		Moderator: ServerModeratorProfile{
			FirstName:       profile.FirstName,
			LastName:        profile.LastName,
			MiddleName:      profile.MiddleName,
			TelephoneNumber: profile.TelephoneNumber,
			Email:           profile.Email,
		},
	}
}

func MapperModeratorProfileWithIDServerToService(profile *ServerModeratorProfileWithID) *ServiceModeratorProfileWithID {
	if profile == nil {
		return nil
	}
	return &ServiceModeratorProfileWithID{
		ID: profile.ID,
		ServiceModeratorProfile: ServiceModeratorProfile{
			FirstName:       profile.Moderator.FirstName,
			LastName:        profile.Moderator.LastName,
			MiddleName:      profile.Moderator.MiddleName,
			TelephoneNumber: profile.Moderator.TelephoneNumber,
			Email:           profile.Moderator.Email,
		},
	}
}

func MapperInitModeratorServerToService(data *ServerInitModeratorData) *ServiceInitModeratorData {
	if data == nil {
		return nil
	}
	return &ServiceInitModeratorData{
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

func MapperModeratorDataDBToService(data *DBModeratorData) *ServiceModeratorData {
	if data == nil {
		return nil
	}
	return &ServiceModeratorData{
		ID:     data.ID,
		Salary: data.Salary,
	}
}
