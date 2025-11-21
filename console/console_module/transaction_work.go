package console_module

import (
	"data_base_project/service_logic"
	"data_base_project/types"
	"fmt"
	"strconv"
)

func PayForContract(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter contract ID:")
	contractIDStr := ""
	fmt.Scanln(&contractIDStr)
	contractID, err := strconv.ParseInt(contractIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	contract, err := serviceModule.ContractService.GetContract(contractID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if contract.PaymentStatus != types.PaymentStatusNull {
		fmt.Println("Contract payment status is not null")
		return
	}
	fmt.Println("Enter transaction amount (note: contract price is", contract.Price, "):")
	amountStr := ""
	fmt.Scanln(&amountStr)
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if amount != contract.Price {
		fmt.Println("Invalid amount")
		return
	}
	transactionID, err := serviceModule.TransactionService.CreateContractPaymentTransaction(amount, contract.RepetitorID, contract.ID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Transaction created successfully with ID:", transactionID)
}

func GetTransactionInfo(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter transaction ID:")
	transactionIDStr := ""
	fmt.Scanln(&transactionIDStr)
	transactionID, err := strconv.ParseInt(transactionIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	transaction, err := serviceModule.TransactionService.GetTransaction(transactionID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Transaction info:")
	fmt.Println("ID:", transaction.ID)
	fmt.Println("Amount:", transaction.Amount)
	fmt.Println("Status:", transaction.Status)
	fmt.Println("Created at:", transaction.CreatedAt)
}

func PrintTransaction(transaction *types.ServiceTransaction) {
	fmt.Println("Transaction info:")
	fmt.Println("ID:", transaction.ID)
	fmt.Println("Amount:", transaction.Amount)
	fmt.Println("Status:", transaction.Status)
	fmt.Println("Created at:", transaction.CreatedAt)
}

func GetTransactionList(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter user ID:")
	userIDStr := ""
	fmt.Scanln(&userIDStr)
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	transactions, err := serviceModule.TransactionService.GetTransactionsList(userID, 0, 100)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Transactions:")
	for _, transaction := range transactions {
		PrintTransaction(&transaction)
	}
}

func TransactionWork(serviceModule *service_logic.ServiceModule) {
	for {
		fmt.Println("Transaction work")
		fmt.Println("1. Pay for contract")
		fmt.Println("2. Get transaction info")
		fmt.Println("3. Get transaction list")
		fmt.Println("4. Approve transaction")
		fmt.Println("5. Exit")
		choiceStr := ""
		fmt.Scanln(&choiceStr)
		choice, err := strconv.Atoi(choiceStr)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if choice < 1 || choice > 5 {
			fmt.Println("Invalid choice")
			continue
		}
		switch choice {
		case 1:
			PayForContract(serviceModule)
		case 2:
			GetTransactionInfo(serviceModule)
		case 3:
			GetTransactionList(serviceModule)
		case 4:
			ApproveTransactionGlobal(serviceModule)
		case 5:
			return
		default:
			fmt.Println("Invalid choice")
		}
	}
}

func ApproveTransactionGlobal(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter transaction ID:")
	txIDStr := ""
	fmt.Scanln(&txIDStr)
	txID, err := strconv.ParseInt(txIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if err := serviceModule.TransactionService.ApproveTransaction(txID); err != nil {
		fmt.Println("Error:", err)
		return
	}
	tx, err := serviceModule.TransactionService.GetTransaction(txID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	PrintTransaction(tx)
}
