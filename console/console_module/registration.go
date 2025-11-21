package console_module

import (
	"bufio"
	"data_base_project/service_logic"
	"data_base_project/types"
	"fmt"
	"os"
	"strconv"
	"time"
)

func GetPersonalData() types.ServicePersonalData {
	personalData := types.ServicePersonalData{}
	fmt.Println("Enter first name:")
	fmt.Scanln(&personalData.FirstName)
	fmt.Println("Enter last name:")
	fmt.Scanln(&personalData.LastName)
	fmt.Println("Enter middle name:")
	fmt.Scanln(&personalData.MiddleName)
	fmt.Println("Enter email:")
	fmt.Scanln(&personalData.Email)
	fmt.Println("Enter telephone number:")
	fmt.Scanln(&personalData.TelephoneNumber)
	passportData := types.ServicePassportData{}
	fmt.Println("Enter passport number:")
	fmt.Scanln(&passportData.PassportNumber)
	fmt.Println("Enter passport series:")
	fmt.Scanln(&passportData.PassportSeries)
	passportData.PassportDate = time.Now()
	fmt.Println("Enter passport issued by:")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input := scanner.Text()
		passportData.PassportIssuedBy = input
	}
	personalData.PassportNumber = passportData.PassportNumber
	personalData.PassportSeries = passportData.PassportSeries
	personalData.PassportDate = passportData.PassportDate
	personalData.PassportIssuedBy = passportData.PassportIssuedBy
	return personalData
}

func GetAuthData() types.ServiceAuthData {
	authData := types.ServiceAuthData{}
	fmt.Println("Enter login:")
	fmt.Scanln(&authData.Login)
	fmt.Println("Enter password:")
	fmt.Scanln(&authData.Password)
	return authData
}

func RegisterClient(serviceModule *service_logic.ServiceModule) {
	personalData := GetPersonalData()
	authData := GetAuthData()
	initData := types.ServiceInitUserData{
		ServicePersonalData: personalData,
		ServiceAuthData:     authData,
	}
	err := serviceModule.ClientService.CreateClient(types.ServiceInitClientData{
		ServiceInitUserData: initData,
	})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	verdict, err := serviceModule.AuthService.Authorize(types.ServiceAuthData{Login: authData.Login, Password: authData.Password})
	if err == nil {
		fmt.Println("Client registered successfully with ID:", verdict.UserID)
	} else {
		fmt.Println("Client registered successfully")
	}
}

func RegisterRepetitor(serviceModule *service_logic.ServiceModule) {
	personalData := GetPersonalData()
	authData := GetAuthData()
	initData := types.ServiceInitUserData{
		ServicePersonalData: personalData,
		ServiceAuthData:     authData,
	}
	err := serviceModule.RepetitorService.CreateRepetitor(types.ServiceInitRepetitorData{
		ServiceInitUserData: initData,
	})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	verdict, err := serviceModule.AuthService.Authorize(types.ServiceAuthData{Login: authData.Login, Password: authData.Password})
	if err == nil {
		fmt.Println("Repetitor registered successfully with ID:", verdict.UserID)
	} else {
		fmt.Println("Repetitor registered successfully")
	}
}

func RegisterModerator(serviceModule *service_logic.ServiceModule) {
	personalData := GetPersonalData()
	authData := GetAuthData()
	initData := types.ServiceInitUserData{
		ServicePersonalData: personalData,
		ServiceAuthData:     authData,
	}
	err := serviceModule.ModeratorService.CreateModerator(types.ServiceInitModeratorData{
		ServiceInitUserData: initData,
	})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	verdict, err := serviceModule.AuthService.Authorize(types.ServiceAuthData{Login: authData.Login, Password: authData.Password})
	if err == nil {
		fmt.Println("Moderator registered successfully with ID:", verdict.UserID)
	} else {
		fmt.Println("Moderator registered successfully")
	}
}

func RegisterAdmin(serviceModule *service_logic.ServiceModule) {
	personalData := GetPersonalData()
	authData := GetAuthData()
	initData := types.ServiceInitUserData{
		ServicePersonalData: personalData,
		ServiceAuthData:     authData,
	}
	err := serviceModule.AdminService.CreateAdmin(types.ServiceInitAdminData{
		ServiceInitUserData: initData,
		Salary:              0,
	})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	verdict, err := serviceModule.AuthService.Authorize(types.ServiceAuthData{Login: authData.Login, Password: authData.Password})
	if err == nil {
		fmt.Println("Admin registered successfully with ID:", verdict.UserID)
	} else {
		fmt.Println("Admin registered successfully")
	}
}

func PrintRegistrationMenu() {
	fmt.Println("Registration:")
	fmt.Println("1. Register as a client")
	fmt.Println("2. Register as a repetitor")
	fmt.Println("3. Register as a moderator")
	fmt.Println("4. Register as an admin")
}

func UserRegistrationWork(serviceModule *service_logic.ServiceModule) {
	PrintRegistrationMenu()
	var choiceStr string
	fmt.Scanln(&choiceStr)
	choice, err := strconv.Atoi(choiceStr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if choice < 1 || choice > 4 {
		fmt.Println("Invalid choice")
		return
	}
	switch choice {
	case 1:
		RegisterClient(serviceModule)
	case 2:
		RegisterRepetitor(serviceModule)
	case 3:
		RegisterModerator(serviceModule)
	case 4:
		RegisterAdmin(serviceModule)
	}
}
