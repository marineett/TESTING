package data_base

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/marcboeker/go-duckdb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetSqlConnectionString() string {
	fmt.Println(os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_NAME"), os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASSWORD"))
	ENV_DATABASE_HOST := os.Getenv("DATABASE_HOST")
	ENV_DATABASE_NAME := os.Getenv("DATABASE_NAME")
	ENV_DATABASE_USER := os.Getenv("DATABASE_USER")
	ENV_DATABASE_PASSWORD := os.Getenv("DATABASE_PASSWORD")
	connectionString := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		ENV_DATABASE_USER,
		ENV_DATABASE_PASSWORD,
		ENV_DATABASE_HOST,
		ENV_DATABASE_NAME)

	return connectionString
}

func GetMongoConnectionString() string {
	ENV_DATABASE_HOST := os.Getenv("DATABASE_HOST")
	ENV_DATABASE_PORT := os.Getenv("MONGO_PORT")
	ENV_DATABASE_NAME := os.Getenv("DATABASE_NAME")
	ENV_DATABASE_USER := os.Getenv("DATABASE_USER")
	ENV_DATABASE_PASSWORD := os.Getenv("DATABASE_PASSWORD")
	connectionString := fmt.Sprintf("mongodb://%s:%s@%s/%s?authSource=admin",
		ENV_DATABASE_USER,
		ENV_DATABASE_PASSWORD,
		ENV_DATABASE_HOST+":"+ENV_DATABASE_PORT,
		ENV_DATABASE_NAME)

	return connectionString
}

