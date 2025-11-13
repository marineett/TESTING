package types

import "time"

func MapperLessonDBToService(lesson *DBLesson) *ServiceLesson {
	if lesson == nil {
		return nil
	}
	return &ServiceLesson{
		ID:         lesson.ID,
		ContractID: lesson.ContractID,
		Duration:   lesson.Duration,
		Format:     lesson.Format,
		CreatedAt:  lesson.CreatedAt,
	}
}

func MapperLessonServiceToDB(lesson *ServiceLesson) *DBLesson {
	if lesson == nil {
		return nil
	}
	return &DBLesson{
		ID:         lesson.ID,
		ContractID: lesson.ContractID,
		Duration:   lesson.Duration,
		Format:     lesson.Format,
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

// --- V2 Lesson mappers ---
func MapperLessonServiceToServerV2(lesson *ServiceLesson) *ServerLessonV2 {
	if lesson == nil {
		return nil
	}
	return &ServerLessonV2{
		ID:          lesson.ID,
		ContractID:  lesson.ContractID,
		DurationMin: lesson.Duration,
		Format:      lesson.Format,
		CreatedAt:   lesson.CreatedAt,
	}
}

func MapperLessonCreateV2ServerToService(contractID int64, req *ServerLessonCreateV2) *ServiceLesson {
	if req == nil {
		return nil
	}
	return &ServiceLesson{
		ContractID: contractID,
		Duration:   req.DurationMin,
		CreatedAt:  time.Now(),
	}
}
