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

func PrintContractTypeMenu() {
	fmt.Println("Enter contract type:")
	fmt.Println("1. Tutoring")
	fmt.Println("2. Translation")
	fmt.Println("3. Writing")
	fmt.Println("4. Design")
	fmt.Println("5. Programming")
}

func CreateContract(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Note: you can't create contract without client already registered")
	contract := types.ServiceContractInitData{}
	fmt.Println("Enter contract client ID:")
	clientIDStr := ""
	fmt.Scanln(&clientIDStr)
	clientID, err := strconv.ParseInt(clientIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	contract.ClientID = clientID
	PrintContractTypeMenu()
	categoryStr := ""
	fmt.Scanln(&categoryStr)
	category, err := strconv.Atoi(categoryStr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if category < 1 || category > 5 {
		fmt.Println("Invalid category")
		return
	}
	contract.ContractCategory = types.ContractCategory(category)
	fmt.Println("Enter contract description:")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input := scanner.Text()
		contract.Description = input
	}
	fmt.Println("Enter contract price:")
	priceStr := ""
	fmt.Scanln(&priceStr)
	price, err := strconv.ParseInt(priceStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if price < 0 {
		fmt.Println("Invalid price")
		return
	}
	contract.Price = price
	fmt.Println("Enter contract commission:")
	commissionStr := ""
	fmt.Scanln(&commissionStr)
	commission, err := strconv.ParseInt(commissionStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if commission < 0 {
		fmt.Println("Invalid commission")
		return
	}
	contract.Commission = commission
	fmt.Println("Enter contract duration(in days):")
	durationStr := ""
	fmt.Scanln(&durationStr)
	duration, err := strconv.ParseInt(durationStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if duration < 0 {
		fmt.Println("Invalid duration")
		return
	}
	contract.StartDate = time.Now()
	contract.Duration = duration
	id, err := serviceModule.ContractService.CreateContract(contract)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Contract created successfully with ID:", id)
}

func AssignRepetitorToContract(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter contract ID:")
	contractIDStr := ""
	fmt.Scanln(&contractIDStr)
	contractID, err := strconv.ParseInt(contractIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Enter repetitor ID:")
	repetitorIDStr := ""
	fmt.Scanln(&repetitorIDStr)
	repetitorID, err := strconv.ParseInt(repetitorIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	err = serviceModule.ContractService.UpdateContractRepetitorID(contractID, repetitorID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Repetitor assigned to contract successfully")
}

func UpdateContractStatus(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter contract ID:")
	contractIDStr := ""
	fmt.Scanln(&contractIDStr)
	contractID, err := strconv.ParseInt(contractIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Enter new status:")
	fmt.Println("1. Pending")
	fmt.Println("2. Active")
	fmt.Println("3. Completed")
	fmt.Println("4. Cancelled")
	fmt.Println("5. Banned")
	statusStr := ""
	fmt.Scanln(&statusStr)
	status, err := strconv.Atoi(statusStr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if status < int(types.ContractStatusPending) || status > int(types.ContractStatusBanned) {
		fmt.Println("Invalid status")
		return
	}
	err = serviceModule.ContractService.UpdateContractStatus(contractID, types.ContractStatus(status))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Contract status updated successfully")
}

func PrintContract(contract *types.ServiceContract) {
	fmt.Println("Contract info:")
	fmt.Println("ID:", contract.ID)
	fmt.Println("Client ID:", contract.ClientID)
	fmt.Println("Repetitor ID:", contract.RepetitorID)
	switch contract.Status {
	case types.ContractStatusPending:
		fmt.Println("Status: Pending")
	case types.ContractStatusActive:
		fmt.Println("Status: Active")
	case types.ContractStatusCompleted:
		fmt.Println("Status: Completed")
	case types.ContractStatusCancelled:
		fmt.Println("Status: Cancelled")
	case types.ContractStatusBanned:
		fmt.Println("Status: Banned")
	}
	switch contract.PaymentStatus {
	case types.PaymentStatusNull:
		fmt.Println("Payment status: Null")
	case types.PaymentStatusPending:
		fmt.Println("Payment status: Pending")
	case types.PaymentStatusPaid:
		fmt.Println("Payment status: Paid")
	case types.PaymentStatusRefused:
		fmt.Println("Payment status: Refused")
	case types.PaymentStatusRefunded:
		fmt.Println("Payment status: Refunded")
	}
	fmt.Println("Client review ID:", contract.ReviewClientID)
	fmt.Println("Repetitor review ID:", contract.ReviewRepetitorID)
	fmt.Println("Price:", contract.Price)
	fmt.Println("Commission:", contract.Commission)
	fmt.Println("Description:", contract.Description)
	fmt.Println("Start date:", contract.StartDate)
	fmt.Println("End date:", contract.EndDate)
}

func GetContractInfo(serviceModule *service_logic.ServiceModule) {
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
	PrintContract(contract)
}

func GetLessonById(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter lesson ID:")
	lessonIDStr := ""
	fmt.Scanln(&lessonIDStr)
	lessonID, err := strconv.ParseInt(lessonIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	lesson, err := serviceModule.LessonService.GetLesson(lessonID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	PrintLesson(lesson)
}

func CreateClientReview(serviceModule *service_logic.ServiceModule) {
	var review types.ServiceReview
	fmt.Println("Enter review text:")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input := scanner.Text()
		review.Comment = input
	}
	fmt.Println("Enter review rating:")
	reviewRatingStr := ""
	fmt.Scanln(&reviewRatingStr)
	reviewRating, err := strconv.Atoi(reviewRatingStr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	review.Rating = reviewRating
	review.CreatedAt = time.Now()
	fmt.Println("Enter client ID:")
	clientIDStr := ""
	fmt.Scanln(&clientIDStr)
	clientID, err := strconv.ParseInt(clientIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	review.ClientID = clientID
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
	if contract.RepetitorID == 0 {
		fmt.Println("Contract has no repetitor")
		return
	}
	review.RepetitorID = contract.RepetitorID
	_, err = serviceModule.ContractService.CreateContractReviewClient(contractID, review)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Review created successfully")
}

func CreateRepetitorReview(serviceModule *service_logic.ServiceModule) {
	var review types.ServiceReview
	fmt.Println("Enter review text:")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input := scanner.Text()
		review.Comment = input
	}
	fmt.Println("Enter review rating:")
	reviewRatingStr := ""
	fmt.Scanln(&reviewRatingStr)
	reviewRating, err := strconv.Atoi(reviewRatingStr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	review.Rating = reviewRating
	review.CreatedAt = time.Now()
	fmt.Println("Enter repetitor ID:")
	repetitorIDStr := ""
	fmt.Scanln(&repetitorIDStr)
	repetitorID, err := strconv.ParseInt(repetitorIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	review.ClientID = 0
	review.RepetitorID = repetitorID
	fmt.Println("Enter contract ID:")
	contractIDStr := ""
	fmt.Scanln(&contractIDStr)
	contractID, err := strconv.ParseInt(contractIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	_, err = serviceModule.ContractService.CreateContractReviewRepetitor(contractID, review)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Review created successfully")
}

func WriteReview(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter review type:")
	fmt.Println("1. Client review")
	fmt.Println("2. Repetitor review")
	reviewTypeStr := ""
	fmt.Scanln(&reviewTypeStr)
	reviewType, err := strconv.Atoi(reviewTypeStr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if reviewType < 1 || reviewType > 2 {
		fmt.Println("Invalid review type")
		return
	}
	switch reviewType {
	case 1:
		CreateClientReview(serviceModule)
	case 2:
		CreateRepetitorReview(serviceModule)
	}
}

func UpdateClientReview(serviceModule *service_logic.ServiceModule) {
	var review types.ServiceReview
	fmt.Println("Enter review text:")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input := scanner.Text()
		review.Comment = input
	}
	fmt.Println("Enter review rating:")
	reviewRatingStr := ""
	fmt.Scanln(&reviewRatingStr)
	reviewRating, err := strconv.Atoi(reviewRatingStr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	review.Rating = reviewRating
	review.CreatedAt = time.Now()
	fmt.Println("Enter client ID:")
	clientIDStr := ""
	fmt.Scanln(&clientIDStr)
	clientID, err := strconv.ParseInt(clientIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	review.ClientID = clientID
	review.RepetitorID = 0
	fmt.Println("Enter contract ID:")
	contractIDStr := ""
	fmt.Scanln(&contractIDStr)
	contractID, err := strconv.ParseInt(contractIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	err = serviceModule.ContractService.UpdateContractReviewClient(contractID, review)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Review created successfully")
}

func UpdateRepetitorReview(serviceModule *service_logic.ServiceModule) {
	var review types.ServiceReview
	fmt.Println("Enter review text:")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input := scanner.Text()
		review.Comment = input
	}
	fmt.Println("Enter review rating:")
	reviewRatingStr := ""
	fmt.Scanln(&reviewRatingStr)
	reviewRating, err := strconv.Atoi(reviewRatingStr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	review.Rating = reviewRating
	review.CreatedAt = time.Now()
	fmt.Println("Enter repetitor ID:")
	repetitorIDStr := ""
	fmt.Scanln(&repetitorIDStr)
	repetitorID, err := strconv.ParseInt(repetitorIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	review.ClientID = 0
	review.RepetitorID = repetitorID
	fmt.Println("Enter contract ID:")
	contractIDStr := ""
	fmt.Scanln(&contractIDStr)
	contractID, err := strconv.ParseInt(contractIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	err = serviceModule.ContractService.UpdateContractReviewRepetitor(contractID, review)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Review created successfully")
}

func UpdateReview(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter review type:")
	fmt.Println("1. Client review")
	fmt.Println("2. Repetitor review")
	reviewTypeStr := ""
	fmt.Scanln(&reviewTypeStr)
	reviewType, err := strconv.Atoi(reviewTypeStr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if reviewType < 1 || reviewType > 2 {
		fmt.Println("Invalid review type")
		return
	}
	switch reviewType {
	case 1:
		UpdateClientReview(serviceModule)
	case 2:
		UpdateRepetitorReview(serviceModule)
	}
}

func GetReviews(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter review type:")
	fmt.Println("1. Client review")
	fmt.Println("2. Repetitor review")
	reviewTypeStr := ""
	fmt.Scanln(&reviewTypeStr)
	reviewType, err := strconv.Atoi(reviewTypeStr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if reviewType < 1 || reviewType > 2 {
		fmt.Println("Invalid review type")
		return
	}
	switch reviewType {
	case 1:
		GetReviewsByClientID(serviceModule)
	case 2:
		GetReviewsByRepetitorID(serviceModule)
	}
}

func PrintReview(review *types.ServiceReview) {
	fmt.Println("Review info:")
	fmt.Println("Client ID:", review.ClientID)
	fmt.Println("Repetitor ID:", review.RepetitorID)
	fmt.Println("Rating:", review.Rating)
	fmt.Println("Comment:", review.Comment)
	fmt.Println("Created at:", review.CreatedAt)
}

func GetReviewsByRepetitorID(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter repetitor ID:")
	repetitorIDStr := ""
	fmt.Scanln(&repetitorIDStr)
	repetitorID, err := strconv.ParseInt(repetitorIDStr, 10, 64)
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
	reviews, err := serviceModule.ReviewService.GetReviewsByRepetitorID(repetitorID, offset, limit)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, review := range reviews {
		PrintReview(&review)
	}
}

func GetReviewsByClientID(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter client ID:")
	clientIDStr := ""
	fmt.Scanln(&clientIDStr)
	clientID, err := strconv.ParseInt(clientIDStr, 10, 64)
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
	reviews, err := serviceModule.ReviewService.GetReviewsByClientID(clientID, offset, limit)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, review := range reviews {
		PrintReview(&review)
	}
}

func PrintLesson(lesson *types.ServiceLesson) {
	fmt.Println("Lesson info:")
	fmt.Println("Contract ID:", lesson.ContractID)
	fmt.Println("Duration:", lesson.Duration)
	fmt.Println("Created at:", lesson.CreatedAt)
}

func AddLesson(serviceModule *service_logic.ServiceModule) {
	lesson := types.ServiceLesson{}
	fmt.Println("Enter contract ID:")
	contractIDStr := ""
	fmt.Scanln(&contractIDStr)
	contractID, err := strconv.ParseInt(contractIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Enter lesson duration:")
	durationStr := ""
	fmt.Scanln(&durationStr)
	duration, err := strconv.ParseInt(durationStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	lesson.ContractID = contractID
	lesson.Duration = duration
	lesson.CreatedAt = time.Now()
	id, err := serviceModule.LessonService.CreateLesson(lesson)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Lesson added successfully with ID:", id)
}

func GetLessons(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Enter contract ID:")
	contractIDStr := ""
	fmt.Scanln(&contractIDStr)
	contractID, err := strconv.ParseInt(contractIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	lessons, err := serviceModule.LessonService.GetLessons(contractID, 0, 100)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, lesson := range lessons {
		PrintLesson(&lesson)
	}
}

func ListContracts(serviceModule *service_logic.ServiceModule) {
	fmt.Println("List contracts (V2)")
	fmt.Println("Enter client_id (0 for any):")
	clientIDStr := ""
	fmt.Scanln(&clientIDStr)
	clientID, err := strconv.ParseInt(clientIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Enter repetitor_id (0 for any):")
	repetitorIDStr := ""
	fmt.Scanln(&repetitorIDStr)
	repetitorID, err := strconv.ParseInt(repetitorIDStr, 10, 64)
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
	contracts, err := serviceModule.ContractService.GetContracts(clientID, repetitorID, offset, limit)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if len(contracts) == 0 {
		fmt.Println("No contracts found")
		return
	}
	for _, c := range contracts {
		fmt.Println("--- Contract ---")
		fmt.Println("ID:", c.ID)
		fmt.Println("Client ID:", c.ClientID)
		if c.RepetitorID != nil {
			fmt.Println("Repetitor ID:", *c.RepetitorID)
		} else {
			fmt.Println("Repetitor ID: <none>")
		}
		fmt.Println("Description:", c.Description)
		fmt.Println("Rate:", c.Rate)
		fmt.Println("Format:", c.Format)
		fmt.Println("Status:", c.Status)
		fmt.Println("CreatedAt:", c.CreatedAt)
	}
}

func CreateContractV2(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Create contract (V2)")
	var req types.ServerContractCreateV2
	fmt.Println("Enter client_id:")
	clientIDStr := ""
	fmt.Scanln(&clientIDStr)
	clientID, err := strconv.ParseInt(clientIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	req.ClientID = clientID
	fmt.Println("Enter description:")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		req.Description = scanner.Text()
	}
	fmt.Println("Enter rate:")
	rateStr := ""
	fmt.Scanln(&rateStr)
	rate, err := strconv.ParseInt(rateStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	req.Rate = rate
	req.Format = "online"
	init := types.MapperContractCreateV2ServerToServiceInit(&req)
	id, err := serviceModule.ContractService.CreateContract(*init)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	contract, err := serviceModule.ContractService.GetContract(id)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	sc := types.MapperContractServiceToServerV2(contract)
	fmt.Println("Contract created:")
	fmt.Println("ID:", sc.ID)
	fmt.Println("Client ID:", sc.ClientID)
	if sc.RepetitorID != nil {
		fmt.Println("Repetitor ID:", *sc.RepetitorID)
	} else {
		fmt.Println("Repetitor ID: <none>")
	}
	fmt.Println("Description:", sc.Description)
	fmt.Println("Rate:", sc.Rate)
	fmt.Println("Format:", sc.Format)
	fmt.Println("Status:", sc.Status)
}

func PatchLesson(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Patch lesson is not supported in current server API")
}

func DeleteLesson(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Delete lesson is not supported in current server API")
}

func ListContractTransactions(serviceModule *service_logic.ServiceModule) {
	fmt.Println("List contract transactions (V2)")
	fmt.Println("Enter contract ID:")
	contractIDStr := ""
	fmt.Scanln(&contractIDStr)
	if _, err := strconv.ParseInt(contractIDStr, 10, 64); err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Enter offset:")
	offsetStr := ""
	fmt.Scanln(&offsetStr)
	if _, err := strconv.ParseInt(offsetStr, 10, 64); err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Enter limit:")
	limitStr := ""
	fmt.Scanln(&limitStr)
	if _, err := strconv.ParseInt(limitStr, 10, 64); err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Transactions: [] (not implemented server-side)")
}

func CreateContractTransaction(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Create contract transaction (V2)")
	fmt.Println("Enter contract ID:")
	contractIDStr := ""
	fmt.Scanln(&contractIDStr)
	if _, err := strconv.ParseInt(contractIDStr, 10, 64); err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Enter amount:")
	amountStr := ""
	fmt.Scanln(&amountStr)
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	tx := types.ServerTransactionV2{ID: 0, ContractID: 0, Amount: amount, Status: types.TransactionStatusPending.String(), CreatedAt: time.Now()}
	fmt.Println("Created:")
	fmt.Println("Amount:", tx.Amount)
	fmt.Println("Status:", tx.Status)
	fmt.Println("CreatedAt:", tx.CreatedAt)
}

func ApproveTransaction(serviceModule *service_logic.ServiceModule) {
	fmt.Println("Approve transaction (V2)")
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
	serverTx := types.MapperTransactionServiceToServerV2(tx)
	fmt.Println("Approved:")
	fmt.Println("ID:", serverTx.ID)
	fmt.Println("Amount:", serverTx.Amount)
	fmt.Println("Status:", serverTx.Status)
	fmt.Println("CreatedAt:", serverTx.CreatedAt)
}

func ListContractReviews(serviceModule *service_logic.ServiceModule) {
	fmt.Println("List contract reviews (V2)")
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
	reviews := make([]types.ServiceReview, 0)
	if contract.ReviewClientID != 0 {
		if r, err := serviceModule.ReviewService.GetReview(contract.ReviewClientID); err == nil {
			reviews = append(reviews, *r)
		}
	}
	if contract.ReviewRepetitorID != 0 {
		if r, err := serviceModule.ReviewService.GetReview(contract.ReviewRepetitorID); err == nil {
			reviews = append(reviews, *r)
		}
	}
	if len(reviews) == 0 {
		fmt.Println("No reviews")
		return
	}
	for _, r := range reviews {
		fmt.Println("--- Review ---")
		fmt.Println("Client ID:", r.ClientID)
		fmt.Println("Repetitor ID:", r.RepetitorID)
		fmt.Println("Score:", r.Rating)
		fmt.Println("Text:", r.Comment)
		fmt.Println("CreatedAt:", r.CreatedAt)
	}
}

func ContractWork(serviceModule *service_logic.ServiceModule) {
	for {
		fmt.Println("Contract work")
		fmt.Println("1. Create contract")
		fmt.Println("2. Assign repetitor to contract")
		fmt.Println("3. Update contract status")
		fmt.Println("4. Get contract info")
		fmt.Println("5. Write review")
		fmt.Println("6. Get reviews")
		fmt.Println("7. Add lesson")
		fmt.Println("8. Get lessons")
		fmt.Println("9. Get lesson by ID ")
		fmt.Println("10. List contract reviews ")
		fmt.Println("11. List contracts")
		fmt.Println("12. Create contract via V2 ")
		fmt.Println("13. Update lesson")
		fmt.Println("14. Delete lesson")
		fmt.Println("15. List contract transactions ")
		fmt.Println("16. Create contract transaction ")
		fmt.Println("17. Approve transaction")
		fmt.Println("18. Exit")
		choiceStr := ""
		fmt.Scanln(&choiceStr)
		choice, err := strconv.Atoi(choiceStr)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if choice < 1 || choice > 18 {
			fmt.Println("Invalid choice")
			continue
		}
		switch choice {
		case 1:
			CreateContract(serviceModule)
		case 2:
			AssignRepetitorToContract(serviceModule)
		case 3:
			UpdateContractStatus(serviceModule)
		case 4:
			GetContractInfo(serviceModule)
		case 5:
			WriteReview(serviceModule)
		case 6:
			GetReviews(serviceModule)
		case 7:
			AddLesson(serviceModule)
		case 8:
			GetLessons(serviceModule)
		case 9:
			GetLessonById(serviceModule)
		case 10:
			ListContractReviews(serviceModule)
		case 11:
			ListContracts(serviceModule)
		case 12:
			CreateContractV2(serviceModule)
		case 13:
			PatchLesson(serviceModule)
		case 14:
			DeleteLesson(serviceModule)
		case 15:
			ListContractTransactions(serviceModule)
		case 16:
			CreateContractTransaction(serviceModule)
		case 17:
			ApproveTransaction(serviceModule)
		case 18:
			return
		default:
			fmt.Println("Invalid choice")
		}
	}
}
