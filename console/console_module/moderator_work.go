package console_module

import (
	"data_base_project/service_logic"
	"data_base_project/types"
	"fmt"
	"strconv"
)

func PrintModeratorData(moderatorData *types.ServiceModeratorProfile) {
	fmt.Println("Moderator info:")
	fmt.Println("First name:", moderatorData.FirstName)
	fmt.Println("Last name:", moderatorData.LastName)
	fmt.Println("Middle name:", moderatorData.MiddleName)
	fmt.Println("Telephone number:", moderatorData.TelephoneNumber)
	fmt.Println("Email:", moderatorData.Email)
	fmt.Println("Salary:", moderatorData.Salary)
	for _, department := range moderatorData.Departments {
		fmt.Println("Department:", department)
	}
}

func PrintModeratorInfo(moderatorID int64, serviceModule *service_logic.ServiceModule) {
	fmt.Println("Moderator info:")
	moderatorData, err := serviceModule.ModeratorService.GetModeratorProfile(moderatorID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	PrintModeratorData(moderatorData)
}

func GetModeratorInfoByID(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter moderator ID:")
	var moderatorIDStr string
	fmt.Scanln(&moderatorIDStr)
	moderatorID, err := strconv.ParseInt(moderatorIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	PrintModeratorInfo(moderatorID, serviceModule)
}

func GetModeratorInfoByLogin(serviceModule *service_logic.ServiceModule) {
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
	if authVerdict.UserType != types.Moderator {
		fmt.Println("Error: wrong login or password")
		return
	}
	moderatorID := authVerdict.UserID
	PrintModeratorInfo(moderatorID, serviceModule)
}

func GetModeratorInfo(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Get moderator info by:")
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
		GetModeratorInfoByLogin(serviceModule)
	case 2:
		GetModeratorInfoByID(serviceModule)
	}
}

func UpdateModeratorPersonalData(serviceModule *service_logic.ServiceModule, moderatorID int64) {
	personalData := GetPersonalData()
	serviceModule.ModeratorService.UpdateModeratorPersonalData(moderatorID, personalData)
}

func UpdateModeratorPersonalDataByID(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter moderator ID:")
	var moderatorIDStr string
	fmt.Scanln(&moderatorIDStr)
	moderatorID, err := strconv.ParseInt(moderatorIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	UpdateModeratorPersonalData(serviceModule, moderatorID)
}

func UpdateModeratorPersonalDataByLogin(serviceModule *service_logic.ServiceModule) {
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
	if authVerdict.UserType != types.Moderator {
		fmt.Println("Error: wrong login or password")
		return
	}
	moderatorID := authVerdict.UserID
	UpdateModeratorPersonalData(serviceModule, moderatorID)
}

func UpdateModeratorInfo(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Update moderator info")
	fmt.Println("1. By ID")
	fmt.Println("2. By Login")
	var choiceStr string
	fmt.Scanln(&choiceStr)
	choice, err := strconv.Atoi(choiceStr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	switch choice {
	case 1:
		UpdateModeratorPersonalDataByID(serviceModule)
	case 2:
		UpdateModeratorPersonalDataByLogin(serviceModule)
	}
}

func ModeratorWork(serviceModule *service_logic.ServiceModule) {
	for {
		fmt.Println("Moderator work")
		fmt.Println("1. Get moderator info")
		fmt.Println("2. Update moderator info")
		fmt.Println("3. List all moderators")
		fmt.Println("4. Update moderator salary")
		fmt.Println("5. Exit")
		var choiceStr string
		fmt.Scanln(&choiceStr)
		choice, err := strconv.Atoi(choiceStr)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if choice < 1 || choice > 5 {
			fmt.Println("Invalid choice")
			continue
		}
		switch choice {
		case 1:
			GetModeratorInfo(serviceModule)
		case 2:
			UpdateModeratorInfo(serviceModule)
		case 3:
			ListAllModerators(serviceModule)
		case 4:
			UpdateModeratorSalary(serviceModule)
		case 5:
			return
		}
	}
}

func ListAllModerators(serviceModule *service_logic.ServiceModule) {
	moderators, err := serviceModule.ModeratorService.GetModerators()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if len(moderators) == 0 {
		fmt.Println("No moderators found")
		return
	}
	for _, m := range moderators {
		fmt.Println("--- Moderator ---")
		fmt.Println("ID:", m.ID)
		fmt.Println("First name:", m.FirstName)
		fmt.Println("Last name:", m.LastName)
		fmt.Println("Middle name:", m.MiddleName)
		fmt.Println("Telephone number:", m.TelephoneNumber)
		fmt.Println("Email:", m.Email)
		fmt.Println("Salary:", m.Salary)
	}
}

func UpdateModeratorSalary(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter moderator ID:")
	idStr := ""
	fmt.Scanln(&idStr)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Enter new salary:")
	salaryStr := ""
	fmt.Scanln(&salaryStr)
	salary, err := strconv.ParseInt(salaryStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if err := serviceModule.ModeratorService.UpdateModeratorSalary(id, salary); err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Salary updated")
}
