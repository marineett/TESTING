package console_module

import (
	"data_base_project/data_base"
	"database/sql"
	"fmt"
	"strconv"
)

func MainMenu(sqlDataBaseModule *data_base.DataBaseModule, db *sql.DB) {

	serviceModule := ChangeDataBaseType(sqlDataBaseModule)
	for {
		fmt.Println("Main menu")
		fmt.Println("1. User registration")
		fmt.Println("2. Contract functionality")
		fmt.Println("3. Department functionality")
		fmt.Println("4. Transaction functionality")
		fmt.Println("5. Chat functionality")
		fmt.Println("6. Admin functionality")
		fmt.Println("7. Client functionality")
		fmt.Println("8. Repetitor functionality")
		fmt.Println("9. Moderator functionality")
		fmt.Println("10. Bench mark")
		fmt.Println("11. Exit")
		choiceStr := ""
		fmt.Scanln(&choiceStr)
		choice, err := strconv.Atoi(choiceStr)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if choice < 1 || choice > 13 {
			fmt.Println("Invalid choice")
			continue
		}
		switch choice {
		case 1:
			UserRegistrationWork(serviceModule)
		case 2:
			ContractWork(serviceModule)
		case 3:
			DepartmentWork(serviceModule)
		case 4:
			TransactionWork(serviceModule)
		case 5:
			ChatWork(serviceModule)
		case 6:
			AdminWork(serviceModule)
		case 7:
			ClientWork(serviceModule)
		case 8:
			RepetitorWork(serviceModule)
		case 9:
			ModeratorWork(serviceModule)
		case 10:
			BenchMarkMenu(serviceModule, db)
		case 11:
			return
		}
	}
}
