package server

const (
	API_V1 = "/api/v1/"
	API_V2 = "/api/v2/"
)

const (
	USER_API         = API_V1 + "user/"
	GUEST_API        = API_V1 + "guest/"
	ADMIN_API        = API_V1 + "admin/"
	MODERATOR_API    = API_V1 + "moderator/"
	CLIENT_API       = API_V1 + "client/"
	REGISTRATION_API = API_V1 + "registration/"
	AUTH_API         = API_V1 + "auth/"
	CONTRACT_API     = API_V1 + "contract/"
	CHAT_API         = API_V1 + "chat/"
	REPETITOR_API    = API_V1 + "repetitor/"
)

const (
	REGISTRATION_CLIENT    = REGISTRATION_API + "client"
	REGISTRATION_REPETITOR = REGISTRATION_API + "repetitor"
	REGISTRATION_ADMIN     = REGISTRATION_API + "admin"
	REGISTRATION_MODERATOR = REGISTRATION_API + "moderator"
)

const (
	GUEST_GET_REPETITORS = GUEST_API + "get_repetitors"
)

const (
	AUTH_AUTHORIZE = AUTH_API + "authorize"
)

const (
	CLIENT_GET_PROFILE     = CLIENT_API + "get_profile"
	CLIENT_CREATE_CONTRACT = CLIENT_API + "create_contract"
	CLIENT_GET_CONTRACTS   = CLIENT_API + "get_contracts"
	CLIENT_MAKE_REVIEW     = CLIENT_API + "make_review"
)

const (
	CONTRACT_GET        = CONTRACT_API + "get"
	CONTRACT_GET_REVIEW = CONTRACT_API + "get_review"
	ADD_LESSON          = CONTRACT_API + "add_lesson"
	GET_LESSONS         = CONTRACT_API + "get_lessons"
)

const (
	REPETITOR_GET_PROFILE             = REPETITOR_API + "get_profile"
	REPETITOR_GET_CONTRACTS           = REPETITOR_API + "get_contracts"
	REPETITOR_GET_AVAILABLE_CONTRACTS = REPETITOR_API + "get_available_contracts"
	REPETITOR_ACCEPT_CONTRACT         = REPETITOR_API + "accept_contract"
	REPETITOR_MAKE_REVIEW             = REPETITOR_API + "make_review"
	REPETITOR_PAY_FOR_CONTRACT        = REPETITOR_API + "pay_for_contract"
	REPETITOR_CANCEL_CONTRACT         = REPETITOR_API + "cancel_contract"
	REPETITOR_COMPLETE_CONTRACT       = REPETITOR_API + "complete_contract"
	REPETITOR_CHANGE_RESUME           = REPETITOR_API + "change_resume"
)

const (
	MODERATOR_GET_PROFILE                = MODERATOR_API + "get_profile"
	MODERATOR_GET_TRANSACTION_TO_APPROVE = MODERATOR_API + "get_transaction_to_approve"
	MODERATOR_APPROVE_TRANSACTION        = MODERATOR_API + "approve_transaction"
	MODERATOR_GET_CONTRACTS              = MODERATOR_API + "get_contracts"
	MODERATOR_BAN_CONTRACT               = MODERATOR_API + "ban_contract"
)

const (
	ADMIN_GET_PROFILE             = ADMIN_API + "get_profile"
	ADMIN_CREATE_DEPARTMENT       = ADMIN_API + "create_department"
	ADMIN_GET_DEPARTMENTS         = ADMIN_API + "get_departments"
	ADMIN_GET_MODERATORS          = ADMIN_API + "get_moderators"
	ADMIN_HIRE_MODERATOR          = ADMIN_API + "hire_moderator"
	ADMIN_FIRE_MODERATOR          = ADMIN_API + "fire_moderator"
	ADMIN_CHANGE_MODERATOR_SALARY = ADMIN_API + "change_moderator_salary"
)

const (
	CHAT_GET_CLIENT_CHATS    = CHAT_API + "get_client_chats"
	CHAT_GET_REPETITOR_CHATS = CHAT_API + "get_repetitor_chats"
	CHAT_GET_MODERATOR_CHATS = CHAT_API + "get_moderator_chats"
	CHAT_START_CM_CHAT       = CHAT_API + "start_cm_chat"
	CHAT_START_RM_CHAT       = CHAT_API + "start_rm_chat"
	CHAT_START_CR_CHAT       = CHAT_API + "start_cr_chat"
	CHAT_GET_CHAT            = CHAT_API + "get_chat"
	CHAT_SEND_MESSAGE        = CHAT_API + "send_message"
	CHAT_GET_MESSAGES        = CHAT_API + "get_messages"
	CHAT_DELETE_CHAT         = CHAT_API + "delete_chat"
	CHAT_CLEAR_MESSAGES      = CHAT_API + "clear_messages"
)

const (
	USER_API_V2         = API_V2 + "user/"
	GUEST_API_V2        = API_V2 + "guest/"
	ADMIN_API_V2        = API_V2 + "admin/"
	MODERATOR_API_V2    = API_V2 + "moderator/"
	CLIENT_API_V2       = API_V2 + "client/"
	REGISTRATION_API_V2 = API_V2 + "registration/"
	AUTH_API_V2         = API_V2 + "auth/"
	CONTRACT_API_V2     = API_V2 + "contract/"
	CHAT_API_V2         = API_V2 + "chat/"
	REPETITOR_API_V2    = API_V2 + "repetitor/"
)

