package main

import (
	"context"
	"data_base_project/data_base"
	"data_base_project/utility_module"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"console/console_module"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	defer utility_module.UnsetEnv()
	var db *sql.DB

	db, err = data_base.CreateSqlConnection(data_base.GetSqlConnectionString())
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

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
		os.Getenv("PENDING_CONTRACT_PAYMENT_TRANSACTIONS_TABLE_NAME"),
		os.Getenv("LESSON_TABLE_NAME"),
	)
	if err != nil {
		log.Fatalf("Error dropping tables: %v", err)
		return
	}
	sqlDataBaseModule := console_module.SqlSetup(db)
	client, err := data_base.CreateMongoConnection(data_base.GetMongoConnectionString())
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer client.Disconnect(context.Background())
	fmt.Println("HOST: ", data_base.GetMongoConnectionString())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}
	mongoDataBaseModule := console_module.MongoSetup(client)
	console_module.MainMenu(sqlDataBaseModule, mongoDataBaseModule, db)
}
