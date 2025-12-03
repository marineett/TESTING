package data_base

import (
	tu "data_base_project/test_database_utility"
	"data_base_project/types"
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/marcboeker/go-duckdb"
)

func setupModeratorTables(db *sql.DB) error {
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
	err = CreateSqlModeratorTable(db, "moderators", "users")
	if err != nil {
		return fmt.Errorf("error creating Moderator table: %v", err)
	}
	return nil
}

func TestCreateSqlModeratorRepositoryCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	ModeratorRepository := CreateSqlModeratorRepository(db, "personal_data", "users", "moderators", "auth", "sequence")
	if ModeratorRepository == nil {
		t.Fatalf("Error creating Moderator repository: %v", err)
	}
}

func TestInsertModeratorCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupModeratorTables(db)
	if err != nil {
		t.Fatalf("Error setting up Moderator tables: %v", err)
	}
	ModeratorRepository := CreateSqlModeratorRepository(db, "personal_data", "users", "moderators", "auth", "sequence")
	ModeratorRepository.InsertModerator(tu.TestModeratorData, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting Moderator: %v", err)
	}
}
func TestGetModeratorCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupModeratorTables(db)
	if err != nil {
		t.Fatalf("Error setting up Moderator tables: %v", err)
	}
	ModeratorRepository := CreateSqlModeratorRepository(db, "personal_data", "users", "moderators", "auth", "sequence")
	ModeratorID, err := ModeratorRepository.InsertModerator(tu.TestModeratorData, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting Moderator: %v", err)
	}
	Moderator, err := ModeratorRepository.GetModerator(ModeratorID)
	if err != nil {
		t.Fatalf("Error getting Moderator: %v, ModeratorID: %v", err, ModeratorID)
	}
	if Moderator.ID != ModeratorID {
		t.Fatalf("Moderator not found: %v", Moderator)
	}
	if Moderator.Salary != tu.TestSalary {
		t.Fatalf("Moderator not found: %v", Moderator)
	}
}

func TestGetModeratorIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupModeratorTables(db)
	if err != nil {
		t.Fatalf("Error setting up Moderator tables: %v", err)
	}
	ModeratorRepository := CreateSqlModeratorRepository(db, "personal_data", "users", "moderators", "auth", "sequence")
	_, err = ModeratorRepository.GetModerator(1)
	if err == nil {
		t.Fatalf("No error getting Moderator: %v", err)
	}
}

func TestUpdateModeratorPersonalDataCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupModeratorTables(db)
	if err != nil {
		t.Fatalf("Error setting up Moderator tables: %v", err)
	}
	if err != nil {
		t.Fatalf("Error creating sequence: %v", err)
	}
	ModeratorRepository := CreateSqlModeratorRepository(db, "personal_data", "users", "moderators", "auth", "sequence")
	moderatorID, err := ModeratorRepository.InsertModerator(tu.TestModeratorData, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting Moderator: %v", err)
	}
	newPersonalData := types.DBPersonalData{
		TelephoneNumber: "88005553536",
		Email:           "test2@test.com",
		FirstName:       "Petr",
		LastName:        "Petrov",
		MiddleName:      "Petrovich",
		DBPassportData: types.DBPassportData{
			PassportNumber:   "1234567890",
			PassportSeries:   "1234",
			PassportDate:     time.Now(),
			PassportIssuedBy: "Moscow",
		},
	}
	err = ModeratorRepository.UpdateModeratorPersonalData(moderatorID, newPersonalData)
	if err != nil {
		t.Fatalf("Error updating Moderator personal data: %v", err)
	}
	moderator, err := ModeratorRepository.GetModerator(moderatorID)
	if err != nil {
		t.Fatalf("Error getting Moderator: %v", err)
	}
	if moderator.ID != moderatorID {
		t.Fatalf("Moderator not found: %v", moderator)
	}
	if moderator.Salary != tu.TestSalary {
		t.Fatalf("Moderator not found: %v", moderator)
	}
}

func TestUpdateModeratorPersonalDataIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupModeratorTables(db)
	if err != nil {
		t.Fatalf("Error creating sequence: %v", err)
	}
	ModeratorRepository := CreateSqlModeratorRepository(db, "personal_data", "users", "moderators", "auth", "sequence")
	err = ModeratorRepository.UpdateModeratorPersonalData(1, tu.TestPD)
	if err == nil {
		t.Fatalf("No error updating Moderator personal data: %v", err)
	}
}

func TestUpdateModeratorPasswordCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupModeratorTables(db)
	if err != nil {
		t.Fatalf("Error creating sequence: %v", err)
	}
	ModeratorRepository := CreateSqlModeratorRepository(db, "personal_data", "users", "moderators", "auth", "sequence")
	ModeratorID, err := ModeratorRepository.InsertModerator(tu.TestModeratorData, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting Moderator: %v", err)
	}
	err = ModeratorRepository.UpdateModeratorPassword(ModeratorID, tu.TestAuthData, "test3")
	if err != nil {
		t.Fatalf("Error updating Moderator password: %v", err)
	}
}

func TestUpdateModeratorPasswordIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupModeratorTables(db)
	if err != nil {
		t.Fatalf("Error creating sequence: %v", err)
	}
	ModeratorRepository := CreateSqlModeratorRepository(db, "personal_data", "users", "moderators", "auth", "sequence")
	err = ModeratorRepository.UpdateModeratorPassword(1, tu.TestAuthData, "test3")
	if err == nil {
		t.Fatalf("No error updating Moderator password: %v", err)
	}
}
func TestUpdateModeratorSalaryCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupModeratorTables(db)
	if err != nil {
		t.Fatalf("Error creating sequence: %v", err)
	}
	err = CreateSqlUserTable(db, "users", "personal_data", "sequence")
	if err != nil {
		t.Fatalf("Error creating user table: %v", err)
	}
	ModeratorRepository := CreateSqlModeratorRepository(db, "personal_data", "users", "moderators", "auth", "sequence")
	ModeratorID, err := ModeratorRepository.InsertModerator(tu.TestModeratorData, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting Moderator: %v", err)
	}
	err = ModeratorRepository.UpdateModeratorSalary(ModeratorID, 100000)
	if err != nil {
		t.Fatalf("Error updating Moderator salary: %v", err)
	}
	Moderator, err := ModeratorRepository.GetModerator(ModeratorID)
	if err != nil {
		t.Fatalf("Error getting Moderator: %v", err)
	}
	if Moderator.ID != ModeratorID {
		t.Fatalf("Moderator not found: %v", Moderator)
	}
	if Moderator.Salary != 100000 {
		t.Fatalf("Moderator not found: %v", Moderator)
	}
}

func TestUpdateModeratorSalaryIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupModeratorTables(db)
	if err != nil {
		t.Fatalf("Error creating sequence: %v", err)
	}
	ModeratorRepository := CreateSqlModeratorRepository(db, "personal_data", "users", "moderators", "auth", "sequence")
	err = ModeratorRepository.UpdateModeratorSalary(1, 100000)
	if err == nil {
		t.Fatalf("No error updating Moderator salary: %v", err)
	}
}

func TestGetModeratorsCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupModeratorTables(db)
	if err != nil {
		t.Fatalf("Error setting up Moderator tables: %v", err)
	}
	ModeratorRepository := CreateSqlModeratorRepository(db, "personal_data", "users", "moderators", "auth", "sequence")
	_, err = ModeratorRepository.InsertModerator(tu.TestModeratorData, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting Moderator: %v", err)
	}
	_, err = ModeratorRepository.InsertModerator(tu.TestModeratorData, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting Moderator: %v", err)
	}
	Moderators, err := ModeratorRepository.GetModerators()
	if err != nil {
		t.Fatalf("Error getting Moderators: %v", err)
	}
	if len(Moderators) != 2 {
		t.Fatalf("Moderators not found: %v", Moderators)
	}
}

func TestGetModeratorsIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	ModeratorRepository := CreateSqlModeratorRepository(db, "personal_data", "users", "moderators", "auth", "sequence")
	_, err = ModeratorRepository.GetModerators()
	if err == nil {
		t.Fatalf("No error getting Moderators: %v", err)
	}
}
