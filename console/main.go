package main

import (
	"data_base_project/data_base"
	"data_base_project/server"
	"data_base_project/service_logic"
	"data_base_project/utility_module"
	"database/sql"
	"fmt"
	"log"
	"os"

	"console/console_module"

	"github.com/joho/godotenv"
	_ "github.com/marcboeker/go-duckdb"
)

func main() {
	logger, err := os.OpenFile("./backend.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logger.Close()
	log.Printf("Log file opened")
	log.SetOutput(logger)
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	defer utility_module.UnsetEnv()
	log.Printf("Connection string: %s", data_base.GetSqlConnectionString())
	//db, err := data_base.CreateSqlConnection(data_base.GetSqlConnectionString())
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()
	/*
		err = data_base.DropTables(
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
		)
		if err != nil {
			log.Fatalf("Error dropping tables: %v", err)
			return
		}
	*/
	err = data_base.CreateSqlTables(
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
	)
	if err != nil {
		log.Fatalf("Error creating tables: %v", err)
		return
	}
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
	serviceModule := service_logic.CreateServiceModule(
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
	data_base.SetupRoles(
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
	)
	fmt.Println("Server starting on port before setup", os.Getenv("BACKEND_PORT"))
	server := server.SetupServer(serviceModule, os.Getenv("BACKEND_PORT"), log.Default())
	fmt.Println("Server starting on port ", os.Getenv("BACKEND_PORT"))
	go func() {
		console_module.MainMenu(data_base.CreateDataBaseModule(
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
		), db)
	}()
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}
