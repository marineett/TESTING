package service_logic

import (
	"data_base_project/data_base"
	"data_base_project/types"
)

type ILessonService interface {
	CreateLesson(lesson types.ServiceLesson) (int64, error)
	GetLessons(contractID int64, from int64, size int64) ([]types.ServiceLesson, error)
	GetLesson(lessonID int64) (*types.ServiceLesson, error)
}

type LessonService struct {
	lessonRepository data_base.ILessonRepository
}

func CreateLessonService(lessonRepository data_base.ILessonRepository) ILessonService {
	return &LessonService{lessonRepository: lessonRepository}
}

func (s *LessonService) CreateLesson(lesson types.ServiceLesson) (int64, error) {
	return s.lessonRepository.InsertLesson(types.DBLesson{
		ContractID: lesson.ContractID,
		Duration:   lesson.Duration,
		CreatedAt:  lesson.CreatedAt,
	})
}

func (s *LessonService) GetLessons(contractID int64, from int64, size int64) ([]types.ServiceLesson, error) {
	lessons, err := s.lessonRepository.GetLessons(contractID, from, size)
	if err != nil {
		return nil, err
	}
	serviceLessons := make([]types.ServiceLesson, len(lessons))
	for i, lesson := range lessons {
		serviceLessons[i] = *types.MapperLessonDBToService(&lesson)
	}
	return serviceLessons, nil
}

func (s *LessonService) GetLesson(lessonID int64) (*types.ServiceLesson, error) {
	lesson, err := s.lessonRepository.GetLesson(lessonID)
	if err != nil {
		return nil, err
	}
	return types.MapperLessonDBToService(lesson), nil
}
