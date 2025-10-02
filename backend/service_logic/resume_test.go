package service_logic

import (
	tu "data_base_project/test_service_utility"
	"testing"
)

func TestCreateResumeCorrectLondon(t *testing.T) {
	resumeRepository := tu.CreateTestResumeRepository()
	resumeService := CreateResumeService(resumeRepository)
	resumeID, err := resumeService.CreateResume(tu.TestResume)
	if err != nil {
		t.Fatalf("Error creating resume: %v", err)
	}
	resumeServiceData, err := resumeService.GetResume(resumeID)
	if err != nil {
		t.Fatalf("Error getting resume: %v", err)
	}
	if resumeServiceData.Title != tu.TestResume.Title {
		t.Fatalf("Resume title not updated: %v", resumeServiceData)
	}
	if resumeServiceData.Description != tu.TestResume.Description {
		t.Fatalf("Resume description not updated: %v", resumeServiceData)
	}
	if len(resumeServiceData.Prices) != len(tu.TestResume.Prices) {
		t.Fatalf("Resume prices not updated: %v", resumeServiceData)
	}
	for key, value := range resumeServiceData.Prices {
		if _, ok := tu.TestResume.Prices[key]; !ok {
			t.Fatalf("Resume prices not updated: %v", resumeServiceData)
		}
		if value != tu.TestResume.Prices[key] {
			t.Fatalf("Resume prices not updated: %v", resumeServiceData)
		}
	}
}

func TestGetResumeIncorrectLondon(t *testing.T) {
	resumeRepository := tu.CreateTestResumeRepository()
	resumeService := CreateResumeService(resumeRepository)
	resumeServiceData, err := resumeService.GetResume(1)
	if err == nil {
		t.Fatalf("No error getting resume: %v", resumeServiceData)
	}
}

func TestUpdateResumeTitleCorrectLondon(t *testing.T) {
	resumeRepository := tu.CreateTestResumeRepository()
	resumeService := CreateResumeService(resumeRepository)
	resumeID, err := resumeService.CreateResume(tu.TestResume)
	if err != nil {
		t.Fatalf("Error creating resume: %v", err)
	}
	err = resumeService.UpdateResumeTitle(resumeID, "new title")
	if err != nil {
		t.Fatalf("Error updating resume title: %v", err)
	}
	resumeServiceData, err := resumeService.GetResume(resumeID)
	if err != nil {
		t.Fatalf("Error getting resume: %v", err)
	}
	if resumeServiceData.Title != "new title" {
		t.Fatalf("Resume title not updated: %v", resumeServiceData)
	}
}

func TestUpdateResumeTitleIncorrectLondon(t *testing.T) {
	resumeRepository := tu.CreateTestResumeRepository()
	resumeService := CreateResumeService(resumeRepository)
	err := resumeService.UpdateResumeTitle(1, "new title")
	if err == nil {
		t.Fatalf("No error updating resume title: %v", err)
	}
}

func TestUpdateResumeDescriptionCorrectLondon(t *testing.T) {
	resumeRepository := tu.CreateTestResumeRepository()
	resumeService := CreateResumeService(resumeRepository)
	resumeID, err := resumeService.CreateResume(tu.TestResume)
	if err != nil {
		t.Fatalf("Error creating resume: %v", err)
	}
	err = resumeService.UpdateResumeDescription(resumeID, "new description")
	if err != nil {
		t.Fatalf("Error updating resume description: %v", err)
	}
	resumeServiceData, err := resumeService.GetResume(resumeID)
	if err != nil {
		t.Fatalf("Error getting resume: %v", err)
	}
	if resumeServiceData.Description != "new description" {
		t.Fatalf("Resume description not updated: %v", resumeServiceData)
	}
}

func TestUpdateResumeDescriptionIncorrectLondon(t *testing.T) {
	resumeRepository := tu.CreateTestResumeRepository()
	resumeService := CreateResumeService(resumeRepository)
	err := resumeService.UpdateResumeDescription(1, "new description")
	if err == nil {
		t.Fatalf("No error updating resume description: %v", err)
	}
}

func TestUpdateResumePricesCorrectLondon(t *testing.T) {
	resumeRepository := tu.CreateTestResumeRepository()
	resumeService := CreateResumeService(resumeRepository)
	resumeID, err := resumeService.CreateResume(tu.TestResume)
	if err != nil {
		t.Fatalf("Error creating resume: %v", err)
	}
	mapPrices := map[string]int{"Java": 400, "C++": 300}
	err = resumeService.UpdateResumePrices(resumeID, mapPrices)
	if err != nil {
		t.Fatalf("Error updating resume prices: %v", err)
	}
	resumeServiceData, err := resumeService.GetResume(resumeID)
	if err != nil {
		t.Fatalf("Error getting resume: %v", err)
	}
	if len(resumeServiceData.Prices) != len(mapPrices) {
		t.Fatalf("Resume prices not updated: %v", resumeServiceData)
	}
	for key, value := range resumeServiceData.Prices {
		if _, ok := mapPrices[key]; !ok {
			t.Fatalf("Resume prices not updated: %v", resumeServiceData)
		}
		if value != mapPrices[key] {
			t.Fatalf("Resume prices not updated: %v", resumeServiceData)
		}
	}
}

func TestUpdateResumePricesIncorrectLondon(t *testing.T) {
	resumeRepository := tu.CreateTestResumeRepository()
	resumeService := CreateResumeService(resumeRepository)
	err := resumeService.UpdateResumePrices(1, map[string]int{"Go": 100, "C++": 200, "Python": 100})
	if err == nil {
		t.Fatalf("No error updating resume prices: %v", err)
	}
}

func TestDeleteResumeCorrectLondon(t *testing.T) {
	resumeRepository := tu.CreateTestResumeRepository()
	resumeService := CreateResumeService(resumeRepository)
	resumeID, err := resumeService.CreateResume(tu.TestResume)
	if err != nil {
		t.Fatalf("Error creating resume: %v", err)
	}
	err = resumeService.DeleteResume(resumeID)
	if err != nil {
		t.Fatalf("Error deleting resume: %v", err)
	}
	resumeServiceData, err := resumeService.GetResume(resumeID)
	if err == nil {
		t.Fatalf("No error getting resume after deletion: %v", resumeServiceData)
	}
}

func TestDeleteResumeIncorrectLondon(t *testing.T) {
	resumeRepository := tu.CreateTestResumeRepository()
	resumeService := CreateResumeService(resumeRepository)
	err := resumeService.DeleteResume(1)
	if err == nil {
		t.Fatalf("No error deleting resume: %v", err)
	}
}
