package data_base

import (
	tu "data_base_project/test_database_utility"
	"data_base_project/types"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/marcboeker/go-duckdb"
)

func setupReviewTables(db *sql.DB) error {
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
		return fmt.Errorf("error creating client table: %v", err)
	}
	err = CreateSqlRepetitorTable(db, "repetitors", "users", "sequence")
	if err != nil {
		return fmt.Errorf("error creating repetitor table: %v", err)
	}
	err = CreateSqlReviewTable(db, "reviews", "users", "sequence")
	if err != nil {
		return fmt.Errorf("error creating review table: %v", err)
	}
	return nil
}

func TestInsertReviewCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	err = setupReviewTables(db)
	if err != nil {
		t.Fatalf("Error setting up review tables: %v", err)
	}
	reviewRepository := CreateSqlReviewRepository(db, "reviews", "sequence")
	clientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	repetitorRepository := CreateSqlRepetitorRepository(db, "personal_data", "users", "repetitors", "auth", "resume", "review", "sequence")
	repetitorID, err := repetitorRepository.InsertRepetitor(tu.TestRepetitor, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting repetitor: %v", err)
	}
	tu.TestReview.RepetitorID = repetitorID
	clientID, err := clientRepository.InsertClient(tu.TestClient, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting client: %v", err)
	}
	tu.TestReview.ClientID = clientID
	_, err = reviewRepository.InsertReview(tu.TestReview)
	if err != nil {
		t.Fatalf("Error inserting review: %v", err)
	}
}

func TestInsertReviewInSeqCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	err = setupReviewTables(db)
	if err != nil {
		t.Fatalf("Error setting up review tables: %v", err)
	}
	reviewRepository := CreateSqlReviewRepository(db, "reviews", "sequence")
	clientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	repetitorRepository := CreateSqlRepetitorRepository(db, "personal_data", "users", "repetitors", "auth", "resume", "review", "sequence")
	repetitorID, err := repetitorRepository.InsertRepetitor(tu.TestRepetitor, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting repetitor: %v", err)
	}
	tu.TestReview.RepetitorID = repetitorID
	clientID, err := clientRepository.InsertClient(tu.TestClient, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting client: %v", err)
	}
	tu.TestReview.ClientID = clientID
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Error beginning transaction: %v", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()
	_, err = reviewRepository.InsertReviewInSeq(tx, tu.TestReview)
	if err != nil {
		t.Fatalf("Error inserting review: %v", err)
	}
}
func TestGetReviewCorrect(t *testing.T) {
	db := SetupDatabase(t)
	defer func() {
		err := db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	err := setupReviewTables(db)
	if err != nil {
		t.Fatalf("Error setting up review tables: %v", err)
	}
	reviewRepository := CreateSqlReviewRepository(db, "reviews", "sequence")
	clientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	repetitorRepository := CreateSqlRepetitorRepository(db, "personal_data", "users", "repetitors", "auth", "resume", "review", "sequence")
	userID, err := clientRepository.InsertClient(tu.TestClient, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting client: %v", err)
	}
	tu.TestReview.ClientID = userID
	userID, err = repetitorRepository.InsertRepetitor(tu.TestRepetitor, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting repetitor: %v", err)
	}
	tu.TestReview.RepetitorID = userID
	reviewID, err := reviewRepository.InsertReview(tu.TestReview)
	if err != nil {
		t.Fatalf("Error inserting review: %v", err)
	}
	review, err := reviewRepository.GetReview(reviewID)
	if err != nil {
		t.Fatalf("Error getting review: %v", err)
	}
	CheckReview(t, review, reviewID, tu.TestReview.Rating, tu.TestReview.Comment)
}

func TestGetReviewIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	err = setupReviewTables(db)
	if err != nil {
		t.Fatalf("Error setting up review tables: %v", err)
	}
	reviewRepository := CreateSqlReviewRepository(db, "reviews", "sequence")
	_, err = reviewRepository.GetReview(1)
	if err == nil {
		t.Fatalf("No error getting review: %v", err)
	}
}

func CheckReview(
	t *testing.T,
	Review *types.DBReview,
	ReviewID int64,
	Rating int,
	Comment string,
) {
	if Review.ID != ReviewID {
		t.Fatalf("Review not found: %v", Review)
	}
	if Review.Rating != Rating {
		t.Fatalf("Review not found: %v", Review)
	}
	if Review.Comment != Comment {
		t.Fatalf("Review not found: %v", Review)
	}
}

func TestUpdateReviewCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	err = setupReviewTables(db)
	if err != nil {
		t.Fatalf("Error setting up review tables: %v", err)
	}
	reviewRepository := CreateSqlReviewRepository(db, "reviews", "sequence")
	clientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	repetitorRepository := CreateSqlRepetitorRepository(db, "personal_data", "users", "repetitors", "auth", "resume", "review", "sequence")
	userID, err := clientRepository.InsertClient(tu.TestClient, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting client: %v", err)
	}
	tu.TestReview.ClientID = userID
	repetitorID, err := repetitorRepository.InsertRepetitor(tu.TestRepetitor, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting repetitor: %v", err)
	}
	tu.TestReview.RepetitorID = repetitorID
	reviewID, err := reviewRepository.InsertReview(tu.TestReview)
	if err != nil {
		t.Fatalf("Error inserting review: %v", err)
	}
	newReview := types.DBReview{
		ID:      reviewID,
		Rating:  1,
		Comment: "new comment",
	}
	err = reviewRepository.UpdateReview(newReview)
	if err != nil {
		t.Fatalf("Error updating review: %v", err)
	}
	review, err := reviewRepository.GetReview(reviewID)
	if err != nil {
		t.Fatalf("Error getting review: %v", err)
	}
	CheckReview(t, review, reviewID, newReview.Rating, newReview.Comment)
}

func TestGetReviewsByRepetitorIDCorrect(t *testing.T) {
	db := SetupTestDb(t)
	defer func() {
		err := db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	err := setupReviewTables(db)
	if err != nil {
		t.Fatalf("Error setting up review tables: %v", err)
	}
	reviewRepository := CreateSqlReviewRepository(db, "reviews", "sequence")
	clientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	repetitorRepository := CreateSqlRepetitorRepository(db, "personal_data", "users", "repetitors", "auth", "resume", "review", "sequence")
	userID, err := clientRepository.InsertClient(tu.TestClient, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting client: %v", err)
	}
	tu.TestReview.ClientID = userID
	repetitorID, err := repetitorRepository.InsertRepetitor(tu.TestRepetitor, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting repetitor: %v", err)
	}
	tu.TestReview.RepetitorID = repetitorID
	_, err = reviewRepository.InsertReview(tu.TestReview)
	if err != nil {
		t.Fatalf("Error inserting review: %v", err)
	}

	_, err = reviewRepository.InsertReview(tu.TestReview)
	if err != nil {
		t.Fatalf("Error inserting review: %v", err)
	}
	userID, err = repetitorRepository.InsertRepetitor(tu.TestRepetitor, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting repetitor: %v", err)
	}
	tu.TestReview.RepetitorID = userID
	_, err = reviewRepository.InsertReview(tu.TestReview)
	if err != nil {
		t.Fatalf("Error inserting review: %v", err)
	}
	_, err = reviewRepository.GetReviewsByRepetitorID(repetitorID, 0, 10)
	if err != nil {
		t.Fatalf("Error getting reviews: %v", err)
	}

}

func TestGetReviewsByClientIDCorrect(t *testing.T) {
	db := SetupTestDb(t)
	defer func() {
		err := db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	err := setupReviewTables(db)
	if err != nil {
		t.Fatalf("Error setting up review tables: %v", err)
	}
	reviewRepository := CreateSqlReviewRepository(db, "reviews", "sequence")
	clientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	repetitorRepository := CreateSqlRepetitorRepository(db, "personal_data", "users", "repetitors", "auth", "resume", "review", "sequence")
	clientID, err := clientRepository.InsertClient(tu.TestClient, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting client: %v", err)
	}
	tu.TestReview.ClientID = clientID
	repetitorID, err := repetitorRepository.InsertRepetitor(tu.TestRepetitor, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting repetitor: %v", err)
	}
	tu.TestReview.RepetitorID = repetitorID
	_, err = reviewRepository.InsertReview(tu.TestReview)
	if err != nil {
		t.Fatalf("Error inserting review: %v", err)
	}

	_, err = reviewRepository.InsertReview(tu.TestReview)
	if err != nil {
		t.Fatalf("Error inserting review: %v", err)
	}
	userID, err := clientRepository.InsertClient(tu.TestClient, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting client: %v", err)
	}
	tu.TestReview.ClientID = userID
	_, err = reviewRepository.InsertReview(tu.TestReview)
	if err != nil {
		t.Fatalf("Error inserting review: %v", err)
	}
	_, err = reviewRepository.GetReviewsByClientID(clientID, 0, 10)
	if err != nil {
		t.Fatalf("Error getting reviews: %v", err)
	}

}
