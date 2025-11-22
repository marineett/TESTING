package console_module

import (
	"data_base_project/service_logic"
	"data_base_project/types"
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

func GenerateUsers(size int64, serviceModule *service_logic.ServiceModule) {
	for i := int64(0); i < size; i++ {
		initClient := types.ServiceInitClientData{
			ServiceInitUserData: types.ServiceInitUserData{
				ServicePersonalData: types.ServicePersonalData{
					FirstName:       fmt.Sprintf("user%d", i),
					LastName:        fmt.Sprintf("user%d", i),
					Email:           fmt.Sprintf("user%d@example.com", i),
					TelephoneNumber: fmt.Sprintf("+7999999999%d", i),
				},
				ServiceAuthData: types.ServiceAuthData{
					Login:    fmt.Sprintf("user%d", i),
					Password: "password",
				},
			},
		}
		serviceModule.ClientService.CreateClient(initClient)
	}
}

func BenchMark(size int64, serviceModule *service_logic.ServiceModule, db *sql.DB) {
	GenerateUsers(size, serviceModule)
	start := time.Now()
	for i := int64(0); i < size; i++ {
		generateUserIndex := rand.Intn(int(size))
		serviceModule.AuthService.Authorize(types.ServiceAuthData{
			Login:    fmt.Sprintf("user%d", generateUserIndex),
			Password: "password",
		})
	}
	elapsed := time.Since(start)
	fmt.Printf("Time taken without index with %d users: %d ms\n", size, elapsed.Milliseconds())
	start = time.Now()
	for i := int64(0); i < size; i++ {
		generateUserIndex := rand.Intn(int(size))
		serviceModule.AuthService.Authorize(types.ServiceAuthData{
			Login:    fmt.Sprintf("user%d", generateUserIndex),
			Password: "password",
		})
	}
	elapsed = time.Since(start)
	fmt.Printf("Time taken with index with %d users: %d ms\n", size, elapsed.Milliseconds())
}

func BenchMarkMenu(serviceModule *service_logic.ServiceModule, db *sql.DB) {
	fmt.Println("Enter size of users:")
	var size int64
	fmt.Scanln(&size)
	BenchMark(size, serviceModule, db)
}