func SetupRoles(db *sql.DB,
	personalDataTableName string,
	userTableName string,
	authTableName string,
	chatTableName string,
	messageTableName string,
	departmentTableName string,
	hireInfoTableName string,
	clientTableName string,
	resumeTableName string,
	reviewTableName string,
	repetitorTableName string,
	contractTableName string,
	adminTableName string,
	moderatorTableName string,
	transactionTableName string,
	pendingContractPaymentTransactionsTableName string,
	lessonTableName string,
	sequenceName string,
	rolesNeeded bool,
) error {
	if !rolesNeeded {
		return nil
	}
	createRolesQuery := `
	DO $$ 
	BEGIN
		IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'guest') THEN
			CREATE ROLE guest;
		END IF;
		
		IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'client') THEN
			CREATE ROLE client;
		END IF;
		
		IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'repetitor') THEN
			CREATE ROLE repetitor;
		END IF;
		
		IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'moderator') THEN
			CREATE ROLE moderator;
		END IF;
		
		IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'admin') THEN
			CREATE ROLE admin;
		END IF;
	END
	$$;
	`
	_, err := db.Exec(createRolesQuery)
	if err != nil {
		return fmt.Errorf("error creating roles: %v", err)
	}

	privilegesQuery := `

	GRANT SELECT ON ` + personalDataTableName + ` TO moderator, admin;
	GRANT INSERT ON ` + personalDataTableName + ` TO moderator, admin;
	GRANT UPDATE ON ` + personalDataTableName + ` TO moderator, admin;
	GRANT DELETE ON ` + personalDataTableName + ` TO admin;

	GRANT SELECT ON ` + userTableName + ` TO moderator, admin;
	GRANT INSERT ON ` + userTableName + ` TO moderator,admin;
	GRANT UPDATE ON ` + userTableName + ` TO moderator, admin;
	GRANT DELETE ON ` + userTableName + ` TO admin;

	GRANT SELECT ON ` + authTableName + ` TO admin;
	GRANT INSERT ON ` + authTableName + ` TO admin;
	GRANT UPDATE ON ` + authTableName + ` TO admin;
	GRANT DELETE ON ` + authTableName + ` TO admin;

	GRANT SELECT ON ` + chatTableName + ` TO client, repetitor, moderator, admin;
	GRANT INSERT ON ` + chatTableName + ` TO moderator, admin;
	GRANT UPDATE ON ` + chatTableName + ` TO client, repetitor, moderator, admin;
	GRANT DELETE ON ` + chatTableName + ` TO moderator, admin;

	GRANT SELECT ON ` + messageTableName + ` TO client, repetitor, moderator, admin;
	GRANT INSERT ON ` + messageTableName + ` TO client, repetitor, moderator, admin;
	GRANT UPDATE ON ` + messageTableName + ` TO client, repetitor, moderator, admin;
	GRANT DELETE ON ` + messageTableName + ` TO moderator, admin;

	GRANT SELECT ON ` + departmentTableName + ` TO moderator, admin;
	GRANT INSERT ON ` + departmentTableName + ` TO admin;
	GRANT UPDATE ON ` + departmentTableName + ` TO admin;
	GRANT DELETE ON ` + departmentTableName + ` TO admin;

	GRANT SELECT ON ` + clientTableName + ` TO client, repetitor, moderator, admin;
	GRANT INSERT ON ` + clientTableName + ` TO client, moderator, admin;
	GRANT UPDATE ON ` + clientTableName + ` TO client, moderator, admin;
	GRANT DELETE ON ` + clientTableName + ` TO admin;

	GRANT SELECT ON ` + resumeTableName + ` TO client, repetitor, moderator, admin;
	GRANT INSERT ON ` + resumeTableName + ` TO repetitor, moderator, admin;
	GRANT UPDATE ON ` + resumeTableName + ` TO repetitor, moderator, admin;
	GRANT DELETE ON ` + resumeTableName + ` TO admin;

	GRANT SELECT ON ` + reviewTableName + ` TO client, repetitor, moderator, admin;
	GRANT INSERT ON ` + reviewTableName + ` TO client, repetitor, moderator, admin;
	GRANT UPDATE ON ` + reviewTableName + ` TO client, repetitor, moderator, admin;
	GRANT DELETE ON ` + reviewTableName + ` TO moderator, admin;

	GRANT SELECT ON ` + repetitorTableName + ` TO guest, client, repetitor, moderator, admin;
	GRANT INSERT ON ` + repetitorTableName + ` TO repetitor, moderator, admin;
	GRANT UPDATE ON ` + repetitorTableName + ` TO repetitor, moderator, admin;
	GRANT DELETE ON ` + repetitorTableName + ` TO admin;

	GRANT SELECT ON ` + contractTableName + ` TO client, repetitor, moderator, admin;
	GRANT INSERT ON ` + contractTableName + ` TO client, repetitor, moderator, admin;
	GRANT UPDATE ON ` + contractTableName + ` TO moderator, admin;
	GRANT DELETE ON ` + contractTableName + ` TO admin;

	GRANT SELECT ON ` + adminTableName + ` TO admin;
	GRANT INSERT ON ` + adminTableName + ` TO admin;
	GRANT UPDATE ON ` + adminTableName + ` TO admin;
	GRANT DELETE ON ` + adminTableName + ` TO admin;

	GRANT SELECT ON ` + moderatorTableName + ` TO moderator, admin;
	GRANT INSERT ON ` + moderatorTableName + ` TO admin;
	GRANT UPDATE ON ` + moderatorTableName + ` TO admin;
	GRANT DELETE ON ` + moderatorTableName + ` TO admin;



	GRANT SELECT ON ` + userTableName + ` TO client, repetitor, moderator, admin;
	GRANT INSERT ON ` + userTableName + ` TO admin;
	GRANT UPDATE ON ` + userTableName + ` TO moderator, admin;
	GRANT DELETE ON ` + userTableName + ` TO admin;

	GRANT SELECT ON ` + transactionTableName + ` TO client, repetitor, moderator, admin;
	GRANT INSERT ON ` + transactionTableName + ` TO client, repetitor, moderator, admin;
	GRANT UPDATE ON ` + transactionTableName + ` TO moderator, admin;
	GRANT DELETE ON ` + transactionTableName + ` TO admin;

	GRANT SELECT ON ` + pendingContractPaymentTransactionsTableName + ` TO moderator, admin;
	GRANT INSERT ON ` + pendingContractPaymentTransactionsTableName + ` TO moderator, admin;
	GRANT UPDATE ON ` + pendingContractPaymentTransactionsTableName + ` TO moderator, admin;
	GRANT DELETE ON ` + pendingContractPaymentTransactionsTableName + ` TO moderator, admin;

	GRANT SELECT ON ` + lessonTableName + ` TO client, repetitor, moderator, admin;
	GRANT INSERT ON ` + lessonTableName + ` TO client, repetitor, moderator, admin;
	GRANT UPDATE ON ` + lessonTableName + ` TO client, repetitor, moderator, admin;
	GRANT DELETE ON ` + lessonTableName + ` TO admin;
	`
	_, err = db.Exec(privilegesQuery)
	if err != nil {
		return fmt.Errorf("error granting privileges: %v", err)
	}
	return nil
}

