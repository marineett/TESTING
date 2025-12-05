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
	router.StrictSlash(false)

	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	router.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("# Prometheus metrics endpoint\n# Metrics are exported via OpenTelemetry\n"))
	})

	router.Handle(API_V2, ApiV2Handler()).Methods("GET")
	router.Handle(strings.TrimSuffix(API_V2, "/"), SwaggerSpecHandler()).Methods("GET")
	router.Handle(OPENAPI_YAML_V2, SwaggerSpecHandler()).Methods("GET")
	router.Handle(DOCUMENTATION_V2, DocumentationHandler()).Methods("GET")
	router.Handle(STATIC_FILES_V2, StaticFileHandler()).Methods("GET")
	router.Handle(STATIC_FILES_V2+"/", StaticFileHandler()).Methods("GET")
	router.Handle(EXACT_STATIC_FILE_V2, StaticFileHandler()).Methods("GET")
	router.Handle(RESERVED_FILES_V2, ReservedStaticFileHandler()).Methods("GET")
	router.Handle(RESERVED_FILES_V2+"/", ReservedStaticFileHandler()).Methods("GET")
	router.Handle(RESERVED_FILES_V2+"/", ReservedStaticFileHandler()).Methods("GET")
	router.Handle(EXACT_RESERVED_FILE_V2, ReservedStaticFileHandler()).Methods("GET")
	router.Handle(LEGACY_ARCHIVE_V2, LegacyArchiveHandler()).Methods("GET")

	router.HandleFunc(AUTH_LOGIN_V2, AuthorizeHandlerV2(service_module.AuthService)).Methods("POST")
	router.HandleFunc(AUTH_LOGIN_V2+"/", AuthorizeHandlerV2(service_module.AuthService)).Methods("POST")

	registrationHandler := RegistrationHandlerV2(
		service_module.ClientService,
		service_module.ModeratorService,
		service_module.AdminService,
		service_module.RepetitorService,
		service_module.AuthService,
		service_module.EmailSender,
	)
	router.HandleFunc(REGISTRATION_API_V2, registrationHandler).Methods("POST")
	router.HandleFunc(REGISTRATION_API_V2+"/", registrationHandler).Methods("POST")

	applyTokenHandler := ApplyTokenHandler(service_module.AuthService)
	router.HandleFunc(APPLY_TOKEN_API_V2, applyTokenHandler).Methods("GET")
	router.HandleFunc(APPLY_TOKEN_API_V2+"/", applyTokenHandler).Methods("GET")
	router.HandleFunc(UPDATE_TOKEN_API_V2, UpdateTokenHandler(service_module.AuthService, service_module.EmailSender)).Methods("POST")
	router.Handle(CONTRACTS_V2, ContractsListHandlerV2(service_module.ContractService)).Methods("GET")
	router.Handle(CONTRACTS_V2, ContractCreateHandlerV2(service_module.ContractService)).Methods("POST")
	router.Handle(EXACT_CONTRACT_V2, ContractGetHandlerV2(service_module.ContractService)).Methods("GET")
	router.Handle(EXACT_CONTRACT_V2, ContractStatusPatchHandlerV2(service_module.ContractService)).Methods("PATCH")
	router.Handle(CONTRACT_LESSONS_V2, ContractLessonsListHandlerV2(service_module.LessonService)).Methods("GET")
	router.Handle(CONTRACT_LESSONS_V2, ContractLessonCreateHandlerV2(service_module.LessonService)).Methods("POST")
	router.Handle(EXACT_LESSON_V2, LessonGetHandlerV2(service_module.LessonService)).Methods("GET")
	router.Handle(EXACT_LESSON_V2, LessonPatchHandlerV2(service_module.LessonService)).Methods("PATCH")
	router.Handle(EXACT_LESSON_V2, LessonDeleteHandlerV2(service_module.LessonService)).Methods("DELETE")
	router.Handle(CONTRACT_REVIEWS_V2, ContractReviewsListHandlerV2(service_module.ReviewService, service_module.ContractService)).Methods("GET")
	router.Handle(CONTRACT_REVIEWS_V2, ContractReviewCreateHandlerV2(service_module.ReviewService, service_module.ContractService)).Methods("POST")
	router.Handle(CONTRACT_TRANSACTIONS_V2, ContractTransactionsListHandlerV2(service_module.TransactionService, service_module.ContractService)).Methods("GET")
	router.Handle(CONTRACT_TRANSACTIONS_V2, ContractTransactionCreateHandlerV2(service_module.TransactionService)).Methods("POST")
	router.Handle(TRANSACTION_APPROVAL_V2, TransactionApproveHandlerV2(service_module.TransactionService)).Methods("PATCH")

	router.Handle(EXACT_CLIENT_V2, ClientGetHandlerV2(service_module.ClientService)).Methods("GET")
	router.Handle(EXACT_REPETITOR_V2, RepetitorGetHandlerV2(service_module.RepetitorService)).Methods("GET")
	router.Handle(EXACT_REPETITOR_V2, RepetitorAssignContractHandlerV2(service_module.ContractService)).Methods("PATCH")

	router.Handle(CHATS_V2, ChatGetChatsHandlerV2(service_module.ChatService)).Methods("GET")
	router.Handle(CHATS_V2, ChatCreateChatHandlerV2(service_module.ChatService)).Methods("POST")
	router.Handle(EXACT_CHAT_V2, ChatGetChatHandlerV2(service_module.ChatService)).Methods("GET")
	router.Handle(EXACT_CHAT_V2, ChatUpdateChatHandlerV2(service_module.ChatService)).Methods("PATCH")
	router.Handle(EXACT_CHAT_V2, ChatDeleteChatHandlerV2(service_module.ChatService)).Methods("DELETE")
	router.Handle(EXACT_CHAT_V2, ChatClearChatHandlerV2(service_module.ChatService)).Methods("PUT")
	router.Handle(EXACT_CHAT_MESSAGES_V2, ChatGetMessagesHandlerV2(service_module.ChatService)).Methods("GET")
	router.Handle(EXACT_CHAT_MESSAGES_V2, ChatSendMessageHandlerV2(service_module.ChatService)).Methods("POST")
	router.Handle(EXACT_MESSAGE_V2, UpdateMessageContentHandlerV2(service_module.ChatService)).Methods("PATCH")
	router.Handle(EXACT_MESSAGE_V2, DeleteMessageHandlerV2(service_module.ChatService)).Methods("DELETE")

	router.Handle(EXACT_ADMIN_V2, AdminGetProfileHandlerV2(service_module.AdminService)).Methods("GET")
	router.Handle(DEPARTMENTS_V2, AdminCreateDepartmentHandlerV2(service_module.DepartmentService)).Methods("POST")
	router.Handle(ADMIN_DEPARTMENTS_V2, AdminListDepartmentsHandlerV2(service_module.DepartmentService)).Methods("GET")
	router.Handle(EXACT_DEPARTMENT_V2, DepartmentReplaceHandlerV2(service_module.DepartmentService)).Methods("PUT")
	router.Handle(EXACT_DEPARTMENT_V2, DepartmentDeleteHandlerV2(service_module.DepartmentService)).Methods("DELETE")
	router.Handle(DEPARTMENT_MODERATOR_V2, DepartmentAssignModeratorHandlerV2(service_module.DepartmentService)).Methods("PUT")
	router.Handle(DEPARTMENT_MODERATOR_V2, DepartmentRemoveModeratorHandlerV2(service_module.DepartmentService)).Methods("DELETE")
	router.Handle(MODERATORS_V2, ModeratorsListHandlerV2(service_module.ModeratorService)).Methods("GET")
	router.Handle(MODERATOR_SALARY_V2, ModeratorSalaryPatchHandlerV2(service_module.ModeratorService)).Methods("PATCH")

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

	handler := CORSMiddleware(router)

	return &http.Server{
		Addr:    addr,
		Handler: handler,
	}
}
