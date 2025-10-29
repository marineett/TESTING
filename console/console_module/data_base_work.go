package console_module

import (
	"data_base_project/data_base"
	"data_base_project/service_logic"
)

func ChangeDataBaseType(sqlDataBaseModule *data_base.DataBaseModule) *service_logic.ServiceModule {
	return service_logic.CreateServiceModule(
		sqlDataBaseModule.UserRepository,
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
}
