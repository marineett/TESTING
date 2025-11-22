package console_module

import (
	"bufio"
	"data_base_project/service_logic"
	"data_base_project/types"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func CreateDepartment(serviceModule *service_logic.ServiceModule) {
	department := types.ServiceDepartmentInitData{}
	fmt.Println("Enter department name:")
	buffer := bufio.NewReader(os.Stdin)
	department.Name, _ = buffer.ReadString('\n')
	department.HeadID = 0
	id, err := serviceModule.DepartmentService.CreateDepartment(department)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Department created with ID:", id)
}

func AssignModeratorToDepartment(serviceModule *service_logic.ServiceModule) {
	var departmentID int64
	fmt.Println("Enter department ID:")
	departmentIDStr := ""
	fmt.Scanln(&departmentIDStr)
	departmentID, err := strconv.ParseInt(departmentIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	var moderatorID int64
	fmt.Println("Enter moderator ID:")
	moderatorIDStr := ""
	fmt.Scanln(&moderatorIDStr)
	moderatorID, err = strconv.ParseInt(moderatorIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	err = serviceModule.DepartmentService.AssignModeratorToDepartment(moderatorID, departmentID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Moderator assigned to department")
}

func AssignAdminToDepartment(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter department ID:")
	departmentIDStr := ""
	fmt.Scanln(&departmentIDStr)
	departmentID, err := strconv.ParseInt(departmentIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Enter admin ID:")
	adminIDStr := ""
	fmt.Scanln(&adminIDStr)
	adminID, err := strconv.ParseInt(adminIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	err = serviceModule.DepartmentService.AssignAdminToDepartment(adminID, departmentID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Admin assigned to department")
}

func GetDepartmentIdByName(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter department name:")
	buffer := bufio.NewReader(os.Stdin)
	name, _ := buffer.ReadString('\n')
	id, err := serviceModule.DepartmentService.GetDepartmentIdByName(name)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Department ID:", id)
}

func PrintDepartment(department *types.ServiceDepartment) {
	fmt.Println("Department info:")
	fmt.Println("ID:", department.ID)
	fmt.Println("Name:", department.Name)
	fmt.Println("Head ID:", department.HeadID)
}

func GetDepartmentById(serviceModule *service_logic.ServiceModule) {
	var idStr string
	fmt.Println("Enter department ID:")
	fmt.Scanln(&idStr)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	department, err := serviceModule.DepartmentService.GetDepartment(id)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	PrintDepartment(&department)
}

func ListAdminDepartments(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter admin ID:")
	adminIDStr := ""
	fmt.Scanln(&adminIDStr)
	adminID, err := strconv.ParseInt(adminIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	departments, err := serviceModule.DepartmentService.GetDepartmentsByHeadIdWithModerators(adminID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if len(departments) == 0 {
		fmt.Println("No departments found")
		return
	}
	for _, d := range departments {
		fmt.Println("--- Department ---")
		fmt.Println("ID:", d.ID)
		fmt.Println("Name:", d.Name)
		fmt.Println("Head ID:", d.HeadID)
		if len(d.Moderators) == 0 {
			fmt.Println("Moderators: <none>")
		} else {
			fmt.Println("Moderators:")
			for _, m := range d.Moderators {
				fmt.Printf("  - ID: %d, Name: %s %s\n", m.ID, m.Moderator.FirstName, m.Moderator.LastName)
			}
		}
	}
}

func ReplaceDepartmentInfo(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter department ID:")
	deptIDStr := ""
	fmt.Scanln(&deptIDStr)
	deptID, err := strconv.ParseInt(deptIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Enter new department name:")
	reader := bufio.NewReader(os.Stdin)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	if name == "" {
		fmt.Println("Name cannot be empty")
		return
	}
	if err := serviceModule.DepartmentService.UpdateDepartmentName(deptID, name); err != nil {
		fmt.Println("Error:", err)
		return
	}
	dept, err := serviceModule.DepartmentService.GetDepartment(deptID)
	if err != nil {
		fmt.Println("Updated, but failed to fetch department:", err)
		return
	}
	PrintDepartment(&dept)
}

func DeleteDepartment(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter department ID:")
	deptIDStr := ""
	fmt.Scanln(&deptIDStr)
	deptID, err := strconv.ParseInt(deptIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if err := serviceModule.DepartmentService.DeleteDepartment(deptID); err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Department deleted")
}

func RemoveModeratorFromDepartment(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter department ID:")
	deptIDStr := ""
	fmt.Scanln(&deptIDStr)
	deptID, err := strconv.ParseInt(deptIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Enter moderator ID:")
	moderatorIDStr := ""
	fmt.Scanln(&moderatorIDStr)
	moderatorID, err := strconv.ParseInt(moderatorIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if err := serviceModule.DepartmentService.FireModeratorFromDepartment(moderatorID, deptID); err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Moderator removed from department")
}

func DepartmentWork(serviceModule *service_logic.ServiceModule) {
	for {
		fmt.Println("Department work")
		fmt.Println("1. Create department")
		fmt.Println("2. Assign moderator to department")
		fmt.Println("3. Assign admin to department")
		fmt.Println("4. Get department ID by name")
		fmt.Println("5. Get department by ID")
		fmt.Println("6. List admin departments")
		fmt.Println("7. Replace department info")
		fmt.Println("8. Delete department")
		fmt.Println("9. Remove moderator from department")
		fmt.Println("10. Exit")
		choiceStr := ""
		fmt.Scanln(&choiceStr)
		choice, err := strconv.Atoi(choiceStr)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if choice < 1 || choice > 10 {
			fmt.Println("Invalid choice")
			continue
		}
		switch choice {
		case 1:
			CreateDepartment(serviceModule)
		case 2:
			AssignModeratorToDepartment(serviceModule)
		case 3:
			AssignAdminToDepartment(serviceModule)
		case 4:
			GetDepartmentIdByName(serviceModule)
		case 5:
			GetDepartmentById(serviceModule)
		case 6:
			ListAdminDepartments(serviceModule)
		case 7:
			ReplaceDepartmentInfo(serviceModule)
		case 8:
			DeleteDepartment(serviceModule)
		case 9:
			RemoveModeratorFromDepartment(serviceModule)
		case 10:
			return
		default:
			fmt.Println("Invalid choice")
		}
	}
}
