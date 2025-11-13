package test_service_utility

import (
	"data_base_project/data_base"
	"data_base_project/types"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"
)

func SetupModule(db *sql.DB) (*data_base.DataBaseModule, error) {
	err := data_base.CreateSqlTables(db,
		"personal_data",
		"users",
		"auth",
		"chat",
		"message",
		"department",
		"hire_info",
		"client",
		"resume",
		"review",
		"repetitor",
		"contract",
		"admin",
		"moderator",
		"transaction",
		"pending_contract_payment_transactions",
		"lesson",
		"sequence",
	)
	if err != nil {
		return nil, fmt.Errorf("error creating tables: %v", err)
	}
	userRepository := data_base.CreateSqlUserRepository(db, "users", "sequence")
	authRepository := data_base.CreateSqlAuthRepository(db, "auth", "sequence")
	adminRepository := data_base.CreateSqlAdminRepository(db, "personal_data", "users", "admin", "auth", "sequence")
	moderatorRepository := data_base.CreateSqlModeratorRepository(db, "personal_data", "users", "moderator", "auth", "sequence")
	clientRepository := data_base.CreateSqlClientRepository(db, "personal_data", "users", "client", "auth", "sequence")
	repetitorRepository := data_base.CreateSqlRepetitorRepository(db, "personal_data", "users", "repetitor", "auth", "resume", "review", "sequence")
	contractRepository := data_base.CreateSqlContractRepository(db, "contract", "sequence")
	reviewRepository := data_base.CreateSqlReviewRepository(db, "review", "sequence")
	chatRepository := data_base.CreateSqlChatRepository(db, "chat", "sequence")
	messageRepository := data_base.CreateSqlMessageRepository(db, "message", "sequence")
	resumeRepository := data_base.CreateSqlResumeRepository(db, "resume", "sequence")
	transactionRepository := data_base.CreateSqlTransactionRepository(db, "transaction", "pending_contract_payment_transactions", "sequence")
	departmentRepository := data_base.CreateSqlDepartmentRepository(db, "department", "hire_info", "sequence")
	personalDataRepository := data_base.CreateSqlPersonalDataRepository(db, "personal_data", "sequence")
	lessonRepository := data_base.CreateSqlLessonRepository(db, "lesson", "contract", "transaction", "sequence")
	return data_base.CreateDataBaseModule(
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
	), nil
}

type TestPersonalDataRepository struct {
	data map[int64]*types.DBPersonalData
}

func (r *TestPersonalDataRepository) InsertPersonalData(personalData types.DBPersonalData) (int64, error) {
	personalData.ID = int64(len(r.data) + 1)
	r.data[personalData.ID] = &personalData
	return personalData.ID, nil
}

func (r *TestPersonalDataRepository) GetPersonalData(id int64) (*types.DBPersonalData, error) {
	value, ok := r.data[id]
	if !ok {
		return nil, errors.New("personal data not found")
	} else {
		return value, nil
	}
}

func (r *TestPersonalDataRepository) UpdatePersonalData(id int64, personalData types.DBPersonalData) error {
	if _, ok := r.data[id]; !ok {
		return errors.New("personal data not found")
	}
	r.data[id] = &personalData
	return nil
}

func (r *TestPersonalDataRepository) InsertPersonalDataInSeq(tx *sql.Tx, personalData types.DBPersonalData) (int64, error) {
	return r.InsertPersonalData(personalData)
}

func CreateTestPersonalDataRepository() *TestPersonalDataRepository {
	return &TestPersonalDataRepository{
		data: make(map[int64]*types.DBPersonalData),
	}
}

type TestUserRepository struct {
	data map[int64]*types.DBUserData
}

func (r *TestUserRepository) InsertUser(user types.DBUserData) (int64, error) {
	user.ID = int64(len(r.data) + 1)
	r.data[user.ID] = &user
	return user.ID, nil
}

