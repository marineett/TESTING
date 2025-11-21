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
		Salary:          profile.Salary,
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
			Salary:          profile.Salary,
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
			Salary:          profile.Moderator.Salary,
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

// Fills profile fields that are available from ServiceModeratorData
func MapperModeratorDataServiceToProfileService(data *ServiceModeratorData) *ServiceModeratorProfile {
	if data == nil {
		return nil
	}
	return &ServiceModeratorProfile{
		Salary:      data.Salary,
		Departments: data.Departments,
	}
}

func MapperDepartmentServiceToServerV2(dep *ServiceDepartment, moderators []ServiceModeratorProfileWithID) *ServerDepartmentV2 {
	if dep == nil {
		return nil
	}
	result := ServerDepartmentV2{ID: dep.ID, Name: dep.Name, HeadID: dep.HeadID}
	result.Moderators = make([]ServerModeratorProfileWithIDV2, 0, len(moderators))
	for _, m := range moderators {
		result.Moderators = append(result.Moderators, ServerModeratorProfileWithIDV2{
			ID: m.ID,
			Moderator: ServerModeratorProfile{
				FirstName:       m.FirstName,
				LastName:        m.LastName,
				MiddleName:      m.MiddleName,
				TelephoneNumber: m.TelephoneNumber,
				Email:           m.Email,
				Salary:          m.Salary,
			},
		})
	}
	return &result
}

func MapperDepartmentCreateV2ServerToService(req *ServerDepartmentCreateV2) *ServiceDepartmentInitData {
	if req == nil {
		return nil
	}
	return &ServiceDepartmentInitData{Name: req.Name, HeadID: req.HeadID}
}

func MapperDepartmentNameUpdateV2ServerToService(req *ServerDepartmentNameUpdateV2, departmentID int64, headID int64) *ServiceDepartmentInitData {
	if req == nil {
		return nil
	}
	return &ServiceDepartmentInitData{ID: departmentID, Name: req.Name, HeadID: headID}
}
