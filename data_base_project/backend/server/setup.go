package server

import (
	"data_base_project/service_logic"
	"log"
	"net/http"
)

func SetupServer(service_module *service_logic.ServiceModule, port string, logger *log.Logger) *http.Server {
	router := http.NewServeMux()

	router.Handle(REGISTRATION_API, SetupRegistrationRouter(
		service_module.AuthService,
		service_module.ModeratorService,
		service_module.ClientService,
		service_module.AdminService,
		service_module.RepetitorService,
		logger,
	))
	router.Handle(AUTH_API, SetupAuthorizeRouter(service_module.AuthService, logger))
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
	logger.Printf("Server starting on port %s", port)
	return &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
}
