package console_module

import (
	"bufio"
	"data_base_project/service_logic"
	"data_base_project/types"
	"fmt"
	"os"
	"strconv"
	"time"
)

func CreateChat(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Create chat")
	fmt.Println("Enter first user type:")
	fmt.Println("1. Client")
	fmt.Println("2. Repetitor")
	fmt.Println("3. Moderator")
	choiceStr := ""
	fmt.Scanln(&choiceStr)
	choice, err := strconv.Atoi(choiceStr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if choice < 1 || choice > 3 {
		fmt.Println("Invalid choice")
		return
	}
	fmt.Println("Enter first user id:")
	firstUserIDStr := ""
	fmt.Scanln(&firstUserIDStr)
	firstUserID, err := strconv.ParseInt(firstUserIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Enter second user type:")
	fmt.Println("1. Client")
	fmt.Println("2. Repetitor")
	fmt.Println("3. Moderator")
	choice2Str := ""
	fmt.Scanln(&choice2Str)
	choice2, err := strconv.Atoi(choice2Str)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if choice2 < 1 || choice2 > 3 {
		fmt.Println("Invalid choice")
		return
	}
	fmt.Println("Enter second user id:")
	secondUserIDStr := ""
	fmt.Scanln(&secondUserIDStr)
	secondUserID, err := strconv.ParseInt(secondUserIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if choice == 1 {
		if choice2 == 1 {
			fmt.Println("Cannot create chat between client and client")
		} else if choice2 == 2 {
			chatID, err := serviceModule.ChatService.CreateCRChat(firstUserID, secondUserID)
			if err != nil {
				fmt.Println("Error creating chat:", err)
				return
			}
			fmt.Println("Chat created successfully with id:", chatID)
		} else if choice2 == 3 {
			chatID, err := serviceModule.ChatService.CreateCMChat(firstUserID, secondUserID)
			if err != nil {
				fmt.Println("Error creating chat:", err)
				return
			}
			fmt.Println("Chat created successfully with id:", chatID)
		}
	} else if choice == 2 {
		if choice2 == 1 {
			chatID, err := serviceModule.ChatService.CreateCRChat(secondUserID, firstUserID)
			if err != nil {
				fmt.Println("Error creating chat:", err)
				return
			}
			fmt.Println("Chat created successfully with id:", chatID)
		} else if choice2 == 2 {
			fmt.Println("Cannot create chat between repetitor and repetitor")
		} else if choice2 == 3 {
			chatID, err := serviceModule.ChatService.CreateRMChat(firstUserID, secondUserID)
			if err != nil {
				fmt.Println("Error creating chat:", err)
				return
			}
			fmt.Println("Chat created successfully with id:", chatID)
		}
	} else if choice == 3 {
		if choice2 == 1 {
			chatID, err := serviceModule.ChatService.CreateCMChat(secondUserID, firstUserID)
			if err != nil {
				fmt.Println("Error creating chat:", err)
				return
			}
			fmt.Println("Chat created successfully with id:", chatID)
		} else if choice2 == 2 {
			chatID, err := serviceModule.ChatService.CreateRMChat(secondUserID, firstUserID)
			if err != nil {
				fmt.Println("Error creating chat:", err)
				return
			}
			fmt.Println("Chat created successfully with id:", chatID)
		}
	}
}

func PrintMessage(message *types.ServiceMessage) {
	fmt.Println("Message id:", message.ID)
	fmt.Println("Message content:", message.Content)
	fmt.Println("Message sender id:", message.SenderID)
	fmt.Println("Message created at:", message.CreatedAt.Format(time.RFC3339))
}

func GetChatHistory(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Get chat history")
	fmt.Println("Enter chat id:")
	chatIDStr := ""
	fmt.Scanln(&chatIDStr)
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Enter messages offset:")
	offsetStr := ""
	fmt.Scanln(&offsetStr)
	offset, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Enter messages limit:")
	limitStr := ""
	fmt.Scanln(&limitStr)
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	history, err := serviceModule.ChatService.GetMessages(chatID, offset, limit)
	if err != nil {
		fmt.Println("Error getting chat history:", err)
		return
	}
	fmt.Println("Chat history:")
	for _, message := range history {
		PrintMessage(&message)
	}
}

func SendMessage(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Send message")
	fmt.Println("Enter chat id:")
	chatIDStr := ""
	fmt.Scanln(&chatIDStr)
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Enter user id:")
	userIDStr := ""
	fmt.Scanln(&userIDStr)
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	fmt.Println("Enter message content:")
	scanner := bufio.NewScanner(os.Stdin)
	var content string
	if scanner.Scan() {
		input := scanner.Text()
		content = input
	}
	err = serviceModule.ChatService.SendMessage(chatID, userID, content)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}
	fmt.Println("Message sent successfully")
}

func ChatWork(serviceModule *service_logic.ServiceModule) {
	for {
		fmt.Println("Chat work")
		fmt.Println("1. Create chat")
		fmt.Println("2. Get chat history")
		fmt.Println("3. Send message")
		fmt.Println("4. Exit")
		choiceStr := ""
		fmt.Scanln(&choiceStr)
		choice, err := strconv.Atoi(choiceStr)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		switch choice {
		case 1:
			CreateChat(serviceModule)
		case 2:
			GetChatHistory(serviceModule)
		case 3:
			SendMessage(serviceModule)
		case 4:
			return
		}
	}
}