func CreateSqlConnection(connectionString string) (*sql.DB, error) {
	log.Println("Attempting to connect to database...")

	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	err = db.Ping()
	if err != nil {
		log.Printf("Error pinging database: %v", err)
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	log.Println("Successfully connected to database!")
	return db, nil
}

func CreateMongoConnection(connectionString string) (*mongo.Client, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}
	return client, nil
}

func createBaseTables(db *sql.DB, personalDataTable, userTableName, authTableName, sequenceName string) error {
	if err := CreateSqlPersonalDataTable(db, personalDataTable, sequenceName); err != nil {
		return err
	}
	if err := CreateSqlUserTable(db, userTableName, personalDataTable, sequenceName); err != nil {
		return err
	}
	if err := CreateSqlAuthTable(db, authTableName, userTableName, sequenceName); err != nil {
		return err
	}
	return nil
}

func createChatTables(db *sql.DB, chatTableName, messageTableName, userTableName string) error {
	if err := CreateSqlChatTable(db, chatTableName, userTableName); err != nil {
		return err
	}
	if err := CreateSqlMessageTable(db, messageTableName, chatTableName, userTableName); err != nil {
		return err
	}
	return nil
}

func createDepartmentTables(db *sql.DB, departmentTableName, hireInfoTableName, userTableName string) error {
	if err := CreateSqlDepartmentTable(db, departmentTableName, hireInfoTableName, userTableName); err != nil {
		return err
	}
	return nil
}

func createUserRoleTables(db *sql.DB, clientTableName, resumeTableName, reviewTableName, repetitorTableName, userTableName, sequenceName string) error {
	if err := CreateSqlClientTable(db, clientTableName, userTableName, sequenceName); err != nil {
		return err
	}
	if err := CreateSqlResumeTable(db, resumeTableName, userTableName, sequenceName); err != nil {
		return err
	}
	if err := CreateSqlReviewTable(db, reviewTableName, userTableName, sequenceName); err != nil {
		return err
	}
	if err := CreateSqlRepetitorTable(db, repetitorTableName, userTableName, resumeTableName); err != nil {
		return err
	}
	return nil
}

func createContractTable(db *sql.DB, contractTableName, userTableName, reviewTableName, repetitorTableName, clientTableName string) error {
	if err := CreateSqlContractTable(db, contractTableName, userTableName, reviewTableName, repetitorTableName, clientTableName); err != nil {
		return err
	}
	return nil
}

func createAdminTables(db *sql.DB, adminTableName, moderatorTableName, userTableName, sequenceName string) error {
	if err := CreateSqlAdminTable(db, adminTableName, userTableName, sequenceName); err != nil {
		return err
	}
	if err := CreateSqlModeratorTable(db, moderatorTableName, userTableName); err != nil {
		return err
	}
	return nil
}

