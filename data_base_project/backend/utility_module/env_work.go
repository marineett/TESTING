package utility_module

import "os"

func UnsetEnv() {
	os.Unsetenv("DATABASE_HOST")
	os.Unsetenv("DATABASE_PORT")
	os.Unsetenv("DATABASE_USER")
	os.Unsetenv("DATABASE_PASSWORD")
	os.Unsetenv("DATABASE_NAME")
	os.Unsetenv("USER_TABLE_NAME")
	os.Unsetenv("PERSONAL_DATA_TABLE_NAME")
	os.Unsetenv("CHAT_TABLE_NAME")
	os.Unsetenv("MESSAGE_TABLE_NAME")
	os.Unsetenv("DEPARTMENT_TABLE_NAME")
	os.Unsetenv("EMPLOYEE_TABLE_NAME")
	os.Unsetenv("CLIENT_TABLE_NAME")
	os.Unsetenv("TRANSACTION_TABLE_NAME")
	os.Unsetenv("REPEATITOR_TABLE_NAME")
	os.Unsetenv("RESUME_TABLE_NAME")
	os.Unsetenv("REVIEW_TABLE_NAME")
	os.Unsetenv("CONTRACT_TABLE_NAME")
	os.Unsetenv("ADMIN_TABLE_NAME")
	os.Unsetenv("MODERATOR_TABLE_NAME")
	os.Unsetenv("HIRE_INFO_TABLE_NAME")
	os.Unsetenv("BACKEND_PORT")
}
