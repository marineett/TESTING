package data_base

import (
	"data_base_project/types"
	"database/sql"
	"fmt"
	"time"
)

type IAdminRepository interface {
	InsertAdmin(admin types.DBAdminData, personalData types.DBPersonalData, auth types.DBAuthData) (int64, error)
	GetAdmin(id int64) (*types.DBAdminData, error)
	UpdateAdminPersonalData(adminId int64, personalData types.DBPersonalData) error
	UpdateAdminPassword(adminId int64, authData types.DBAuthData, newPassword string) error
	UpdateAdminDepartment(adminId int64, departmentId int64) error
	UpdateAdminSalary(adminId int64, salary int64) error
}

func CreateSqlAdminTable(db *sql.DB, adminTableName string, userTableName string, sequenceName string) error {
	query := `
	CREATE TABLE IF NOT EXISTS ` + adminTableName + ` (
		id INTEGER PRIMARY KEY,
		department_id INTEGER NOT NULL,
		salary INTEGER NOT NULL,
		FOREIGN KEY (id) REFERENCES ` + userTableName + `(id)
	)`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table %s: %v", adminTableName, err)
	}
	return nil
}

type SqlAdminRepository struct {
	db                     *sql.DB
	adminTable             string
	userRepository         IUserRepository
	personalDataRepository IPersonalDataRepository
	authRepository         IAuthRepository
}

func CreateSqlAdminRepository(db *sql.DB, personalDataTable string, userTable string, adminTable string, authTable string, sequenceName string) *SqlAdminRepository {
	return &SqlAdminRepository{
		db:                     db,
		adminTable:             adminTable,
		userRepository:         CreateSqlUserRepository(db, userTable, sequenceName),
		personalDataRepository: CreateSqlPersonalDataRepository(db, personalDataTable, sequenceName),
		authRepository:         CreateSqlAuthRepository(db, authTable, sequenceName),
	}
}

func (r *SqlAdminRepository) InsertAdmin(admin types.DBAdminData, personalData types.DBPersonalData, auth types.DBAuthData) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	personalData.ID, err = r.personalDataRepository.InsertPersonalDataInSeq(tx, personalData)
	if err != nil {
		return 0, err
	}
	userID, err := r.userRepository.InsertUserInSeq(tx, types.DBUserData{
		PersonalDataID:   personalData.ID,
		RegistrationDate: time.Now(),
		LastLoginDate:    time.Now(),
	})
	if err != nil {
		return 0, err
	}
	_, err = r.authRepository.InsertAuthInSeq(tx, types.DBAuthInfo{
		UserID:   userID,
		UserType: types.Admin,
		Login:    auth.Login,
		Password: auth.Password,
	})
	if err != nil {
		return 0, err
	}
	query := `
	INSERT INTO ` + r.adminTable + ` (id, department_id, salary)
	VALUES ($1, $2, $3)
	`
	_, err = tx.Exec(query, userID, 0, admin.Salary)
	if err != nil {
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (r *SqlAdminRepository) GetAdmin(id int64) (*types.DBAdminData, error) {
	query := `
	SELECT id, department_id, salary FROM ` + r.adminTable + ` WHERE id = $1
	`
	var admin types.DBAdminData
	err := r.db.QueryRow(query, id).Scan(&admin.ID, &admin.DepartmentID, &admin.Salary)
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *SqlAdminRepository) UpdateAdminPersonalData(adminId int64, personalData types.DBPersonalData) error {
	adminData, err := r.userRepository.GetUser(adminId)
	if err != nil {
		return err
	}
	return r.personalDataRepository.UpdatePersonalData(adminData.PersonalDataID, personalData)
}

func (r *SqlAdminRepository) UpdateAdminPassword(adminId int64, authData types.DBAuthData, newPassword string) error {
	return r.authRepository.ChangePassword(adminId, authData, newPassword)
}

func (r *SqlAdminRepository) UpdateAdminDepartment(adminId int64, departmentId int64) error {
	query := `
	UPDATE ` + r.adminTable + ` SET department_id = $1 WHERE id = $2
	`
	result, err := r.db.Exec(query, departmentId, adminId)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("admin not found")
	}
	return nil
}

func (r *SqlAdminRepository) UpdateAdminSalary(admin_id int64, salary int64) error {
	query := `
	UPDATE ` + r.adminTable + ` SET salary = $1 WHERE id = $2
	`
	result, err := r.db.Exec(query, salary, admin_id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("admin not found")
	}
	return nil
}
