package console_module

import (
	"data_base_project/service_logic"
	"data_base_project/types"
	"fmt"
	"strconv"
)

func PrintAdminData(adminData *types.ServiceAdminProfile) {
	fmt.Println("Admin info:")
	fmt.Println("First name:", adminData.FirstName)
	fmt.Println("Last name:", adminData.LastName)
	fmt.Println("Middle name:", adminData.MiddleName)
	fmt.Println("Telephone number:", adminData.TelephoneNumber)
	fmt.Println("Email:", adminData.Email)
	fmt.Println("Salary:", adminData.Salary)
}

func PrintAdminInfo(adminID int64, serviceModule *service_logic.ServiceModule) {
	fmt.Println("Admin info:")
	adminData, err := serviceModule.AdminService.GetAdminProfile(adminID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	PrintAdminData(adminData)
}

func GetAdminInfoByID(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter admin ID:")

	adminIDStr := ""
	fmt.Scanln(&adminIDStr)

	adminID, err := strconv.ParseInt(adminIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	PrintAdminInfo(adminID, serviceModule)
}

func GetAdminInfoByLogin(serviceModule *service_logic.ServiceModule) {
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
	if authVerdict.UserType != types.Admin {
		fmt.Println("Error: wrong login or password")
		return
	}
	adminID := authVerdict.UserID
	PrintAdminInfo(adminID, serviceModule)
}

func GetAdminInfo(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Get admin info by:")
	fmt.Println("1. Login")
	fmt.Println("2. ID")
	choiceStr := ""
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
		GetAdminInfoByLogin(serviceModule)
	case 2:
		GetAdminInfoByID(serviceModule)
	}
}

func UpdateAdminPersonalDataMenu(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Update admin personal data")
	fmt.Println("1. By ID")
	fmt.Println("2. By Login")
	choiceStr := ""
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
		UpdateAdminPersonalDataByID(serviceModule)
	case 2:
		UpdateAdminPersonalDataByLogin(serviceModule)
	}
}

func UpdateAdminSalaryByID(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter admin ID:")
	adminIDStr := ""
	fmt.Scanln(&adminIDStr)
	adminID, err := strconv.ParseInt(adminIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Enter salary:")
	salaryStr := ""
	fmt.Scanln(&salaryStr)
	salary, err := strconv.ParseInt(salaryStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if salary < 0 {
		fmt.Println("Salary cannot be negative")
		return
	}
	err = serviceModule.AdminService.UpdateAdminSalary(adminID, salary)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func UpdateAdminSalaryByLogin(serviceModule *service_logic.ServiceModule) {
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
	fmt.Println("Enter salary:")
	salaryStr := ""
	fmt.Scanln(&salaryStr)
	salary, err := strconv.ParseInt(salaryStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if salary < 0 {
		fmt.Println("Salary cannot be negative")
		return
	}
	err = serviceModule.AdminService.UpdateAdminSalary(authVerdict.UserID, salary)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func UpdateAdminSalaryMenu(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Update admin salary")
	fmt.Println("1. By ID")
	fmt.Println("2. By Login")
	var choice int
	fmt.Scanln(&choice)
	switch choice {
	case 1:
		UpdateAdminSalaryByID(serviceModule)
	case 2:
		UpdateAdminSalaryByLogin(serviceModule)
	}
}

func UpdateAdminDepartmentByID(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter admin ID:")
	adminIDStr := ""
	fmt.Scanln(&adminIDStr)
	adminID, err := strconv.ParseInt(adminIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Enter department ID:")
	departmentIDStr := ""
	fmt.Scanln(&departmentIDStr)
	departmentID, err := strconv.ParseInt(departmentIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	err = serviceModule.AdminService.UpdateAdminDepartment(adminID, departmentID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func UpdateAdminDepartmentByLogin(serviceModule *service_logic.ServiceModule) {
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
	fmt.Println("Enter department ID:")
	departmentIDStr := ""
	fmt.Scanln(&departmentIDStr)
	departmentID, err := strconv.ParseInt(departmentIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	err = serviceModule.AdminService.UpdateAdminDepartment(authVerdict.UserID, departmentID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func UpdateAdminDepartmentMenu(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Update admin department")
	fmt.Println("1. By ID")
	fmt.Println("2. By Login")
	choiceStr := ""
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
		UpdateAdminDepartmentByID(serviceModule)
	case 2:
		UpdateAdminDepartmentByLogin(serviceModule)
	}
}

func UpdateAdminPersonalData(serviceModule *service_logic.ServiceModule, adminID int64) {
	personalData := GetPersonalData()
	err := serviceModule.AdminService.UpdateAdminPersonalData(adminID, personalData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func UpdateAdminPersonalDataByID(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter admin ID:")
	adminIDStr := ""
	fmt.Scanln(&adminIDStr)
	adminID, err := strconv.ParseInt(adminIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	UpdateAdminPersonalData(serviceModule, adminID)
}

func UpdateAdminPersonalDataByLogin(serviceModule *service_logic.ServiceModule) {
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
	if authVerdict.UserType != types.Admin {
		fmt.Println("Error: wrong login or password")
		return
	}
	adminID := authVerdict.UserID
	UpdateAdminPersonalData(serviceModule, adminID)
}

func UpdateAdminInfo(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Update admin info")
	fmt.Println("1. Update personal data")
	fmt.Println("2. Update department")
	fmt.Println("3. Update salary")
	choiceStr := ""
	fmt.Scanln(&choiceStr)
	choice, err := strconv.Atoi(choiceStr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if choice < 1 || choice > 3 {
		fmt.Println("Invalid choice")
		return
	}
	switch choice {
	case 1:
		UpdateAdminPersonalDataMenu(serviceModule)
	case 2:
		UpdateAdminDepartmentMenu(serviceModule)
	case 3:
		UpdateAdminSalaryMenu(serviceModule)
	}
}

func AdminWork(serviceModule *service_logic.ServiceModule) {
	for {
		fmt.Println("Admin work")
		fmt.Println("1. Get admin info")
		fmt.Println("2. Update admin info")
		fmt.Println("3. Exit")
		choiceStr := ""
		fmt.Scanln(&choiceStr)
		choice, err := strconv.Atoi(choiceStr)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if choice < 1 || choice > 3 {
			fmt.Println("Invalid choice")
			return
		}
		switch choice {
		case 1:
			GetAdminInfo(serviceModule)
		case 2:
			UpdateAdminInfo(serviceModule)
		case 3:
			return
		}
	}
}
