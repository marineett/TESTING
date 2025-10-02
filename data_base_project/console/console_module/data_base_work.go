package console_module

import (
	"data_base_project/data_base"
	"data_base_project/service_logic"
	"fmt"
)

func ChangeDataBaseType(sqlDataBaseModule *data_base.DataBaseModule, mongoDataBaseModule *data_base.DataBaseModule) *service_logic.ServiceModule {
	fmt.Println("Choose data_base type:")
	fmt.Println("1. SQL")
	fmt.Println("2. MongoDB")
	fmt.Println("Enter data_base type:")
	var dataBaseType string
	fmt.Scanln(&dataBaseType)
	var serviceModule *service_logic.ServiceModule
	if dataBaseType == "1" {
		serviceModule = service_logic.CreateServiceModule(
			sqlDataBaseModule.AuthRepository,
			sqlDataBaseModule.AdminRepository,
			sqlDataBaseModule.ModeratorRepository,
			sqlDataBaseModule.ClientRepository,
			sqlDataBaseModule.RepetitorRepository,
			sqlDataBaseModule.ContractRepository,
			sqlDataBaseModule.ReviewRepository,
			sqlDataBaseModule.ChatRepository,
			sqlDataBaseModule.MessageRepository,
			sqlDataBaseModule.ResumeRepository,
			sqlDataBaseModule.TransactionRepository,
			sqlDataBaseModule.DepartmentRepository,
			sqlDataBaseModule.PersonalDataRepository,
			sqlDataBaseModule.LessonRepository,
		)
	} else if dataBaseType == "2" {
		serviceModule = service_logic.CreateServiceModule(
			mongoDataBaseModule.AuthRepository,
			mongoDataBaseModule.AdminRepository,
			mongoDataBaseModule.ModeratorRepository,
			mongoDataBaseModule.ClientRepository,
			mongoDataBaseModule.RepetitorRepository,
			mongoDataBaseModule.ContractRepository,
			mongoDataBaseModule.ReviewRepository,
			mongoDataBaseModule.ChatRepository,
			mongoDataBaseModule.MessageRepository,
			mongoDataBaseModule.ResumeRepository,
			mongoDataBaseModule.TransactionRepository,
			mongoDataBaseModule.DepartmentRepository,
			mongoDataBaseModule.PersonalDataRepository,
			mongoDataBaseModule.LessonRepository,
		)
	}
	return serviceModule
}
