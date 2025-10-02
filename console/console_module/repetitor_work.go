package console_module

import (
	"data_base_project/service_logic"
	"data_base_project/types"
	"fmt"
	"strconv"
)

func PrintRepetitorData(repetitorData *types.ServiceRepetitorProfile) {
	fmt.Println("Repetitor info:")
	fmt.Println("First name:", repetitorData.FirstName)
	fmt.Println("Last name:", repetitorData.LastName)
	fmt.Println("Middle name:", repetitorData.MiddleName)
	fmt.Println("Telephone number:", repetitorData.TelephoneNumber)
	fmt.Println("Email:", repetitorData.Email)
	fmt.Println("Mean rating:", repetitorData.MeanRating)
	for _, review := range repetitorData.Reviews {
		PrintReview(&review)
	}
}

func PrintRepetitorInfo(repetitorID int64, serviceModule *service_logic.ServiceModule) {
	fmt.Println("Repetitor info:")
	repetitorData, err := serviceModule.RepetitorService.GetRepetitorProfile(repetitorID, 0, 100)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	PrintRepetitorData(repetitorData)
}

func GetRepetitorInfoByID(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter repetitor ID:")
	var repetitorIDStr string
	fmt.Scanln(&repetitorIDStr)
	repetitorID, err := strconv.ParseInt(repetitorIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	PrintRepetitorInfo(repetitorID, serviceModule)
}

func GetRepetitorInfoByLogin(serviceModule *service_logic.ServiceModule) {
	var authData types.ServiceAuthData
	fmt.Println("Enter login:")
	fmt.Scanln(&authData.Login)
	fmt.Println("Enter password:")
	fmt.Scanln(&authData.Password)
	authVerdict, err := serviceModule.AuthService.Authorize(authData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if authVerdict.UserType != types.Repetitor {
		fmt.Println("Error: wrong login or password")
		return
	}
	repetitorID := authVerdict.UserID
	PrintRepetitorInfo(repetitorID, serviceModule)
}

func GetRepetitorInfo(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Get repetitor info by:")
	fmt.Println("1. Login")
	fmt.Println("2. ID")
	var choiceStr string
	fmt.Scanln(&choiceStr)
	choice, err := strconv.Atoi(choiceStr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if choice < 1 || choice > 2 {
		fmt.Println("Invalid choice")
		return
	}
	switch choice {
	case 1:
		GetRepetitorInfoByLogin(serviceModule)
	case 2:
		GetRepetitorInfoByID(serviceModule)
	}
}

func UpdateRepetitorPersonalData(serviceModule *service_logic.ServiceModule, repetitorID int64) {
	personalData := GetPersonalData()
	serviceModule.RepetitorService.UpdateRepetitorPersonalData(repetitorID, personalData)
}

func UpdateRepetitorPersonalDataByID(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter repetitor ID:")
	var repetitorIDStr string
	fmt.Scanln(&repetitorIDStr)
	repetitorID, err := strconv.ParseInt(repetitorIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	UpdateRepetitorPersonalData(serviceModule, repetitorID)
}

func UpdateRepetitorPersonalDataByLogin(serviceModule *service_logic.ServiceModule) {
	var authData types.ServiceAuthData
	fmt.Println("Enter login:")
	fmt.Scanln(&authData.Login)
	fmt.Println("Enter password:")
	fmt.Scanln(&authData.Password)
	authVerdict, err := serviceModule.AuthService.Authorize(authData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if authVerdict.UserType != types.Repetitor {
		fmt.Println("Error: wrong login or password")
		return
	}
	repetitorID := authVerdict.UserID
	UpdateRepetitorPersonalData(serviceModule, repetitorID)
}

func UpdateRepetitorInfo(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Update repetitor info")
	fmt.Println("1. By ID")
	fmt.Println("2. By Login")
	var choiceStr string
	fmt.Scanln(&choiceStr)
	choice, err := strconv.Atoi(choiceStr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if choice < 1 || choice > 2 {
		fmt.Println("Invalid choice")
		return
	}
	switch choice {
	case 1:
		UpdateRepetitorPersonalDataByID(serviceModule)
	case 2:
		UpdateRepetitorPersonalDataByLogin(serviceModule)
	}
}

func RepetitorWork(serviceModule *service_logic.ServiceModule) {
	for {
		fmt.Println("Repetitor work")
		fmt.Println("1. Get repetitor info")
		fmt.Println("2. Update repetitor info")
		fmt.Println("3. Exit")
		var choiceStr string
		fmt.Scanln(&choiceStr)
		choice, err := strconv.Atoi(choiceStr)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		switch choice {
		case 1:
			GetRepetitorInfo(serviceModule)
		case 2:
			UpdateRepetitorInfo(serviceModule)
		case 3:
			return
		}
	}
}
