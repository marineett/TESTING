package service_logic

import "data_base_project/data_base"

type ServiceModule struct {
	AuthService         IAuthService
	AdminService        IAdminService
	ModeratorService    IModeratorService
	ClientService       IClientService
	RepetitorService    IRepetitorService
	ContractService     IContractService
	ChatService         IChatService
	ResumeService       IResumeService
	TransactionService  ITransactionService
	DepartmentService   IDepartmentService
	PersonalDataService IPersonalDataService
	ReviewService       IReviewService
	LessonService       ILessonService
}

func CreateServiceModule(
	userRepository data_base.IUserRepository,
	authRepository data_base.IAuthRepository,
	adminRepository data_base.IAdminRepository,
	moderatorRepository data_base.IModeratorRepository,
	clientRepository data_base.IClientRepository,
	repetitorRepository data_base.IRepetitorRepository,
	contractRepository data_base.IContractRepository,
	reviewRepository data_base.IReviewRepository,
	chatRepository data_base.IChatRepository,
	messageRepository data_base.IMessageRepository,
	resumeRepository data_base.IResumeRepository,
	transactionRepository data_base.ITransactionRepository,
	departmentRepository data_base.IDepartmentRepository,
	personalDataRepository data_base.IPersonalDataRepository,
	lessonRepository data_base.ILessonRepository,
) *ServiceModule {
	return &ServiceModule{
		AuthService:         CreateAuthService(authRepository),
		AdminService:        CreateAdminService(adminRepository, userRepository, personalDataRepository),
		ModeratorService:    CreateModeratorService(moderatorRepository, personalDataRepository, userRepository, departmentRepository),
		ClientService:       CreateClientService(clientRepository, personalDataRepository, userRepository, reviewRepository),
		RepetitorService:    CreateRepetitorService(repetitorRepository, personalDataRepository, userRepository, reviewRepository, resumeRepository),
		ContractService:     CreateContractService(contractRepository, reviewRepository),
		ChatService:         CreateChatService(chatRepository, messageRepository),
		ResumeService:       CreateResumeService(resumeRepository),
		TransactionService:  CreateTransactionService(transactionRepository),
		DepartmentService:   CreateDepartmentService(departmentRepository, moderatorRepository, userRepository, personalDataRepository),
		PersonalDataService: CreatePersonalDataService(personalDataRepository),
		ReviewService:       CreateReviewService(reviewRepository),
		LessonService:       CreateLessonService(lessonRepository),
	}
}
