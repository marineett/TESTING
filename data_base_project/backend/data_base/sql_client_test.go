package data_base

import (
	tu "data_base_project/test_database_utility"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/marcboeker/go-duckdb"
)

func setupClientTables(db *sql.DB) error {
	err := CreateSqlSequence(db, "sequence")
	if err != nil {
		return fmt.Errorf("error creating sequence: %v", err)
	}
	err = CreateSqlPersonalDataTable(db, "personal_data", "sequence")
	if err != nil {
		return fmt.Errorf("error creating personal data table: %v", err)
	}
	err = CreateSqlUserTable(db, "users", "personal_data", "sequence")
	if err != nil {
		return fmt.Errorf("error creating user table: %v", err)
	}
	err = CreateSqlAuthTable(db, "auth", "users", "sequence")
	if err != nil {
		return fmt.Errorf("error creating auth table: %v", err)
	}
	err = CreateSqlClientTable(db, "clients", "users", "sequence")
	if err != nil {
		return fmt.Errorf("error creating Client table: %v", err)
	}
	return nil
}

func TestCreateSqlClientRepositoryCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	ClientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	if ClientRepository == nil {
		t.Fatalf("Error creating Client repository: %v", err)
	}
}

func TestInsertClientCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupClientTables(db)
	if err != nil {
		t.Fatalf("Error setting up Client tables: %v", err)
	}
	ClientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	ClientRepository.InsertClient(tu.TestClient, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting Client: %v", err)
	}
}
func TestGetClientCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupClientTables(db)
	if err != nil {
		t.Fatalf("Error setting up Client tables: %v", err)
	}
	ClientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	ClientID, err := ClientRepository.InsertClient(tu.TestClient, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting Client: %v", err)
	}
	Client, err := ClientRepository.GetClient(ClientID)
	if err != nil {
		t.Fatalf("Error getting Client: %v, ClientID: %v", err, ClientID)
	}
	if Client.ID != ClientID {
		t.Fatalf("Client not found: %v", Client)
	}
	if Client.SummaryRating != tu.TestClient.SummaryRating {
		t.Fatalf("Client not found: %v", Client)
	}
	if Client.ReviewsCount != tu.TestClient.ReviewsCount {
		t.Fatalf("Client not found: %v", Client)
	}
}

func TestGetClientIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupClientTables(db)
	if err != nil {
		t.Fatalf("Error setting up Client tables: %v", err)
	}
	ClientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	_, err = ClientRepository.GetClient(1)
	if err == nil {
		t.Fatalf("No error getting Client: %v", err)
	}
}

func TestUpdateClientPersonalDataCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupClientTables(db)
	if err != nil {
		t.Fatalf("Error setting up Client tables: %v", err)
	}
	if err != nil {
		t.Fatalf("Error creating sequence: %v", err)
	}
	ClientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	ClientID, err := ClientRepository.InsertClient(tu.TestClient, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting Client: %v", err)
	}
	err = ClientRepository.UpdateClientPersonalData(ClientID, tu.TestPD)
	if err != nil {
		t.Fatalf("Error updating Client personal data: %v", err)
	}
	Client, err := ClientRepository.GetClient(ClientID)
	if err != nil {
		t.Fatalf("Error getting Client: %v", err)
	}
	if Client.ID != ClientID {
		t.Fatalf("Client not found: %v", Client)
	}
	if Client.SummaryRating != tu.TestClient.SummaryRating {
		t.Fatalf("Client not found: %v", Client)
	}
	if Client.ReviewsCount != tu.TestClient.ReviewsCount {
		t.Fatalf("Client not found: %v", Client)
	}
}

func TestUpdateClientPersonalDataIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupClientTables(db)
	if err != nil {
		t.Fatalf("Error creating sequence: %v", err)
	}
	ClientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	err = ClientRepository.UpdateClientPersonalData(1, tu.TestPD)
	if err == nil {
		t.Fatalf("No error updating Client personal data: %v", err)
	}
}

func TestUpdateClientPasswordCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupClientTables(db)
	if err != nil {
		t.Fatalf("Error creating sequence: %v", err)
	}
	ClientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	ClientID, err := ClientRepository.InsertClient(tu.TestClient, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting Client: %v", err)
	}
	err = ClientRepository.UpdateClientPassword(ClientID, tu.TestAuthData, "test3")
	if err != nil {
		t.Fatalf("Error updating Client password: %v", err)
	}
}

func TestUpdateClientPasswordIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupClientTables(db)
	if err != nil {
		t.Fatalf("Error creating sequence: %v", err)
	}
	ClientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	err = ClientRepository.UpdateClientPassword(1, tu.TestAuthData, "test3")
	if err == nil {
		t.Fatalf("No error updating Client password: %v", err)
	}
}
