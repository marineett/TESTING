package server

import (
	"data_base_project/service_logic"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func SetupServer(service_module *service_logic.ServiceModule, port string, logger *log.Logger) *http.Server {
	router := mux.NewRouter()
	router.StrictSlash(true)
	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	router.HandleFunc(AUTH_LOGIN_V2, AuthorizeHandlerV2(service_module.AuthService)).Methods("POST")

	router.Handle(strings.TrimSuffix(REGISTRATION_API_V2, "/"), RegistrationHandlerV2(
		service_module.ClientService,
		service_module.ModeratorService,
		service_module.AdminService,
		service_module.RepetitorService,
		service_module.AuthService,
	)).Methods("POST")

	router.Handle(CONTRACTS_V2, JWTAuthMiddleware(ContractsListHandlerV2(service_module.ContractService))).Methods("GET")
	router.Handle(CONTRACTS_V2, JWTAuthMiddleware(ContractCreateHandlerV2(service_module.ContractService))).Methods("POST")
	router.Handle(EXACT_CONTRACT_V2, JWTAuthMiddleware(ContractGetHandlerV2(service_module.ContractService))).Methods("GET")
	router.Handle(EXACT_CONTRACT_V2, JWTAuthMiddleware(ContractStatusPatchHandlerV2(service_module.ContractService))).Methods("PATCH")
	router.Handle(CONTRACT_LESSONS_V2, JWTAuthMiddleware(ContractLessonsListHandlerV2(service_module.LessonService))).Methods("GET")
	router.Handle(CONTRACT_LESSONS_V2, JWTAuthMiddleware(ContractLessonCreateHandlerV2(service_module.LessonService))).Methods("POST")
	router.Handle(CONTRACT_REVIEWS_V2, JWTAuthMiddleware(ContractReviewsListHandlerV2(service_module.ReviewService, service_module.ContractService))).Methods("GET")
	router.Handle(CONTRACT_REVIEWS_V2, JWTAuthMiddleware(ContractReviewCreateHandlerV2(service_module.ReviewService, service_module.ContractService))).Methods("POST")
	router.Handle(CONTRACT_TRANSACTIONS_V2, JWTAuthMiddleware(ContractTransactionsListHandlerV2(service_module.TransactionService))).Methods("GET")
	router.Handle(CONTRACT_TRANSACTIONS_V2, JWTAuthMiddleware(ContractTransactionCreateHandlerV2(service_module.TransactionService))).Methods("POST")
	router.Handle(TRANSACTION_APPROVAL_V2, JWTAuthMiddleware(TransactionApproveHandlerV2(service_module.TransactionService))).Methods("PATCH")

	router.Handle(EXACT_CLIENT_V2, JWTAuthMiddleware(ClientGetHandlerV2(service_module.ClientService))).Methods("GET")
	router.Handle(EXACT_REPETITOR_V2, JWTAuthMiddleware(RepetitorGetHandlerV2(service_module.RepetitorService))).Methods("GET")

	router.Handle(CHATS_V2, JWTAuthMiddleware(ChatGetChatsHandlerV2(service_module.ChatService))).Methods("GET")
	router.Handle(CHATS_V2, JWTAuthMiddleware(ChatCreateChatHandlerV2(service_module.ChatService))).Methods("POST")
	router.Handle(EXACT_CHAT_V2, JWTAuthMiddleware(ChatGetChatHandlerV2(service_module.ChatService))).Methods("GET")
	router.Handle(EXACT_CHAT_V2, JWTAuthMiddleware(ChatUpdateChatHandlerV2(service_module.ChatService))).Methods("PATCH")
	router.Handle(EXACT_CHAT_V2, JWTAuthMiddleware(ChatDeleteChatHandlerV2(service_module.ChatService))).Methods("DELETE")
	router.Handle(EXACT_CHAT_V2, JWTAuthMiddleware(ChatClearChatHandlerV2(service_module.ChatService))).Methods("PUT")
	router.Handle(EXACT_CHAT_MESSAGES_V2, JWTAuthMiddleware(ChatGetMessagesHandlerV2(service_module.ChatService))).Methods("GET")
	router.Handle(EXACT_CHAT_MESSAGES_V2, JWTAuthMiddleware(ChatSendMessageHandlerV2(service_module.ChatService))).Methods("POST")
	router.Handle(EXACT_MESSAGE_V2, JWTAuthMiddleware(UpdateMessageContentHandlerV2(service_module.ChatService))).Methods("PATCH")
	router.Handle(EXACT_MESSAGE_V2, JWTAuthMiddleware(DeleteMessageHandlerV2(service_module.ChatService))).Methods("DELETE")

	router.Handle(EXACT_ADMIN_V2, JWTAuthMiddleware(AdminGetProfileHandlerV2(service_module.AdminService))).Methods("GET")
	router.Handle(DEPARTMENTS_V2, JWTAuthMiddleware(AdminCreateDepartmentHandlerV2(service_module.DepartmentService))).Methods("POST")
	router.Handle(ADMIN_DEPARTMENTS_V2, JWTAuthMiddleware(AdminListDepartmentsHandlerV2(service_module.DepartmentService))).Methods("GET")
	router.Handle(EXACT_DEPARTMENT_V2, JWTAuthMiddleware(DepartmentReplaceHandlerV2(service_module.DepartmentService))).Methods("PUT")
	router.Handle(EXACT_DEPARTMENT_V2, JWTAuthMiddleware(DepartmentDeleteHandlerV2(service_module.DepartmentService))).Methods("DELETE")
	router.Handle(DEPARTMENT_MODERATOR_V2, JWTAuthMiddleware(DepartmentAssignModeratorHandlerV2(service_module.DepartmentService))).Methods("PUT")
	router.Handle(DEPARTMENT_MODERATOR_V2, JWTAuthMiddleware(DepartmentRemoveModeratorHandlerV2(service_module.DepartmentService))).Methods("DELETE")
	router.Handle(MODERATORS_V2, JWTAuthMiddleware(ModeratorsListHandlerV2(service_module.ModeratorService))).Methods("GET")
	router.Handle(MODERATOR_SALARY_V2, JWTAuthMiddleware(ModeratorSalaryPatchHandlerV2(service_module.ModeratorService))).Methods("PATCH")

	router.Handle(REGISTRATION_API, SetupRegistrationRouter(
		service_module.AuthService,
		service_module.ModeratorService,
		service_module.ClientService,
		service_module.AdminService,
		service_module.RepetitorService,
		logger,
	))
	router.Handle(AUTH_API, SetupAuthorizeRouterV1(service_module.AuthService, logger))
	router.Handle(CONTRACT_API, SetupContractRouter(
		service_module.ContractService,
		service_module.ReviewService,
		service_module.LessonService,
		logger,
	))
	router.Handle(CLIENT_API, SetupClientRouter(
		service_module.ClientService,
		service_module.ContractService,
		logger,
	))
	router.Handle(REPETITOR_API, SetupRepetitorRouter(
		service_module.RepetitorService,
		service_module.ContractService,
		service_module.TransactionService,
		service_module.ResumeService,
		logger,
	))
	router.Handle(MODERATOR_API, SetupModeratorRouter(
		service_module.TransactionService,
		service_module.ContractService,
		service_module.ModeratorService,
		logger,
	))
	router.Handle(ADMIN_API, SetupAdminRouter(
		service_module.AdminService,
		service_module.DepartmentService,
		service_module.ModeratorService,
		logger,
	))
	router.Handle(CHAT_API, SetupChatRouter(
		service_module.ChatService,
		logger,
	))
	router.Handle(GUEST_API, SetupGuestRouter(
		service_module.RepetitorService,
		logger,
	))

	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	addr := ":" + port
	logger.Printf("Server starting on addr %s", addr)
	return &http.Server{
		Addr:    addr,
		Handler: CORSMiddleware(router),
	}
}
