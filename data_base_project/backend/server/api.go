package server

const (
	API = "/api/"
)

const (
	USER_API         = API + "user/"
	GUEST_API        = API + "guest/"
	ADMIN_API        = API + "admin/"
	MODERATOR_API    = API + "moderator/"
	CLIENT_API       = API + "client/"
	REGISTRATION_API = API + "registration/"
	AUTH_API         = API + "auth/"
	CONTRACT_API     = API + "contract/"
	CHAT_API         = API + "chat/"
	REPETITOR_API    = API + "repetitor/"
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
)