func (r *TestUserRepository) GetUser(id int64) (*types.DBUserData, error) {
	value, ok := r.data[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return value, nil
}

func (r *TestUserRepository) InsertUserInSeq(tx *sql.Tx, user types.DBUserData) (int64, error) {
	return r.InsertUser(user)
}

func CreateTestUserRepository() *TestUserRepository {
	return &TestUserRepository{
		data: make(map[int64]*types.DBUserData),
	}
}

type TestAuthRepository struct {
	data map[int64]*types.DBAuthInfo
}

func (r *TestAuthRepository) InsertAuth(auth types.DBAuthInfo) (int64, error) {
	auth.ID = int64(len(r.data) + 1)
	r.data[auth.ID] = &auth
	return auth.ID, nil
}

func (r *TestAuthRepository) InsertAuthInSeq(tx *sql.Tx, auth types.DBAuthInfo) (int64, error) {
	return r.InsertAuth(auth)
}

func CreateTestAuthRepository() *TestAuthRepository {
	return &TestAuthRepository{
		data: make(map[int64]*types.DBAuthInfo),
	}
}

func (r *TestAuthRepository) TestGetAuth(id int64) (*types.DBAuthInfo, error) {
	value, ok := r.data[id]
	if !ok {
		return nil, errors.New("auth not found")
	}
	return value, nil
}

func (r *TestAuthRepository) ChangePassword(userId int64, authData types.DBAuthData, newPassword string) error {
	if _, ok := r.data[userId]; !ok {
		return errors.New("auth not found")
	}
	if r.data[userId].Login != authData.Login || r.data[userId].Password != authData.Password {
		return errors.New("invalid auth data for user with id " + strconv.FormatInt(userId, 10))
	}
	r.data[userId].Password = newPassword
	return nil
}

func (r *TestAuthRepository) Authorize(authData types.DBAuthData) (types.DBAuthVerdict, error) {
	for _, auth := range r.data {
		if auth.Login == authData.Login && auth.Password == authData.Password {
			return types.DBAuthVerdict{
				UserID:   auth.UserID,
				UserType: auth.UserType,
			}, nil
		}
	}
	return types.DBAuthVerdict{}, errors.New("invalid auth data")
}

func (r *TestAuthRepository) CheckLogin(login string) (bool, error) {
	for _, auth := range r.data {
		if auth.Login == login {
			return true, nil
		}
	}
	return false, nil
}

type TestAdminRepository struct {
	data                   map[int64]*types.DBAdminData
	personalDataRepository data_base.IPersonalDataRepository
	authRepository         data_base.IAuthRepository
	userRepository         data_base.IUserRepository
}

func CreateTestAdminRepository(
	personalDataRepository data_base.IPersonalDataRepository,
	authRepository data_base.IAuthRepository,
	userRepository data_base.IUserRepository,
) *TestAdminRepository {
	return &TestAdminRepository{
		data:                   make(map[int64]*types.DBAdminData),
		personalDataRepository: personalDataRepository,
		authRepository:         authRepository,
		userRepository:         userRepository,
	}
}

func (r *TestAdminRepository) InsertAdmin(admin types.DBAdminData, personalData types.DBPersonalData, auth types.DBAuthData) (int64, error) {

	personalDataID, err := r.personalDataRepository.InsertPersonalData(personalData)
	if err != nil {
		return 0, err
	}
	userID, err := r.userRepository.InsertUser(types.DBUserData{
		PersonalDataID:   personalDataID,
		RegistrationDate: time.Now(),
		LastLoginDate:    time.Now(),
	})
	if err != nil {
		return 0, err
	}
	_, err = r.authRepository.InsertAuth(types.DBAuthInfo{
		UserID:   userID,
		UserType: types.Admin,
		Login:    auth.Login,
		Password: auth.Password,
	})
	if err != nil {
		return 0, err
	}
	admin.ID = userID
	r.data[admin.ID] = &admin
	return admin.ID, nil
}

func (r *TestAdminRepository) GetAdmin(id int64) (*types.DBAdminData, error) {
	value, ok := r.data[id]
	if !ok {
		return nil, errors.New("admin not found")
	}
	return value, nil
}

func (r *TestAdminRepository) UpdateAdminPersonalData(adminId int64, personalData types.DBPersonalData) error {
	if _, ok := r.data[adminId]; !ok {
		return errors.New("admin not found")
	}
	r.personalDataRepository.UpdatePersonalData(adminId, personalData)
	return nil
}

func (r *TestAdminRepository) UpdateAdminPassword(adminId int64, authData types.DBAuthData, newPassword string) error {
	if _, ok := r.data[adminId]; !ok {
		return errors.New("admin not found")
	}
	r.authRepository.ChangePassword(adminId, authData, newPassword)
	return nil
}

func (r *TestAdminRepository) UpdateAdminDepartment(adminId int64, departmentId int64) error {
	if _, ok := r.data[adminId]; !ok {
		return errors.New("admin not found")
	}
	r.data[adminId].DepartmentID = departmentId
	return nil
}

func (r *TestAdminRepository) UpdateAdminSalary(adminId int64, salary int64) error {
	if _, ok := r.data[adminId]; !ok {
		return errors.New("admin not found")
	}
	r.data[adminId].Salary = salary
	return nil
}

type TestModeratorRepository struct {
	data                   map[int64]*types.DBModeratorData
	personalDataRepository data_base.IPersonalDataRepository
	authRepository         data_base.IAuthRepository
	userRepository         data_base.IUserRepository
}

func CreateTestModeratorRepository(
	personalDataRepository data_base.IPersonalDataRepository,
	authRepository data_base.IAuthRepository,
	userRepository data_base.IUserRepository,
) *TestModeratorRepository {
	return &TestModeratorRepository{
		data:                   make(map[int64]*types.DBModeratorData),
		personalDataRepository: personalDataRepository,
		authRepository:         authRepository,
		userRepository:         userRepository,
	}
}

func (r *TestModeratorRepository) GetModeratorsByDepartmentID(departmentID int64) ([]int64, error) {
	return nil, nil
}

func (r *TestModeratorRepository) InsertModerator(moderator types.DBModeratorData, personalData types.DBPersonalData, auth types.DBAuthData) (int64, error) {
	personalDataID, err := r.personalDataRepository.InsertPersonalData(personalData)
	if err != nil {
		return 0, err
	}
	userID, err := r.userRepository.InsertUser(types.DBUserData{
		PersonalDataID:   personalDataID,
		RegistrationDate: time.Now(),
		LastLoginDate:    time.Now(),
	})
	if err != nil {
		return 0, err
	}
	_, err = r.authRepository.InsertAuth(types.DBAuthInfo{
		UserID:   userID,
		UserType: types.Moderator,
		Login:    auth.Login,
		Password: auth.Password,
	})
	if err != nil {
		return 0, err
	}
	moderator.ID = userID
	r.data[moderator.ID] = &moderator
	return moderator.ID, nil
}

func (r *TestModeratorRepository) GetModerator(id int64) (*types.DBModeratorData, error) {
	value, ok := r.data[id]
	if !ok {
		return nil, errors.New("moderator not found")
	}
	return value, nil
}

func (r *TestModeratorRepository) UpdateModeratorPersonalData(moderatorId int64, personalData types.DBPersonalData) error {
	if _, ok := r.data[moderatorId]; !ok {
		return errors.New("moderator not found")
	}
	r.personalDataRepository.UpdatePersonalData(moderatorId, personalData)
	return nil
}

func (r *TestModeratorRepository) UpdateModeratorPassword(moderatorId int64, authData types.DBAuthData, newPassword string) error {
	if _, ok := r.data[moderatorId]; !ok {
		return errors.New("moderator not found")
	}
	r.authRepository.ChangePassword(moderatorId, authData, newPassword)
	return nil
}

func (r *TestModeratorRepository) UpdateModeratorSalary(moderatorId int64, salary int64) error {
	if _, ok := r.data[moderatorId]; !ok {
		return errors.New("moderator not found")
	}
	r.data[moderatorId].Salary = salary
	return nil
}

func (r *TestModeratorRepository) GetModerators() ([]int64, error) {
	return nil, nil
}

type TestRepetiorRepository struct {
	data                   map[int64]*types.DBRepetitorData
	personalDataRepository data_base.IPersonalDataRepository
	authRepository         data_base.IAuthRepository
	userRepository         data_base.IUserRepository
	resumeRepository       data_base.IResumeRepository
}

func CreateTestRepetiorRepository(
	personalDataRepository data_base.IPersonalDataRepository,
	authRepository data_base.IAuthRepository,
	userRepository data_base.IUserRepository,
	resumeRepository data_base.IResumeRepository,
) *TestRepetiorRepository {
	return &TestRepetiorRepository{
		data:                   make(map[int64]*types.DBRepetitorData),
		personalDataRepository: personalDataRepository,
		authRepository:         authRepository,
		userRepository:         userRepository,
		resumeRepository:       resumeRepository,
	}
}

func (r *TestRepetiorRepository) InsertRepetitor(repetitor types.DBRepetitorData, personalData types.DBPersonalData, auth types.DBAuthData) (int64, error) {
	personalDataID, err := r.personalDataRepository.InsertPersonalData(personalData)
	if err != nil {
		return 0, err
	}
	userID, err := r.userRepository.InsertUser(types.DBUserData{
		PersonalDataID:   personalDataID,
		RegistrationDate: time.Now(),
		LastLoginDate:    time.Now(),
	})
	if err != nil {
		return 0, err
	}
	_, err = r.authRepository.InsertAuth(types.DBAuthInfo{
		UserID:   userID,
		UserType: types.Repetitor,
		Login:    auth.Login,
		Password: auth.Password,
	})
	if err != nil {
		return 0, err
	}
	_, err = r.resumeRepository.InsertResume(types.DBResume{
		RepetitorID: userID,
	})
	if err != nil {
		return 0, err
	}
	repetitor.ID = userID
	r.data[repetitor.ID] = &repetitor
	return repetitor.ID, nil
}

func (r *TestRepetiorRepository) GetRepetitor(id int64) (*types.DBRepetitorData, error) {
	value, ok := r.data[id]
	if !ok {
		return nil, errors.New("repetitor not found")
	}
	return value, nil
}

func (r *TestRepetiorRepository) UpdateRepetitorPersonalData(repetitorId int64, personalData types.DBPersonalData) error {
	if _, ok := r.data[repetitorId]; !ok {
		return errors.New("repetitor not found")
	}
	r.personalDataRepository.UpdatePersonalData(repetitorId, personalData)
	return nil
}

func (r *TestRepetiorRepository) UpdateRepetitorPassword(repetitorId int64, authData types.DBAuthData, newPassword string) error {
	if _, ok := r.data[repetitorId]; !ok {
		return errors.New("repetitor not found")
	}
	r.authRepository.ChangePassword(repetitorId, authData, newPassword)
	return nil
}

func (r *TestRepetiorRepository) GetRepetitorsIds(repetitorsOffset int64, repetitorsLimit int64) ([]int64, error) {
	return nil, nil
}

type TestClientRepository struct {
	data                   map[int64]*types.DBClientData
	personalDataRepository data_base.IPersonalDataRepository
	authRepository         data_base.IAuthRepository
	userRepository         data_base.IUserRepository
}

func CreateTestClientRepository(
	personalDataRepository data_base.IPersonalDataRepository,
	authRepository data_base.IAuthRepository,
	userRepository data_base.IUserRepository,
) *TestClientRepository {
	return &TestClientRepository{
		data:                   make(map[int64]*types.DBClientData),
		personalDataRepository: personalDataRepository,
		authRepository:         authRepository,
		userRepository:         userRepository,
	}
}

func (r *TestClientRepository) InsertClient(client types.DBClientData, personalData types.DBPersonalData, auth types.DBAuthData) (int64, error) {
	personalDataID, err := r.personalDataRepository.InsertPersonalData(personalData)
	if err != nil {
		return 0, err
	}
	userID, err := r.userRepository.InsertUser(types.DBUserData{
		PersonalDataID:   personalDataID,
		RegistrationDate: time.Now(),
		LastLoginDate:    time.Now(),
	})
	if err != nil {
		return 0, err
	}
	_, err = r.authRepository.InsertAuth(types.DBAuthInfo{
		UserID:   userID,
		UserType: types.Client,
		Login:    auth.Login,
		Password: auth.Password,
	})
	if err != nil {
		return 0, err
	}
	client.ID = userID
	r.data[client.ID] = &client
	return client.ID, nil
}

func (r *TestClientRepository) GetClient(id int64) (*types.DBClientData, error) {
	value, ok := r.data[id]
	if !ok {
		return nil, errors.New("client not found")
	}
	return value, nil
}

func (r *TestClientRepository) UpdateClientPersonalData(clientId int64, personalData types.DBPersonalData) error {
	if _, ok := r.data[clientId]; !ok {
		return errors.New("client not found")
	}
	r.personalDataRepository.UpdatePersonalData(clientId, personalData)
	return nil
}

func (r *TestClientRepository) UpdateClientPassword(clientId int64, authData types.DBAuthData, newPassword string) error {
	if _, ok := r.data[clientId]; !ok {
		return errors.New("client not found")
	}
	r.authRepository.ChangePassword(clientId, authData, newPassword)
	return nil
}

type TestContractRepository struct {
	data map[int64]*types.DBContract
}

func CreateTestContractRepository() *TestContractRepository {
	return &TestContractRepository{
		data: make(map[int64]*types.DBContract),
	}
}

func (r *TestContractRepository) GetContracts(clientID int64, repetitorID int64, from int64, size int64) ([]types.DBContract, error) {
	contracts := make([]types.DBContract, 0)
	for _, contract := range r.data {
		if contract.ClientID == clientID && contract.RepetitorID == repetitorID {
			contracts = append(contracts, *contract)
		}
	}
	return contracts[from:min(from+size, int64(len(contracts)))], nil
}

func (r *TestContractRepository) InsertContract(contract types.DBContract) (int64, error) {
	contract.ID = int64(len(r.data) + 1)
	r.data[contract.ID] = &contract
	return contract.ID, nil
}

func (r *TestContractRepository) GetContract(id int64) (*types.DBContract, error) {
	value, ok := r.data[id]
	if !ok {
		return nil, errors.New("contract not found")
	}
	return value, nil
}

func (r *TestContractRepository) UpdateContract(contract types.DBContract) error {
	if _, ok := r.data[contract.ID]; !ok {
		return errors.New("contract not found")
	}
	r.data[contract.ID] = &contract
	return nil
}

func (r *TestContractRepository) DeleteContract(id int64) error {
	if _, ok := r.data[id]; !ok {
		return errors.New("contract not found")
	}
	delete(r.data, id)
	return nil
}

func (r *TestContractRepository) GetContractsByRepetitorID(repetitorID int64, from int64, size int64, status types.ContractStatus) ([]types.DBContract, error) {
	contracts := make([]types.DBContract, 0)
	for _, contract := range r.data {
		if contract.RepetitorID == repetitorID && contract.Status == status {
			contracts = append(contracts, *contract)
		}
	}
	return contracts[from:min(from+size, int64(len(contracts)))], nil
}

func (r *TestContractRepository) GetContractsByClientID(clientID int64, from int64, size int64, status types.ContractStatus) ([]types.DBContract, error) {
	contracts := make([]types.DBContract, 0)
	for _, contract := range r.data {
		if contract.ClientID == clientID && contract.Status == status {
			contracts = append(contracts, *contract)
		}
	}
	return contracts[from:min(from+size, int64(len(contracts)))], nil
}

func (r *TestContractRepository) UpdateContractStatus(id int64, status types.ContractStatus) error {
	if _, ok := r.data[id]; !ok {
		return errors.New("contract not found")
	}
	r.data[id].Status = status
	return nil
}

func (r *TestContractRepository) UpdateContractRepetitorID(contractID int64, repetitorID int64) error {
	if _, ok := r.data[contractID]; !ok {
		return errors.New("contract not found")
	}
	r.data[contractID].RepetitorID = repetitorID
	return nil
}

func (r *TestContractRepository) UpdateContractReviewClientIDInSeq(tx *sql.Tx, id int64, reviewClientID int64) error {
	return r.UpdateContractReviewClientID(id, reviewClientID)
}

func (r *TestContractRepository) UpdateContractPaymentStatus(id int64, paymentStatus types.PaymentStatus) error {
	if _, ok := r.data[id]; !ok {
		return errors.New("contract not found")
	}
	r.data[id].PaymentStatus = paymentStatus
	return nil
}

func (r *TestContractRepository) UpdateContractReviewClientID(id int64, reviewClientID int64) error {
	if _, ok := r.data[id]; !ok {
		return errors.New("contract not found")
	}
	r.data[id].ReviewClientID = reviewClientID
	return nil
}

func (r *TestContractRepository) UpdateContractReviewRepetitorID(id int64, reviewRepetitorID int64) error {
	if _, ok := r.data[id]; !ok {
		return errors.New("contract not found")
	}
	r.data[id].ReviewRepetitorID = reviewRepetitorID
	return nil
}

func (r *TestContractRepository) GetContractList(from int64, size int64, status types.ContractStatus) ([]types.DBContract, error) {
	contracts := make([]types.DBContract, 0)
	for _, contract := range r.data {
		if contract.Status == status {
			contracts = append(contracts, *contract)
		}
	}
	return contracts[from:min(from+size, int64(len(contracts)))], nil
}

func (r *TestContractRepository) GetAllContracts(from int64, size int64) ([]types.DBContract, error) {
	contracts := make([]types.DBContract, 0)
	for _, contract := range r.data {
		contracts = append(contracts, *contract)
	}
	return contracts[from:min(from+size, int64(len(contracts)))], nil
}

func (r *TestContractRepository) BeginTx() (*sql.Tx, error) {
	return nil, nil
}

type TestContractReviewRepository struct {
	data map[int64]*types.DBReview
}

func CreateTestContractReviewRepository() *TestContractReviewRepository {
	return &TestContractReviewRepository{
		data: make(map[int64]*types.DBReview),
	}
}

type TestResumeRepository struct {
	data map[int64]*types.DBResume
}

func CreateTestResumeRepository() *TestResumeRepository {
	return &TestResumeRepository{
		data: make(map[int64]*types.DBResume),
	}
}

func (r *TestResumeRepository) InsertResume(resume types.DBResume) (int64, error) {
	resume.ID = int64(len(r.data) + 1)
	r.data[resume.ID] = &resume
	return resume.ID, nil
}

func (r *TestResumeRepository) InsertResumeInSeq(tx *sql.Tx, resume types.DBResume) (int64, error) {
	return r.InsertResume(resume)
}

func (r *TestResumeRepository) GetResume(id int64) (*types.DBResume, error) {
	value, ok := r.data[id]
	if !ok {
		return nil, errors.New("resume not found")
	}
	return value, nil
}

func (r *TestResumeRepository) UpdateResumeTitle(id int64, title string, updatedAt time.Time) error {
	if _, ok := r.data[id]; !ok {
		return errors.New("resume not found")
	}
	r.data[id].Title = title
	r.data[id].UpdatedAt = updatedAt
	return nil
}

func (r *TestResumeRepository) UpdateResumeDescription(id int64, description string, updatedAt time.Time) error {
	if _, ok := r.data[id]; !ok {
		return errors.New("resume not found")
	}
	r.data[id].Description = description
	r.data[id].UpdatedAt = updatedAt
	return nil
}

func (r *TestResumeRepository) UpdateResumePrices(id int64, prices map[string]int, updatedAt time.Time) error {
	if _, ok := r.data[id]; !ok {
		return errors.New("resume not found")
	}
	r.data[id].Prices = prices
	r.data[id].UpdatedAt = updatedAt
	return nil
}

func (r *TestResumeRepository) DeleteResume(id int64) error {
	if _, ok := r.data[id]; !ok {
		return errors.New("resume not found")
	}
	delete(r.data, id)
	return nil
}

type TestReviewRepository struct {
	data map[int64]*types.DBReview
}

func CreateTestReviewRepository() *TestReviewRepository {
	return &TestReviewRepository{
		data: make(map[int64]*types.DBReview),
	}
}

func (r *TestReviewRepository) InsertReview(review types.DBReview) (int64, error) {
	review.ID = int64(len(r.data) + 1)
	r.data[review.ID] = &review
	return review.ID, nil
}

func (r *TestReviewRepository) GetReview(id int64) (*types.DBReview, error) {
	value, ok := r.data[id]
	if !ok {
		return nil, errors.New("review not found")
	}
	return value, nil
}

func (r *TestReviewRepository) UpdateReview(review types.DBReview) error {
	if _, ok := r.data[review.ID]; !ok {
		return errors.New("review not found")
	}
	r.data[review.ID] = &review
	return nil
}

func (r *TestReviewRepository) GetReviewsByRepetitorID(repetitorID int64, from int64, size int64) ([]types.DBReview, error) {
	reviews := make([]types.DBReview, 0)
	for _, review := range r.data {
		if review.RepetitorID == repetitorID {
			reviews = append(reviews, *review)
		}
	}
	return reviews[from:min(from+size, int64(len(reviews)))], nil
}

func (r *TestReviewRepository) GetReviewsByClientID(clientID int64, from int64, size int64) ([]types.DBReview, error) {
	reviews := make([]types.DBReview, 0)
	for _, review := range r.data {
		if review.ClientID == clientID {
			reviews = append(reviews, *review)
		}
	}
	return reviews[from:min(from+size, int64(len(reviews)))], nil
}

func (r *TestReviewRepository) InsertReviewInSeq(tx *sql.Tx, review types.DBReview) (int64, error) {
	return r.InsertReview(review)
}

func (r *TestReviewRepository) BeginTx() (*sql.Tx, error) {
	return nil, nil
}

type TestLessonRepository struct {
	data map[int64]*types.DBLesson
}

func CreateTestLessonRepository() *TestLessonRepository {
	return &TestLessonRepository{
		data: make(map[int64]*types.DBLesson),
	}
}

func (r *TestLessonRepository) UpdateLesson(lessonID int64, duration *int64, format *string) error {
	if _, ok := r.data[lessonID]; !ok {
		return errors.New("lesson not found")
	}
	r.data[lessonID].Duration = *duration
	r.data[lessonID].Format = *format
	return nil
}

func (r *TestLessonRepository) DeleteLesson(lessonID int64) error {
	if _, ok := r.data[lessonID]; !ok {
		return errors.New("lesson not found")
	}
	delete(r.data, lessonID)
	return nil
}

func (r *TestLessonRepository) InsertLesson(lesson types.DBLesson) (int64, error) {
	if lesson.ContractID == 0 {
		return 0, errors.New("contract id is required")
	}
	lesson.ID = int64(len(r.data) + 1)
	r.data[lesson.ID] = &lesson
	return lesson.ID, nil
}

func (r *TestLessonRepository) GetLesson(id int64) (*types.DBLesson, error) {
	value, ok := r.data[id]
	if !ok {
		return nil, errors.New("lesson not found")
	}
	return value, nil
}

func (r *TestLessonRepository) GetLessons(contractID int64, from int64, size int64) ([]types.DBLesson, error) {
	lessons := make([]types.DBLesson, 0)
	for _, lesson := range r.data {
		if lesson.ContractID == contractID {
			lessons = append(lessons, *lesson)
		}
	}
	return lessons[from:min(from+size, int64(len(lessons)))], nil
}

func (r *TestLessonRepository) InsertLessonInSeq(tx *sql.Tx, lesson types.DBLesson) (int64, error) {
	return r.InsertLesson(lesson)
}

func (r *TestLessonRepository) BeginTx() (*sql.Tx, error) {
	return nil, nil
}

type TestTransactionRepository struct {
	data map[int64]*types.DBTransaction
}

func CreateTestTransactionRepository() *TestTransactionRepository {
	return &TestTransactionRepository{
		data: make(map[int64]*types.DBTransaction),
	}
}

func (r *TestTransactionRepository) GetContractTransactionsList(contract_id int64, from int64, size int64) ([]types.DBTransaction, error) {
	transactions := make([]types.DBTransaction, 0)
	for _, transaction := range r.data {
		if transaction.ContractID == contract_id {
			transactions = append(transactions, *transaction)
		}
	}
	return transactions[from:min(from+size, int64(len(transactions)))], nil
}

func (r *TestTransactionRepository) InsertTransaction(transaction types.DBTransaction) (int64, error) {
	transaction.ID = int64(len(r.data) + 1)
	r.data[transaction.ID] = &transaction
	return transaction.ID, nil
}

func (r *TestTransactionRepository) GetTransaction(id int64) (*types.DBTransaction, error) {
	value, ok := r.data[id]
	if !ok {
		return nil, errors.New("transaction not found")
	}
	return value, nil
}

func (r *TestTransactionRepository) GetTransactionsList(userId int64, from int64, size int64) ([]types.DBTransaction, error) {
	transactions := make([]types.DBTransaction, 0)
	for _, transaction := range r.data {
		if transaction.UserID == userId {
			transactions = append(transactions, *transaction)
		}
	}
	return transactions[from:min(from+size, int64(len(transactions)))], nil
}

func (r *TestTransactionRepository) GetPendingContractPaymentTransaction() (*types.DBPendingContractPaymentTransaction, error) {
	for _, transaction := range r.data {
		if transaction.Type == types.TransactionTypeContractPayment && transaction.Status == types.TransactionStatusPending {
			return &types.DBPendingContractPaymentTransaction{
				ID:            transaction.ID,
				UserID:        transaction.UserID,
				Amount:        transaction.Amount,
				CreatedAt:     transaction.CreatedAt,
				TransactionID: transaction.ID,
			}, nil
		}
	}
	return nil, errors.New("pending contract payment transaction not found")
}

func (r *TestTransactionRepository) InsertPendingContractPaymentTransaction(transactionPendingContractPayment types.DBPendingContractPaymentTransaction, transaction types.DBTransaction) (int64, error) {
	return r.InsertTransaction(transaction)
}

func (r *TestTransactionRepository) UpdateTransactionStatus(transactionId int64, status types.TransactionStatus) error {
	if _, ok := r.data[transactionId]; !ok {
		return errors.New("transaction not found")
	}
	r.data[transactionId].Status = status
	return nil
}

func (r *TestTransactionRepository) ApproveTransaction(transactionId int64) error {
	if _, ok := r.data[transactionId]; !ok {
		return errors.New("transaction not found")
	}
	r.data[transactionId].Status = types.TransactionStatusPaid
	r.data[transactionId].CreatedAt = time.Now()
	return nil
}

type TestDepartmentRepository struct {
	data map[int64]*types.DBDepartment
}

func CreateTestDepartmentRepository() *TestDepartmentRepository {
	return &TestDepartmentRepository{
		data: make(map[int64]*types.DBDepartment),
	}
}

func (r *TestDepartmentRepository) UpdateDepartmentName(departmentID int64, name string) error {
	if _, ok := r.data[departmentID]; !ok {
		return errors.New("department not found")
	}
	r.data[departmentID].Name = name
	return nil
}

func (r *TestDepartmentRepository) DeleteDepartment(departmentID int64) error {
	if _, ok := r.data[departmentID]; !ok {
		return errors.New("department not found")
	}
	delete(r.data, departmentID)
	return nil
}

func (r *TestDepartmentRepository) GetModeratorsByDepartmentID(departmentID int64) ([]int64, error) {
	return nil, nil
}

func (r *TestDepartmentRepository) InsertDepartment(department types.DBDepartment) (int64, error) {
	department.ID = int64(len(r.data) + 1)
	r.data[department.ID] = &department
	return department.ID, nil
}

func (r *TestDepartmentRepository) GetDepartment(id int64) (*types.DBDepartment, error) {
	value, ok := r.data[id]
	if !ok {
		return nil, errors.New("department not found")
	}
	return value, nil
}

func (r *TestDepartmentRepository) GetDepartmentsByHeadID(headID int64) ([]types.DBDepartment, error) {
	departments := make([]types.DBDepartment, 0)
	for _, department := range r.data {
		if department.HeadID == headID {
			departments = append(departments, *department)
		}
	}
	return departments, nil
}

func (r *TestDepartmentRepository) GetDepartmentIdByName(name string) (int64, error) {
	for _, department := range r.data {
		if department.Name == name {
			return department.ID, nil
		}
	}
	return 0, nil
}

func (r *TestDepartmentRepository) GetDepartmentUsersIDs(id int64) ([]int64, error) {
	users := make([]int64, 0)
	for _, department := range r.data {
		if department.HeadID == id {
			users = append(users, department.HeadID)
		}
	}
	return users, nil
}

func (r *TestDepartmentRepository) ChangeDepartmentHead(id int64, headID int64) error {
	if _, ok := r.data[id]; !ok {
		return errors.New("department not found")
	}
	r.data[id].HeadID = headID
	return nil
}

func (r *TestDepartmentRepository) HireInfoInsert(hireInfo types.DBHireInfo) error {
	return nil
}

func (r *TestDepartmentRepository) HireInfoDelete(userId int64, departmentId int64) error {
	return nil
}

func (r *TestDepartmentRepository) GetUserDepartmentsIDs(userId int64) ([]int64, error) {
	return nil, nil
}

type TestChatRepository struct {
	data map[int64]*types.DBChat
}

func CreateTestChatRepository() *TestChatRepository {
	return &TestChatRepository{
		data: make(map[int64]*types.DBChat),
	}
}

func (r *TestChatRepository) UpdateChat(chatID int64, chatStatus string) error {
	if _, ok := r.data[chatID]; !ok {
		return errors.New("chat not found")
	}
	r.data[chatID].Status = chatStatus
	return nil
}

func (r *TestChatRepository) DeleteChat(id int64) error {
	if _, ok := r.data[id]; !ok {
		return errors.New("chat not found")
	}
	delete(r.data, id)
	return nil
}

func (r *TestChatRepository) InsertChat(chat types.DBChat) (int64, error) {
	chat.ID = int64(len(r.data) + 1)
	r.data[chat.ID] = &chat
	return chat.ID, nil
}

func (r *TestChatRepository) GetChat(id int64) (*types.DBChat, error) {
	value, ok := r.data[id]
	if !ok {
		return nil, errors.New("chat not found")
	}
	return value, nil
}

func (r *TestMessageRepository) DeleteMessages(chatID int64) error {
	if _, ok := r.data[chatID]; !ok {
		return errors.New("messages not found")
	}
	delete(r.data, chatID)
	return nil
}

func (r *TestChatRepository) GetChatListByClientID(clientID int64, from int64, size int64) ([]types.DBChat, error) {
	chats := make([]types.DBChat, 0)
	for _, chat := range r.data {
		if chat.ClientID == clientID {
			chats = append(chats, *chat)
		}
	}
	return chats, nil
}

func (r *TestChatRepository) GetChatListByRepetitorID(repetitorID int64, from int64, size int64) ([]types.DBChat, error) {
	chats := make([]types.DBChat, 0)
	for _, chat := range r.data {
		if chat.RepetitorID == repetitorID {
			chats = append(chats, *chat)
		}
	}
	return chats, nil
}

func (r *TestChatRepository) GetChatListByModeratorID(moderatorID int64, from int64, size int64) ([]types.DBChat, error) {
	chats := make([]types.DBChat, 0)
	for _, chat := range r.data {
		if chat.ModeratorID == moderatorID {
			chats = append(chats, *chat)
		}
	}
	return chats, nil
}

func (r *TestChatRepository) GetChatIdByCIDAndMID(clientID int64, moderatorID int64) (int64, error) {
	for _, chat := range r.data {
		if chat.ClientID == clientID && chat.ModeratorID == moderatorID {
			return chat.ID, nil
		}
	}
	return 0, nil
}

func (r *TestChatRepository) GetChatIdByCIDAndRID(clientID int64, repetitorID int64) (int64, error) {
	for _, chat := range r.data {
		if chat.ClientID == clientID && chat.RepetitorID == repetitorID {
			return chat.ID, nil
		}
	}
	return 0, nil
}

func (r *TestChatRepository) GetChatIdByMIDAndRID(moderatorID int64, repetitorID int64) (int64, error) {
	for _, chat := range r.data {
		if chat.ModeratorID == moderatorID && chat.RepetitorID == repetitorID {
			return chat.ID, nil
		}
	}
	return 0, nil
}

type TestMessageRepository struct {
	data map[int64]*types.DBMessage
}

func CreateTestMessageRepository() *TestMessageRepository {
	return &TestMessageRepository{
		data: make(map[int64]*types.DBMessage),
	}
}

func (r *TestMessageRepository) UpdateMessageContent(messageID int64, content string) error {
	if _, ok := r.data[messageID]; !ok {
		return errors.New("message not found")
	}
	r.data[messageID].Content = content
	return nil
}

func (r *TestMessageRepository) GetMessage(messageID int64) (*types.DBMessage, error) {
	if _, ok := r.data[messageID]; !ok {
		return nil, errors.New("message not found")
	}
	return r.data[messageID], nil
}

func (r *TestMessageRepository) DeleteMessage(messageID int64) error {
	if _, ok := r.data[messageID]; !ok {
		return errors.New("message not found")
	}
	delete(r.data, messageID)
	return nil
}

func (r *TestMessageRepository) InsertMessage(message types.DBMessage) (int64, error) {
	message.ID = int64(len(r.data) + 1)
	r.data[message.ID] = &message
	return message.ID, nil
}

func (r *TestMessageRepository) GetMessages(chatID int64, from int64, size int64) ([]types.DBMessage, error) {
	messages := make([]types.DBMessage, 0)
	for _, message := range r.data {
		if message.ChatID == chatID {
			messages = append(messages, *message)
		}
	}
	return messages, nil
}
