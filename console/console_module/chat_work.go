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
			fmt.Println("Chat created successfully with ID:", chatID)
		} else if choice2 == 3 {
			chatID, err := serviceModule.ChatService.CreateCMChat(firstUserID, secondUserID)
			if err != nil {
				fmt.Println("Error creating chat:", err)
				return
			}
			fmt.Println("Chat created successfully with ID:", chatID)
		}
	} else if choice == 2 {
		if choice2 == 1 {
			chatID, err := serviceModule.ChatService.CreateCRChat(secondUserID, firstUserID)
			if err != nil {
				fmt.Println("Error creating chat:", err)
				return
			}
			fmt.Println("Chat created successfully with ID:", chatID)
		} else if choice2 == 2 {
			fmt.Println("Cannot create chat between repetitor and repetitor")
		} else if choice2 == 3 {
			chatID, err := serviceModule.ChatService.CreateRMChat(firstUserID, secondUserID)
			if err != nil {
				fmt.Println("Error creating chat:", err)
				return
			}
			fmt.Println("Chat created successfully with ID:", chatID)
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
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
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
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Enter message content:")
	scanner := bufio.NewScanner(os.Stdin)
	var content string
	if scanner.Scan() {
		input := scanner.Text()
		content = input
	}
	_, err = serviceModule.ChatService.SendMessage(chatID, userID, content)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}
	fmt.Println("Message sent successfully")
}

func ListUserChats(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter role:")
	fmt.Println("1. Client")
	fmt.Println("2. Repetitor")
	fmt.Println("3. Moderator")
	roleStr := ""
	fmt.Scanln(&roleStr)
	role, err := strconv.Atoi(roleStr)
	if err != nil || role < 1 || role > 3 {
		fmt.Println("Invalid role")
		return
	}
	fmt.Println("Enter user id:")
	userIDStr := ""
	fmt.Scanln(&userIDStr)
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Enter offset:")
	offsetStr := ""
	fmt.Scanln(&offsetStr)
	offset, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Enter limit:")
	limitStr := ""
	fmt.Scanln(&limitStr)
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	var chats []types.ServiceChat
	switch role {
	case 1:
		chats, err = serviceModule.ChatService.GetChatListByClientID(userID, offset, limit)
	case 2:
		chats, err = serviceModule.ChatService.GetChatListByRepetitorID(userID, offset, limit)
	case 3:
		chats, err = serviceModule.ChatService.GetChatListByModeratorID(userID, offset, limit)
	}
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if len(chats) == 0 {
		fmt.Println("No chats found")
		return
	}
	for _, c := range chats {
		fmt.Println("--- Chat ---")
		fmt.Println("ID:", c.ID)
		fmt.Println("ClientID:", c.ClientID)
		fmt.Println("RepetitorID:", c.RepetitorID)
		fmt.Println("ModeratorID:", c.ModeratorID)
		fmt.Println("Status:", c.Status)
		fmt.Println("Type:", c.Type)
		fmt.Println("CreatedAt:", c.CreatedAt.Format(time.RFC3339))
	}
}

func ReplaceChat(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter chat id:")
	chatIDStr := ""
	fmt.Scanln(&chatIDStr)
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Enter new chat status:")
	status := ""
	fmt.Scanln(&status)
	if err := serviceModule.ChatService.UpdateChat(chatID, status); err != nil {
		fmt.Println("Error:", err)
		return
	}
	chat, err := serviceModule.ChatService.GetChat(chatID)
	if err != nil {
		fmt.Println("Updated, but failed to fetch chat:", err)
		return
	}
	fmt.Println("Chat updated:")
	fmt.Println("ID:", chat.ID)
	fmt.Println("Status:", chat.Status)
}

func UpdateChatStatus(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter chat id:")
	chatIDStr := ""
	fmt.Scanln(&chatIDStr)
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Enter new chat status:")
	status := ""
	fmt.Scanln(&status)
	if err := serviceModule.ChatService.UpdateChat(chatID, status); err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Status updated")
}

func DeleteChat(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter chat id:")
	chatIDStr := ""
	fmt.Scanln(&chatIDStr)
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if err := serviceModule.ChatService.DeleteChat(chatID); err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Chat deleted")
}

func UpdateMessageContent(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter message id:")
	messageIDStr := ""
	fmt.Scanln(&messageIDStr)
	messageID, err := strconv.ParseInt(messageIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Enter new content:")
	scanner := bufio.NewScanner(os.Stdin)
	content := ""
	if scanner.Scan() {
		content = scanner.Text()
	}
	if err := serviceModule.ChatService.UpdateMessageContent(messageID, content); err != nil {
		fmt.Println("Error:", err)
		return
	}
	msg, err := serviceModule.ChatService.GetMessage(messageID)
	if err == nil {
		PrintMessage(msg)
	} else {
		fmt.Println("Content updated")
	}
}

func DeleteMessage(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter message id:")
	messageIDStr := ""
	fmt.Scanln(&messageIDStr)
	messageID, err := strconv.ParseInt(messageIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if err := serviceModule.ChatService.DeleteMessage(messageID); err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Message deleted")
}

func GetChatById(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter chat id:")
	chatIDStr := ""
	fmt.Scanln(&chatIDStr)
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	chat, err := serviceModule.ChatService.GetChat(chatID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Chat:")
	fmt.Println("ID:", chat.ID)
	fmt.Println("ClientID:", chat.ClientID)
	fmt.Println("RepetitorID:", chat.RepetitorID)
	fmt.Println("ModeratorID:", chat.ModeratorID)
	fmt.Println("Status:", chat.Status)
	fmt.Println("Type:", chat.Type)
	fmt.Println("CreatedAt:", chat.CreatedAt.Format(time.RFC3339))
}

func ChatWork(serviceModule *service_logic.ServiceModule) {
	for {
		fmt.Println("Chat work")
		fmt.Println("1. Create chat")
		fmt.Println("2. Get chat history")
		fmt.Println("3. Send message")
		fmt.Println("4. List user chats")
		fmt.Println("5. Get chat by ID")
		fmt.Println("6. Replace chat (PUT)")
		fmt.Println("7. Update chat status (PATCH)")
		fmt.Println("8. Delete chat (DELETE)")
		fmt.Println("9. Update message content (PATCH)")
		fmt.Println("10. Delete message (DELETE)")
		fmt.Println("11. Exit")
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
			ListUserChats(serviceModule)
		case 5:
			GetChatById(serviceModule)
		case 6:
			ReplaceChat(serviceModule)
		case 7:
			UpdateChatStatus(serviceModule)
		case 8:
			DeleteChat(serviceModule)
		case 9:
			UpdateMessageContent(serviceModule)
		case 10:
			DeleteMessage(serviceModule)
		case 11:
			return
		}
	}
}
