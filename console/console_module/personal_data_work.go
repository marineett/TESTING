package console_module

import (
	"data_base_project/types"
	"fmt"
)

func PrintPersonalData(personalData *types.ServicePersonalData) {
	fmt.Println("Personal data:")
	fmt.Println("First name:", personalData.FirstName)
	fmt.Println("Last name:", personalData.LastName)
	fmt.Println("Middle name:", personalData.MiddleName)
	fmt.Println("Telephone number:", personalData.TelephoneNumber)
	fmt.Println("Email:", personalData.Email)
	fmt.Println("Passport number:", personalData.PassportNumber)
	fmt.Println("Passport series:", personalData.PassportSeries)
	fmt.Println("Passport date:", personalData.PassportDate)
	fmt.Println("Passport issued by:", personalData.PassportIssuedBy)
}
