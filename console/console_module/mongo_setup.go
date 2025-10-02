package console_module

import (
	"data_base_project/data_base"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
)

func MongoSetup(client *mongo.Client) *data_base.DataBaseModule {

	authRepository := data_base.CreateMongoAuthRepository(
		client,
		os.Getenv("DATABASE_NAME"),
		os.Getenv("AUTH_TABLE_NAME"),
		os.Getenv("MONGO_IDS_TABLE"),
	)
	adminRepository := data_base.CreateMongoAdminRepository(
		client,
		os.Getenv("DATABASE_NAME"),
		os.Getenv("PERSONAL_DATA_TABLE_NAME"),
		os.Getenv("USER_TABLE_NAME"),
		os.Getenv("ADMIN_TABLE_NAME"),
		os.Getenv("AUTH_TABLE_NAME"),
		os.Getenv("MONGO_IDS_TABLE"),
	)
	moderatorRepository := data_base.CreateMongoModeratorRepository(
		client,
		os.Getenv("DATABASE_NAME"),
		os.Getenv("PERSONAL_DATA_TABLE_NAME"),
		os.Getenv("USER_TABLE_NAME"),
		os.Getenv("MODERATOR_TABLE_NAME"),
		os.Getenv("AUTH_TABLE_NAME"),
		os.Getenv("MONGO_IDS_TABLE"),
	)
	clientRepository := data_base.CreateMongoClientRepository(
		client,
		os.Getenv("DATABASE_NAME"),
		os.Getenv("PERSONAL_DATA_TABLE_NAME"),
		os.Getenv("USER_TABLE_NAME"),
		os.Getenv("CLIENT_TABLE_NAME"),
		os.Getenv("AUTH_TABLE_NAME"),
		os.Getenv("MONGO_IDS_TABLE"),
	)
	repetitorRepository := data_base.CreateMongoRepetitorRepository(
		client,
		os.Getenv("DATABASE_NAME"),
		os.Getenv("PERSONAL_DATA_TABLE_NAME"),
		os.Getenv("USER_TABLE_NAME"),
		os.Getenv("REPEATITOR_TABLE_NAME"),
		os.Getenv("AUTH_TABLE_NAME"),
		os.Getenv("RESUME_TABLE_NAME"),
		os.Getenv("REVIEW_TABLE_NAME"),
		os.Getenv("MONGO_IDS_TABLE"),
	)
	chatRepository := data_base.CreateMongoChatRepository(
		client,
		os.Getenv("DATABASE_NAME"),
		os.Getenv("CHAT_TABLE_NAME"),
		os.Getenv("MONGO_IDS_TABLE"),
	)
	messageRepository := data_base.CreateMongoMessageRepository(
		client,
		os.Getenv("DATABASE_NAME"),
		os.Getenv("MESSAGE_TABLE_NAME"),
		os.Getenv("MONGO_IDS_TABLE"),
	)
	resumeRepository := data_base.CreateMongoResumeRepository(
		client,
		os.Getenv("DATABASE_NAME"),
		os.Getenv("RESUME_TABLE_NAME"),
		os.Getenv("MONGO_IDS_TABLE"),
	)
	reviewRepository := data_base.CreateMongoReviewRepository(
		client,
		os.Getenv("DATABASE_NAME"),
		os.Getenv("REVIEW_TABLE_NAME"),
		os.Getenv("MONGO_IDS_TABLE"),
	)
	contractRepository := data_base.CreateMongoContractRepository(
		client,
		os.Getenv("DATABASE_NAME"),
		os.Getenv("CONTRACT_TABLE_NAME"),
		os.Getenv("MONGO_IDS_TABLE"),
	)
	transactionRepository := data_base.CreateMongoTransactionRepository(
		client,
		os.Getenv("DATABASE_NAME"),
		os.Getenv("TRANSACTION_TABLE_NAME"),
		os.Getenv("PENDING_CONTRACT_PAYMENT_TRANSACTION_TABLE_NAME"),
		os.Getenv("MONGO_IDS_TABLE"),
	)
	departmentRepository := data_base.CreateMongoDepartmentRepository(
		client,
		os.Getenv("DATABASE_NAME"),
		os.Getenv("DEPARTMENT_TABLE_NAME"),
		os.Getenv("HIRE_INFO_TABLE_NAME"),
		os.Getenv("MONGO_IDS_TABLE"),
	)
	personalDataRepository := data_base.CreateMongoPersonalDataRepository(
		client,
		os.Getenv("DATABASE_NAME"),
		os.Getenv("PERSONAL_DATA_TABLE_NAME"),
		os.Getenv("MONGO_IDS_TABLE"),
	)
	lessonRepository := data_base.CreateMongoLessonRepository(
		client,
		os.Getenv("DATABASE_NAME"),
		os.Getenv("LESSON_TABLE_NAME"),
		os.Getenv("CONTRACT_TABLE_NAME"),
		os.Getenv("TRANSACTION_TABLE_NAME"),
		os.Getenv("MONGO_IDS_TABLE"),
	)

	return data_base.CreateDataBaseModule(
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
	)
}
