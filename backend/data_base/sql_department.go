package data_base

import (
	"data_base_project/types"
	"database/sql"
	"fmt"
	"log"

	_ "go.mongodb.org/mongo-driver/mongo"
)

type IDepartmentRepository interface {
	InsertDepartment(department types.DBDepartment) (int64, error)
	DeleteDepartment(departmentID int64) error
	GetDepartmentIdByName(name string) (int64, error)
	GetDepartmentsByHeadID(headID int64) ([]types.DBDepartment, error)
	GetDepartment(id int64) (*types.DBDepartment, error)
	ChangeDepartmentHead(id int64, headID int64) error
	HireInfoInsert(hireInfo types.DBHireInfo) error
	HireInfoDelete(userId int64, departmentId int64) error
	GetUserDepartmentsIDs(userId int64) ([]int64, error)
	GetDepartmentUsersIDs(departmentId int64) ([]int64, error)
	UpdateDepartmentName(departmentID int64, name string) error
}

func CreateSqlDepartmentTable(db *sql.DB, departmentTable string, hireInfoTable string, userTableName string) error {
	query := `
	CREATE TABLE IF NOT EXISTS ` + departmentTable + ` (
		id INTEGER PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		head_id INTEGER NOT NULL
	)
	`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	query = `
	CREATE TABLE IF NOT EXISTS ` + hireInfoTable + ` (
		id INTEGER PRIMARY KEY,
		department_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL
	)
	`
	_, err = db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

type SqlDepartmentRepository struct {
	db              *sql.DB
	departmentTable string
	hireInfoTable   string
	sequenceName    string
}

func CreateSqlDepartmentRepository(db *sql.DB, departmentTable string, hireInfoTable string, sequenceName string) *SqlDepartmentRepository {
	return &SqlDepartmentRepository{
		db:              db,
		departmentTable: departmentTable,
		hireInfoTable:   hireInfoTable,
		sequenceName:    sequenceName,
	}
}

func (r *SqlDepartmentRepository) InsertDepartment(department types.DBDepartment) (int64, error) {
	var id int64
	err := r.db.QueryRow("SELECT nextval('" + r.sequenceName + "')").Scan(&id)
	if err != nil {
		return 0, err
	}
	query := `
	INSERT INTO ` + r.departmentTable + ` (id, name, head_id) VALUES ($1, $2, $3)`
	_, err = r.db.Exec(query, id, department.Name, department.HeadID)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *SqlDepartmentRepository) GetDepartment(id int64) (*types.DBDepartment, error) {
	query := `
	SELECT id, name, head_id FROM ` + r.departmentTable + ` WHERE id = $1
	`
	var department types.DBDepartment
	err := r.db.QueryRow(query, id).Scan(&department.ID, &department.Name, &department.HeadID)
	if err != nil {
		return nil, err
	}
	return &department, nil
}

func (r *SqlDepartmentRepository) GetDepartmentsByHeadID(headID int64) ([]types.DBDepartment, error) {
	query := `
	SELECT id, name, head_id FROM ` + r.departmentTable + ` WHERE head_id = $1
	`
	var departments []types.DBDepartment
	rows, err := r.db.Query(query, headID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var department types.DBDepartment
		err := rows.Scan(&department.ID, &department.Name, &department.HeadID)
		if err != nil {
			return nil, err
		}
		departments = append(departments, department)
	}
	return departments, nil
}

func (r *SqlDepartmentRepository) ChangeDepartmentHead(id int64, headID int64) error {
	query := `
	UPDATE ` + r.departmentTable + ` SET head_id = $1 WHERE id = $2
	`
	result, err := r.db.Exec(query, headID, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("department not found")
	}
	return nil
}

func (r *SqlDepartmentRepository) HireInfoInsert(hireInfo types.DBHireInfo) error {
	var id int64
	err := r.db.QueryRow("SELECT nextval('" + r.sequenceName + "')").Scan(&id)
	if err != nil {
		return err
	}
	query := `
	INSERT INTO ` + r.hireInfoTable + ` (id, department_id, user_id) VALUES ($1, $2, $3)
	`
	_, err = r.db.Exec(query, id, hireInfo.DepartmentID, hireInfo.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (r *SqlDepartmentRepository) HireInfoDelete(userId int64, departmentId int64) error {
	query := `
	DELETE FROM ` + r.hireInfoTable + ` WHERE user_id = $1 AND department_id = $2
	`
	result, err := r.db.Exec(query, userId, departmentId)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("hire info not found")
	}
	return nil
}

func (r *SqlDepartmentRepository) GetUserDepartmentsIDs(userId int64) ([]int64, error) {
	query := `
	SELECT department_id FROM ` + r.hireInfoTable + ` WHERE user_id = $1
	`
	var departments []int64
	rows, err := r.db.Query(query, userId)
	if err != nil {
		log.Printf("Error getting user departments IDs: %v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var departmentID int64
		err := rows.Scan(&departmentID)
		if err != nil {
			log.Printf("Error getting user departments IDs: %v", err)
			return nil, err
		}
		departments = append(departments, departmentID)
	}
	return departments, nil
}

func (r *SqlDepartmentRepository) GetDepartmentUsersIDs(departmentId int64) ([]int64, error) {
	query := `
	SELECT user_id FROM ` + r.hireInfoTable + ` WHERE department_id = $1
	`
	var users []int64
	rows, err := r.db.Query(query, departmentId)
	if err != nil {
		log.Printf("Error getting department users IDs: %v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var userID int64
		err := rows.Scan(&userID)
		if err != nil {
			return nil, err
		}
		users = append(users, userID)
	}
	return users, nil
}

func (r *SqlDepartmentRepository) GetDepartmentIdByName(name string) (int64, error) {
	query := `
	SELECT id FROM ` + r.departmentTable + ` WHERE name = $1
	`
	var id int64
	err := r.db.QueryRow(query, name).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *SqlDepartmentRepository) UpdateDepartmentName(departmentID int64, name string) error {
	query := `
	UPDATE ` + r.departmentTable + ` SET name = $1 WHERE id = $2
	`
	_, err := r.db.Exec(query, name, departmentID)
	if err != nil {
		return err
	}
	return nil
}

func (r *SqlDepartmentRepository) DeleteDepartment(departmentID int64) error {
	query := `
	DELETE FROM ` + r.departmentTable + ` WHERE id = $1
	`
	_, err := r.db.Exec(query, departmentID)
	if err != nil {
		return err
	}
	return nil
}
