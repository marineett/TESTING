package types

func MapperRepetitorDBToService(repetitor *DBRepetitorData) *ServiceRepetitorData {
	if repetitor == nil {
		return nil
	}
	return &ServiceRepetitorData{
		ID:         repetitor.ID,
		MeanRating: repetitor.SummaryRating / float64(repetitor.ReviewsCount),
		ResumeID:   repetitor.ResumeID,
	}
}

func MapperRepetitorServiceToDB(repetitor *ServiceRepetitorData) *DBRepetitorData {
	if repetitor == nil {
		return nil
	}
	return &DBRepetitorData{
		ID:            repetitor.ID,
		SummaryRating: repetitor.MeanRating,
		ReviewsCount:  1,
		ResumeID:      repetitor.ResumeID,
	}
}

func MapperResumeDBToService(resume *DBResume) *ServiceResume {
	if resume == nil {
		return nil
	}
	return &ServiceResume{
		RepetitorID: resume.RepetitorID,
		Title:       resume.Title,
		Description: resume.Description,
		Prices:      resume.Prices,
		CreatedAt:   resume.CreatedAt,
		UpdatedAt:   resume.UpdatedAt,
	}
}

func MapperResumeServiceToDB(resume *ServiceResume) *DBResume {
	if resume == nil {
		return nil
	}
	return &DBResume{
		RepetitorID: resume.RepetitorID,
		Title:       resume.Title,
		Description: resume.Description,
		Prices:      resume.Prices,
		CreatedAt:   resume.CreatedAt,
		UpdatedAt:   resume.UpdatedAt,
	}
}

func MapperRepetitorServiceToServerInit(repetitor *ServiceInitRepetitorData) *ServerInitRepetitorData {
	if repetitor == nil {
		return nil
	}
	serverInitUserData := MapperUserServiceToServerInit(&repetitor.ServiceInitUserData)
	return &ServerInitRepetitorData{
		ServerInitUserData: *serverInitUserData,
	}
}

func MapperRepetitorServerInitToService(repetitor *ServerInitRepetitorData) *ServiceInitRepetitorData {
	if repetitor == nil {
		return nil
	}
	serviceInitUserData := MapperUserServerInitToService(&repetitor.ServerInitUserData)
	return &ServiceInitRepetitorData{
		ServiceInitUserData: *serviceInitUserData,
	}
}

func MapperResumeServiceToServer(resume *ServiceResume) *ServerResume {
	if resume == nil {
		return nil
	}
	return &ServerResume{
		RepetitorID: resume.RepetitorID,
		Title:       resume.Title,
		Description: resume.Description,
		Prices:      resume.Prices,
		CreatedAt:   resume.CreatedAt,
		UpdatedAt:   resume.UpdatedAt,
	}
}

func MapperResumeServerToService(resume *ServerResume) *ServiceResume {
	if resume == nil {
		return nil
	}
	return &ServiceResume{
		RepetitorID: resume.RepetitorID,
		Title:       resume.Title,
		Description: resume.Description,
		Prices:      resume.Prices,
		CreatedAt:   resume.CreatedAt,
		UpdatedAt:   resume.UpdatedAt,
	}
}

func MapperRepetitorProfileServiceToServer(profile *ServiceRepetitorProfile) *ServerRepetitorProfile {
	if profile == nil {
		return nil
	}
	return &ServerRepetitorProfile{
		FirstName:         profile.FirstName,
		LastName:          profile.LastName,
		MiddleName:        profile.MiddleName,
		TelephoneNumber:   profile.TelephoneNumber,
		Email:             profile.Email,
		MeanRating:        profile.MeanRating,
		ResumeTitle:       profile.ResumeTitle,
		ResumeDescription: profile.ResumeDescription,
		ResumePrices:      profile.ResumePrices,
	}
}

func MapperRepetitorProfileServerToService(profile *ServerRepetitorProfile) *ServiceRepetitorProfile {
	if profile == nil {
		return nil
	}
	return &ServiceRepetitorProfile{
		FirstName:         profile.FirstName,
		LastName:          profile.LastName,
		MiddleName:        profile.MiddleName,
		TelephoneNumber:   profile.TelephoneNumber,
		Email:             profile.Email,
		MeanRating:        profile.MeanRating,
		ResumeTitle:       profile.ResumeTitle,
		ResumeDescription: profile.ResumeDescription,
		ResumePrices:      profile.ResumePrices,
	}
}

func MapperRepetitorViewServiceToServer(view *ServiceRepetitorView) *ServerRepetitorView {
	if view == nil {
		return nil
	}
	return &ServerRepetitorView{
		FirstName:  view.FirstName,
		MeanRating: view.MeanRating,
	}
}

func MapperRepetitorViewServerToService(view *ServerRepetitorView) *ServiceRepetitorView {
	if view == nil {
		return nil
	}
	return &ServiceRepetitorView{
		FirstName:  view.FirstName,
		MeanRating: view.MeanRating,
	}
}

func MapperInitRepetitorServerToService(data *ServerInitRepetitorData) *ServiceInitRepetitorData {
	if data == nil {
		return nil
	}
	return &ServiceInitRepetitorData{
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

// --- V2 Repetitor mappers ---
func MapperRepetitorProfileServiceToServerV2(profile *ServiceRepetitorProfile) *ServerRepetitorProfileV2 {
	if profile == nil {
		return nil
	}
	return &ServerRepetitorProfileV2{
		FirstName:       profile.FirstName,
		LastName:        profile.LastName,
		MiddleName:      profile.MiddleName,
		Email:           profile.Email,
		TelephoneNumber: profile.TelephoneNumber,
		Raiting:         profile.MeanRating,
		Resume:          profile.ResumeDescription,
	}
}

func MapperRepetitorProfileV2ServerToService(profile *ServerRepetitorProfileV2) *ServiceRepetitorProfile {
	if profile == nil {
		return nil
	}
	return &ServiceRepetitorProfile{
		FirstName:         profile.FirstName,
		LastName:          profile.LastName,
		MiddleName:        profile.MiddleName,
		TelephoneNumber:   profile.TelephoneNumber,
		Email:             profile.Email,
		MeanRating:        profile.Raiting,
		ResumeTitle:       "",
		ResumeDescription: profile.Resume,
		ResumePrices:      map[string]int{},
	}
}

func MapperRepetitorProfileWithIDServiceToServerV2(id int64, profile *ServiceRepetitorProfile) *ServerRepetitorProfileWithIDV2 {
	if profile == nil {
		return &ServerRepetitorProfileWithIDV2{ID: id}
	}
	return &ServerRepetitorProfileWithIDV2{
		ID:        id,
		Repetitor: *MapperRepetitorProfileServiceToServerV2(profile),
	}
}

func MapperRepetitorProfileWithIDServerV2ToService(wrapper *ServerRepetitorProfileWithIDV2) (int64, *ServiceRepetitorProfile) {
	if wrapper == nil {
		return 0, nil
	}
	return wrapper.ID, MapperRepetitorProfileV2ServerToService(&wrapper.Repetitor)
}
