package service_logic

import (
	"data_base_project/data_base"
	"data_base_project/types"
	"time"
)

type IResumeService interface {
	GetResume(resumeID int64) (*types.ServiceResume, error)
	CreateResume(resume types.ServiceResume) (int64, error)
	UpdateResumeTitle(resumeID int64, title string) error
	UpdateResumeDescription(resumeID int64, description string) error
	UpdateResumePrices(resumeID int64, prices map[string]int) error
	DeleteResume(resumeID int64) error
}

type ResumeService struct {
	resumeRepository data_base.IResumeRepository
}

func CreateResumeService(resumeRepository data_base.IResumeRepository) *ResumeService {
	return &ResumeService{resumeRepository: resumeRepository}
}

func (r *ResumeService) CreateResume(resume types.ServiceResume) (int64, error) {
	return r.resumeRepository.InsertResume(*types.MapperResumeServiceToDB(&resume))
}

func (r *ResumeService) UpdateResumeTitle(resumeID int64, title string) error {
	return r.resumeRepository.UpdateResumeTitle(resumeID, title, time.Now())
}

func (r *ResumeService) UpdateResumeDescription(resumeID int64, description string) error {
	return r.resumeRepository.UpdateResumeDescription(resumeID, description, time.Now())
}

func (r *ResumeService) GetResume(resumeID int64) (*types.ServiceResume, error) {
	dbResume, err := r.resumeRepository.GetResume(resumeID)
	if err != nil {
		return nil, err
	}
	return types.MapperResumeDBToService(dbResume), nil
}

func (r *ResumeService) UpdateResumePrices(resumeID int64, prices map[string]int) error {
	return r.resumeRepository.UpdateResumePrices(resumeID, prices, time.Now())
}

func (r *ResumeService) DeleteResume(resumeID int64) error {
	return r.resumeRepository.DeleteResume(resumeID)
}

func (r *ResumeService) BuildUpInfo(repetitorID int64) (*types.ServiceRepetitorData, error) {
	return nil, nil
}