const (
	AUTH_REGISTRATION = AUTH_API_V2 + "registration"
	AUTH_LOGIN_V2     = AUTH_API_V2 + "login"
)

const (
	CLIENTS_V2      = API_V2 + "clients"
	EXACT_CLIENT_V2 = CLIENTS_V2 + "/{clientId}"
)

const (
	REPETITORS_V2      = API_V2 + "repetitors"
	EXACT_REPETITOR_V2 = REPETITORS_V2 + "/{repetitorId}"
)

const (
	CHATS_V2               = API_V2 + "chats"
	EXACT_CHAT_V2          = CHATS_V2 + "/{chatId}"
	EXACT_CHAT_MESSAGES_V2 = EXACT_CHAT_V2 + "/messages"
)

const (
	MESSAGES_V2      = API_V2 + "messages"
	EXACT_MESSAGE_V2 = MESSAGES_V2 + "/{messageId}"
)

const (
	LESSONS_V2      = API_V2 + "lessons"
	EXACT_LESSON_V2 = LESSONS_V2 + "/{lessonId}"
)

const (
	CONTRACTS_V2             = API_V2 + "contracts"
	EXACT_CONTRACT_V2        = CONTRACTS_V2 + "/{contractId}"
	CONTRACT_LESSONS_V2      = EXACT_CONTRACT_V2 + "/lessons"
	CONTRACT_REVIEWS_V2      = EXACT_CONTRACT_V2 + "/reviews"
	CONTRACT_TRANSACTIONS_V2 = EXACT_CONTRACT_V2 + "/transactions"
)

const (
	TRANSACTIONS_V2         = API_V2 + "transactions"
	EXACT_TRANSACTION_V2    = TRANSACTIONS_V2 + "/{transactionId}"
	TRANSACTION_APPROVAL_V2 = EXACT_TRANSACTION_V2 + "/approval"
)

const (
	DEPARTMENTS_V2       = API_V2 + "departments"
	EXACT_DEPARTMENT_V2  = DEPARTMENTS_V2 + "/{departmentId}"
	ADMINS_V2            = API_V2 + "admins"
	EXACT_ADMIN_V2       = ADMINS_V2 + "/{adminId}"
	ADMIN_DEPARTMENTS_V2 = EXACT_ADMIN_V2 + "/departments"
	MODERATORS_V2        = API_V2 + "moderators"
	EXACT_MODERATOR_V2   = MODERATORS_V2 + "/{moderatorId}"
	MODERATOR_SALARY_V2  = EXACT_MODERATOR_V2 + "/salary"
	LEGACY_ARCHIVE_V2    = API_V2 + "legacy/"
)

const (
	DEPARTMENT_MODERATOR_V2 = DEPARTMENTS_V2 + "/{departmentId}/moderators/{moderatorId}"
)

const (
	STATIC_FILES_V2        = API_V2 + "static"
	EXACT_STATIC_FILE_V2   = STATIC_FILES_V2 + "/{filename}"
	OPENAPI_YAML_V2        = STATIC_FILES_V2 + "/openapi.yaml"
	README_V2              = STATIC_FILES_V2 + "/README.md"
	DOCUMENTATION_V2       = STATIC_FILES_V2 + "/documentation"
	RESERVED_FILES_V2      = API_V2 + "reserved"
	EXACT_RESERVED_FILE_V2 = RESERVED_FILES_V2 + "/{filename}"
)

const (
	ERR_CODE_BAD_REQUEST  = "bad_request"
	ERR_CODE_UNAUTHORIZED = "unauthorized"
	ERR_CODE_FORBIDDEN    = "forbidden"
	ERR_CODE_NOT_FOUND    = "not_found"
	ERR_CODE_CONFLICT     = "conflict"
	ERR_CODE_SERVER_ERROR = "server_error"
)

const (
	ERR_MSG_INVALID_PARAMS = "Invalid request parameters"
	ERR_MSG_INVALID_BODY   = "Invalid request body"
)

const (
	ERR_MSG_NO_TOKEN      = "Authorization token is missing"
	ERR_MSG_BAD_TOKEN     = "Token is invalid"
	ERR_MSG_TOKEN_EXPIRED = "Token expired"
)

const (
	ERR_MSG_NOT_PARTICIPANT = "User is not a participant of this chat"
	ERR_MSG_NOT_MODERATOR   = "Moderator role required"
	ERR_MSG_NOT_HEAD        = "Department head role required"
)

const (
	ERR_MSG_CHAT_NOT_FOUND        = "Chat not found"
	ERR_MSG_MESSAGE_NOT_FOUND     = "Message not found"
	ERR_MSG_DEPARTMENT_NOT_FOUND  = "Department not found"
	ERR_MSG_PROFILE_NOT_FOUND     = "Profile not found"
	ERR_MSG_CONTRACT_NOT_FOUND    = "Contract not found"
	ERR_MSG_TRANSACTION_NOT_FOUND = "Transaction not found"
	ERR_MSG_LESSON_NOT_FOUND      = "Lesson not found"
	ERR_MSG_REVIEW_NOT_FOUND      = "Review not found"
)

const (
	ERR_MSG_CHAT_ALREADY_EXISTS             = "Chat already exists for these participants"
	ERR_MSG_DEPARTMENT_NAME_TAKEN           = "Department name already in use"
	ERR_MSG_MODERATOR_ALREADY_IN_DEPARTMENT = "Moderator already assigned to this department"
)

const (
	ERR_MSG_SERVER_UNEXPECTED = "Unexpected server error"
)
