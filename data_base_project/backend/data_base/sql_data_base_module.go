package data_base

type DataBaseModule struct {
	UserRepository         IUserRepository
	AuthRepository         IAuthRepository
	AdminRepository        IAdminRepository
	ModeratorRepository    IModeratorRepository
	ClientRepository       IClientRepository
	RepetitorRepository    IRepetitorRepository
	ContractRepository     IContractRepository
	ReviewRepository       IReviewRepository
	ChatRepository         IChatRepository
	MessageRepository      IMessageRepository
	ResumeRepository       IResumeRepository
	TransactionRepository  ITransactionRepository
	DepartmentRepository   IDepartmentRepository
	PersonalDataRepository IPersonalDataRepository
	LessonRepository       ILessonRepository
}

func CreateDataBaseModule(
	userRepository IUserRepository,
	authRepository IAuthRepository,
	adminRepository IAdminRepository,
	moderatorRepository IModeratorRepository,
	clientRepository IClientRepository,
	repetitorRepository IRepetitorRepository,
	contractRepository IContractRepository,
	reviewRepository IReviewRepository,
	chatRepository IChatRepository,
	messageRepository IMessageRepository,
	resumeRepository IResumeRepository,
	transactionRepository ITransactionRepository,
	departmentRepository IDepartmentRepository,
	personalDataRepository IPersonalDataRepository,
	lessonRepository ILessonRepository) *DataBaseModule {
	return &DataBaseModule{
		userRepository,
		authRepository,
		adminRepository,
		moderatorRepository,
		clientRepository,
		repetitorRepository,
		contractRepository,
		reviewRepository,
		chatRepository,
		messageRepository,
		resumeRepository,
		transactionRepository,
		departmentRepository,
		personalDataRepository,
		lessonRepository,
	}
}
