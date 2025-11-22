package main

import (
	"console2/ui"
	"fmt"
)

func main() {
	base := "http://127.0.0.1:8000"
	token := ""
	for {
		fmt.Println("Setup:")
		fmt.Printf("Base URL [%s]: ", base)
		if s := ui.PromptWithDefault("", base); s != base {
			base = s
		}
		runMainMenu(base, &token)
		exit := ui.PromptWithDefault("Exit? (y/n)", "y")
		if exit == "y" || exit == "Y" {
			return
		}
	}
}

func runMainMenu(base string, token *string) {
	for {
		fmt.Println("Main menu:")
		fmt.Println("1) Auth (login/register, set token)")
		fmt.Println("2) Contracts")
		fmt.Println("3) Chats")
		fmt.Println("4) Moderators")
		fmt.Println("5) Departments")
		fmt.Println("6) Admins")
		fmt.Println("7) Repetitors")
		fmt.Println("0) Exit")
		ch := ui.PromptIntRequired("choice")
		switch ch {
		case 0:
			return
		case 1:
			ui.RunAuthMenu(base, token)
		case 2:
			ui.RunContractsMenu(base, token)
		case 3:
			ui.RunChatsMenu(base, token)
		case 4:
			ui.RunModeratorsMenu(base, token)
		case 5:
			ui.RunDepartmentsMenu(base, token)
		case 6:
			ui.RunAdminsMenu(base, token)
		case 7:
			ui.RunRepetitorsMenu(base, token)
		default:
			fmt.Println("unknown choice")
		}
	}
}
