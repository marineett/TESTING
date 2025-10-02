package console_module

import (
	"data_base_project/service_logic"
	"data_base_project/types"
	"fmt"
	"strconv"
)

func PrintClient(clientData *types.ServiceClientProfile) {
	fmt.Println("Client info:")
	fmt.Println("First name:", clientData.FirstName)
	fmt.Println("Last name:", clientData.LastName)
	fmt.Println("Middle name:", clientData.MiddleName)
	fmt.Println("Telephone number:", clientData.TelephoneNumber)
	fmt.Println("Email:", clientData.Email)
	fmt.Println("Mean rating:", clientData.MeanRating)
	for _, review := range clientData.Reviews {
		PrintReview(&review)
	}
}

func PrintClientInfo(clientID int64, serviceModule *service_logic.ServiceModule) {
	fmt.Println("Client info:")
	clientData, err := serviceModule.ClientService.GetClientProfile(clientID, 0, 100)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	PrintClient(clientData)
}

func GetClientInfoByID(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter client ID:")
	clientIDStr := ""
	fmt.Scanln(&clientIDStr)
	clientID, err := strconv.ParseInt(clientIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Print(clientID)
	PrintClientInfo(clientID, serviceModule)
}

func GetClientInfoByLogin(serviceModule *service_logic.ServiceModule) {
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
	if authVerdict.UserType != types.Client {
		fmt.Println("Error: wrong login or password")
		return
	}
	clientID := authVerdict.UserID
	PrintClientInfo(clientID, serviceModule)
}

func GetClientInfo(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Get client info by:")
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
		GetClientInfoByLogin(serviceModule)
	case 2:
		GetClientInfoByID(serviceModule)
	}
}

func UpdateClientPersonalData(serviceModule *service_logic.ServiceModule, clientID int64) {
	personalData := GetPersonalData()
	serviceModule.ClientService.UpdateClientPersonalData(clientID, personalData)
}

func UpdateClientPersonalDataByID(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter client ID:")
	clientIDStr := ""
	fmt.Scanln(&clientIDStr)
	clientID, err := strconv.ParseInt(clientIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	UpdateClientPersonalData(serviceModule, clientID)
}

func UpdateClientPersonalDataByLogin(serviceModule *service_logic.ServiceModule) {
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
	if authVerdict.UserType != types.Client {
		fmt.Println("Error: wrong login or password")
		return
	}
	clientID := authVerdict.UserID
	UpdateClientPersonalData(serviceModule, clientID)
}

func UpdateClientInfo(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Update client info")
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
		UpdateClientPersonalDataByID(serviceModule)
	case 2:
		UpdateClientPersonalDataByLogin(serviceModule)
	}
}

func ClientWork(serviceModule *service_logic.ServiceModule) {
	for {
		fmt.Println("Client work")
		fmt.Println("1. Get client info")
		fmt.Println("2. Update client info")
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
			GetClientInfo(serviceModule)
		case 2:
			UpdateClientInfo(serviceModule)
		case 3:
			return
		}
	}
}
