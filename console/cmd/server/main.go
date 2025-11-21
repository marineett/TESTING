package main

import (
	"console/server_http"
	"data_base_project/data_base"
	"data_base_project/service_logic"
	"data_base_project/utility_module"
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	defer utility_module.UnsetEnv()

	db, err := data_base.CreateSqlConnection(data_base.GetSqlConnectionString())
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	ensureTables(db)

	services := buildServices(db)

	if os.Getenv("CONSOLE_SERVER_PORT") == "" {
		os.Setenv("CONSOLE_SERVER_PORT", "8081")
	}
	_ = server_http.StartHTTPServer(services, log.Default())

	select {}
}

func ensureTables(db *sql.DB) {
	if err := data_base.CreateSqlTables(
		db,
		os.Getenv("PERSONAL_DATA_TABLE_NAME"),
		os.Getenv("USER_TABLE_NAME"),
		os.Getenv("AUTH_TABLE_NAME"),
		os.Getenv("CHAT_TABLE_NAME"),
		os.Getenv("MESSAGE_TABLE_NAME"),
		os.Getenv("DEPARTMENT_TABLE_NAME"),
		os.Getenv("HIRE_INFO_TABLE_NAME"),
		os.Getenv("CLIENT_TABLE_NAME"),
		os.Getenv("RESUME_TABLE_NAME"),
		os.Getenv("REVIEW_TABLE_NAME"),
		os.Getenv("REPEATITOR_TABLE_NAME"),
		os.Getenv("CONTRACT_TABLE_NAME"),
		os.Getenv("ADMIN_TABLE_NAME"),
		os.Getenv("MODERATOR_TABLE_NAME"),
		os.Getenv("TRANSACTION_TABLE_NAME"),
		os.Getenv("PENDING_CONTRACT_PAYMENT_TRANSACTIONS"),
		os.Getenv("LESSON_TABLE_NAME"),
		os.Getenv("SEQUENCE_NAME"),
	); err != nil {
		log.Fatalf("Error creating tables: %v", err)
	}
}

func buildServices(db *sql.DB) *service_logic.ServiceModule {
	userRepository := data_base.CreateSqlUserRepository(
		db,
		os.Getenv("USER_TABLE_NAME"),
		os.Getenv("SEQUENCE_NAME"),
	)
	adminRepository := data_base.CreateSqlAdminRepository(
		db,
		os.Getenv("PERSONAL_DATA_TABLE_NAME"),
		os.Getenv("USER_TABLE_NAME"),
		os.Getenv("ADMIN_TABLE_NAME"),
		os.Getenv("AUTH_TABLE_NAME"),
		os.Getenv("SEQUENCE_NAME"),
	)
	moderatorRepository := data_base.CreateSqlModeratorRepository(
		db,
		os.Getenv("PERSONAL_DATA_TABLE_NAME"),
		os.Getenv("USER_TABLE_NAME"),
		os.Getenv("MODERATOR_TABLE_NAME"),
		os.Getenv("AUTH_TABLE_NAME"),
		os.Getenv("SEQUENCE_NAME"),
	)
	clientRepository := data_base.CreateSqlClientRepository(
		db,
		os.Getenv("PERSONAL_DATA_TABLE_NAME"),
		os.Getenv("USER_TABLE_NAME"),
		os.Getenv("CLIENT_TABLE_NAME"),
		os.Getenv("AUTH_TABLE_NAME"),
		os.Getenv("SEQUENCE_NAME"),
	)
	repetitorRepository := data_base.CreateSqlRepetitorRepository(
		db,
		os.Getenv("PERSONAL_DATA_TABLE_NAME"),
		os.Getenv("USER_TABLE_NAME"),
		os.Getenv("REPEATITOR_TABLE_NAME"),
		os.Getenv("AUTH_TABLE_NAME"),
		os.Getenv("RESUME_TABLE_NAME"),
		os.Getenv("REVIEW_TABLE_NAME"),
		os.Getenv("SEQUENCE_NAME"),
	)
	contractRepository := data_base.CreateSqlContractRepository(
		db,
		os.Getenv("CONTRACT_TABLE_NAME"),
		os.Getenv("SEQUENCE_NAME"),
	)
	chatRepository := data_base.CreateSqlChatRepository(
		db,
		os.Getenv("CHAT_TABLE_NAME"),
		os.Getenv("SEQUENCE_NAME"),
	)
	messageRepository := data_base.CreateSqlMessageRepository(
		db,
		os.Getenv("MESSAGE_TABLE_NAME"),
		os.Getenv("SEQUENCE_NAME"),
	)
	resumeRepository := data_base.CreateSqlResumeRepository(
		db,
		os.Getenv("RESUME_TABLE_NAME"),
		os.Getenv("SEQUENCE_NAME"),
	)
	transactionRepository := data_base.CreateSqlTransactionRepository(
		db,
		os.Getenv("TRANSACTION_TABLE_NAME"),
		os.Getenv("PENDING_CONTRACT_PAYMENT_TRANSACTIONS"),
		os.Getenv("SEQUENCE_NAME"),
	)
	reviewRepository := data_base.CreateSqlReviewRepository(
		db,
		os.Getenv("REVIEW_TABLE_NAME"),
		os.Getenv("SEQUENCE_NAME"),
	)
	authRepository := data_base.CreateSqlAuthRepository(
		db,
		os.Getenv("AUTH_TABLE_NAME"),
		os.Getenv("SEQUENCE_NAME"),
	)
	departmentRepository := data_base.CreateSqlDepartmentRepository(
		db,
		os.Getenv("DEPARTMENT_TABLE_NAME"),
		os.Getenv("HIRE_INFO_TABLE_NAME"),
		os.Getenv("SEQUENCE_NAME"),
	)
	personalDataRepository := data_base.CreateSqlPersonalDataRepository(
		db,
		os.Getenv("PERSONAL_DATA_TABLE_NAME"),
		os.Getenv("SEQUENCE_NAME"),
	)
	lessonRepository := data_base.CreateSqlLessonRepository(
		db,
		os.Getenv("LESSON_TABLE_NAME"),
		os.Getenv("CONTRACT_TABLE_NAME"),
		os.Getenv("TRANSACTION_TABLE_NAME"),
		os.Getenv("SEQUENCE_NAME"),
	)

	return service_logic.CreateServiceModule(
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
	)
}