func createTransactionTables(db *sql.DB, transactionTableName, pendingContractPaymentTransactionsTableName, userTableName, sequenceName string) error {
	if err := CreateSqlTransactionTable(db, transactionTableName, userTableName, pendingContractPaymentTransactionsTableName, sequenceName); err != nil {
		return err
	}
	if err := CreateSqlPendingContractPaymentTransactionsTable(db, pendingContractPaymentTransactionsTableName, userTableName, transactionTableName, sequenceName); err != nil {
		return err
	}
	return nil
}

func createLessonTable(db *sql.DB, lessonTableName, contractTableName, transactionTableName string) error {
	if err := CreateSqlLessonTable(db, lessonTableName, contractTableName, transactionTableName); err != nil {
		return err
	}
	return nil
}

func createSequence(db *sql.DB, sequenceName string) error {
	if err := CreateSqlSequence(db, sequenceName); err != nil {
		return err
	}
	return nil
}

func CreateSqlTables(db *sql.DB,
	personalDataTable string,
	userTableName string,
	authTableName string,
	chatTableName string,
	messageTableName string,
	departmentTableName string,
	hireInfoTableName string,
	clientTableName string,
	resumeTableName string,
	reviewTableName string,
	repetitorTableName string,
	contractTableName string,
	adminTableName string,
	moderatorTableName string,
	transactionTableName string,
	pendingContractPaymentTransactionsTableName string,
	lessonTableName string,
	sequenceName string,
) error {
	if err := createBaseTables(db, personalDataTable, userTableName, authTableName, sequenceName); err != nil {
		return err
	}
	if err := createChatTables(db, chatTableName, messageTableName, userTableName); err != nil {
		return err
	}
	if err := createDepartmentTables(db, departmentTableName, hireInfoTableName, userTableName); err != nil {
		return err
	}
	if err := createUserRoleTables(db, clientTableName, resumeTableName, reviewTableName, repetitorTableName, userTableName, sequenceName); err != nil {
		return err
	}
	if err := createContractTable(db, contractTableName, userTableName, reviewTableName, repetitorTableName, clientTableName); err != nil {
		return err
	}
	if err := createAdminTables(db, adminTableName, moderatorTableName, userTableName, sequenceName); err != nil {
		return err
	}
	if err := createTransactionTables(db, transactionTableName, pendingContractPaymentTransactionsTableName, userTableName, sequenceName); err != nil {
		return err
	}
	if err := createLessonTable(db, lessonTableName, contractTableName, transactionTableName); err != nil {
		return err
	}
	if err := createSequence(db, sequenceName); err != nil {
		return err
	}
	return nil
}

func DropTables(db *sql.DB,
	personalDataTable string,
	userTableName string,
	authTableName string,
	chatTableName string,
	messageTableName string,
	departmentTableName string,
	hireInfoTableName string,
	clientTableName string,
	resumeTableName string,
	reviewTableName string,
	repetitorTableName string,
	contractTableName string,
	adminTableName string,
	moderatorTableName string,
	transactionTableName string,
	pendingContractPaymentTransactionsTableName string,
	lessonTableName string) error {
	query := `
	DROP TABLE IF EXISTS ` +
		personalDataTable + `, ` +
		userTableName + `, ` +
		authTableName + `, ` +
		chatTableName + `, ` +
		messageTableName + `, ` +
		departmentTableName + `, ` +
		hireInfoTableName + `, ` +
		clientTableName + `, ` +
		resumeTableName + `, ` +
		reviewTableName + `, ` +
		repetitorTableName + `, ` +
		contractTableName + `, ` +
		adminTableName + `, ` +
		moderatorTableName + `, ` +
		transactionTableName + `, ` +
		pendingContractPaymentTransactionsTableName + `, ` +
		lessonTableName + ` CASCADE;
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error dropping tables: %v", err)
	}
	return nil
}
