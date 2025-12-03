package main

import (
	"data_base_project/data_base"
	"data_base_project/server"
	"data_base_project/service_logic"
	"data_base_project/utility_module"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	_ "github.com/marcboeker/go-duckdb"
)

func main() {
	// Пытаемся загрузить .env, но не падаем, если файла нет (в Docker-окружении
	// переменные обычно приходят из docker-compose env / Kubernetes и т.п.).
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not loaded: %v (using existing environment)", err)
	} else {
		defer utility_module.UnsetEnv()
	}

	// Логируем одновременно в stdout (для docker logs / CI) и в файл внутри контейнера.
	log.Println("backend starting, configuring logger...")
	if logger, err := os.OpenFile("./backend.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err != nil {
		log.Printf("Warning: cannot open backend.log file: %v (logging only to stdout)", err)
	} else {
		mw := io.MultiWriter(os.Stdout, logger)
		log.SetOutput(mw)
		defer logger.Close()
		log.Println("Log file opened, logging to stdout and backend.log")
	}

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
		service_logic.CreateEmailSender(),
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
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
