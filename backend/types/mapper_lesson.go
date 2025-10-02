package types

func MapperLessonDBToService(lesson *DBLesson) *ServiceLesson {
	if lesson == nil {
		return nil
	}
	return &ServiceLesson{
		ContractID: lesson.ContractID,
		Duration:   lesson.Duration,
		CreatedAt:  lesson.CreatedAt,
	}
}

func MapperLessonServiceToDB(lesson *ServiceLesson) *DBLesson {
	if lesson == nil {
		return nil
	}
	return &DBLesson{
		ContractID: lesson.ContractID,
		Duration:   lesson.Duration,
		CreatedAt:  lesson.CreatedAt,
	}
}

func MapperLessonServiceToServer(lesson *ServiceLesson) *ServerLesson {
	if lesson == nil {
		return nil
	}
	return &ServerLesson{
		ContractID: lesson.ContractID,
		Duration:   lesson.Duration,
		CreatedAt:  lesson.CreatedAt,
	}
}

func MapperLessonServerToService(lesson *ServerLesson) *ServiceLesson {
	if lesson == nil {
		return nil
	}
	return &ServiceLesson{
		ContractID: lesson.ContractID,
		Duration:   lesson.Duration,
		CreatedAt:  lesson.CreatedAt,
	}
}
