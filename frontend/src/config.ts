export const API_ENDPOINTS = {
    AUTH: {
        LOGIN: '/api/auth/authorize',
        REGISTER: {
            CLIENT: '/api/registration/client',
            REPETITOR: '/api/registration/repetitor',
            MODERATOR: '/api/registration/moderator',
            ADMIN: '/api/registration/admin',
        }
    },
    CLIENT: {
        CONTRACTS: (clientId: number, offset: number, limit: number, status: number) => 
            `/api/client/get_contracts?client_id=${clientId}&offset=${offset}&limit=${limit}&status=${status}`,
        CREATE_CONTRACT: '/api/client/create_contract',
        MAKE_REVIEW: '/api/client/make_review',
        GET_REVIEWS: (repetitorId: number) => `/api/contract/get_reviews?repetitor_id=${repetitorId}`,
        GET_REVIEW: '/api/contract/get_review',
        GET_PROFILE: (clientId: number, offset: number, limit: number) =>
            `/api/client/get_profile?id=${clientId}&reviews_offset=${offset}&reviews_limit=${limit}`
    },
    REPETITOR: {
        GET_CONTRACTS: (repetitorId: number, offset: number, limit: number, status: number) =>
            `/api/repetitor/get_contracts?repetitor_id=${repetitorId}&offset=${offset}&limit=${limit}&status=${status}`,
        GET_AVAILABLE_CONTRACTS: (offset: number, limit: number, status: number) =>
            `/api/repetitor/get_available_contracts?offset=${offset}&limit=${limit}&status=${status}`,
        RESPOND_TO_CONTRACT: (contractId: number) =>
            `/api/repetitor/respond_to_contract?contract_id=${contractId}`,
        ACCEPT_CONTRACT: (repetitorId: number, contractId: number) =>
            `/api/repetitor/accept_contract?repetitor_id=${repetitorId}&contract_id=${contractId}`,
        GET_REVIEWS: (repetitorId: number) => `/api/repetitor/get_reviews?repetitor_id=${repetitorId}`,
        CONTRACTS: (repetitorId: number, offset: number, limit: number, status: number) => 
            `/api/contract/get_contracts?repetitor_id=${repetitorId}&offset=${offset}&limit=${limit}&status=${status}`,
        UPDATE_CONTRACT_STATUS: '/api/contract/update_status',
        MAKE_REVIEW: '/api/repetitor/make_review',
        PAY_FOR_CONTRACT: (contractId: number, amount: number) =>
            `/api/repetitor/pay_for_contract?contract_id=${contractId}&amount=${amount}`,
        GET_PROFILE: (id: number, reviews_offset: number, reviews_limit: number) => `/api/repetitor/get_profile?id=${id}&reviews_offset=${reviews_offset}&reviews_limit=${reviews_limit}`,
        CHANGE_RESUME: (id: number) => `/api/repetitor/change_resume?id=${id}`,
        CANCEL_CONTRACT: (repetitorId: number, contractId: number) =>
            `/api/repetitor/cancel_contract?id=${repetitorId}&c_id=${contractId}`,
        GET_LESSONS: (contractId: number, offset: number, limit: number) =>
            `/api/contract/get_lessons?contract_id=${contractId}&lessons_offset=${offset}&lessons_size=${limit}`,
        ADD_LESSON: '/api/contract/add_lesson',
        COMPLETE_CONTRACT: (repetitorId: number, contractId: number) =>
            `/api/repetitor/complete_contract?id=${repetitorId}&c_id=${contractId}`
    },
    MODERATOR: {
        CONTRACTS: (moderatorId: number, offset: number, limit: number) => 
            `/api/moderator/contracts?moderator_id=${moderatorId}&offset=${offset}&limit=${limit}`,
        GET_TRANSACTION_TO_APPROVE: '/api/moderator/get_transaction_to_approve',
        APPROVE_TRANSACTION: (transactionId: number) =>
            `/api/moderator/approve_transaction?transaction_id=${transactionId}`,
        REJECT_TRANSACTION: (transactionId: number) =>
            `/api/moderator/reject_transaction?transaction_id=${transactionId}`,
        GET_PROFILE: (id: number) => `/api/moderator/get_profile?id=${id}`,
        GET_CONTRACTS: (from: number, size: number) =>
            `/api/moderator/get_contracts?from=${from}&size=${size}`,
        BAN_CONTRACT: (contractId: number) =>
            `/api/moderator/ban_contract?id=${contractId}`
    },
    ADMIN: {
        CONTRACTS: (adminId: number, offset: number, limit: number) => 
            `/api/admin/contracts?admin_id=${adminId}&offset=${offset}&limit=${limit}`,
        GET_PROFILE: (id: number) => `/api/admin/get_profile?id=${id}`,
        GET_DEPARTMENTS: (id: number) => `/api/admin/get_departments?id=${id}`,
        GET_MODERATORS: (id: number) => `/api/admin/get_moderators?id=${id}`,
        HIRE_MODERATOR: (adminId: number, departmentId: number, moderatorId: number) =>
            `/api/admin/hire_moderator?id=${adminId}&d_id=${departmentId}&m_id=${moderatorId}`,
        FIRE_MODERATOR: (adminId: number, departmentId: number, moderatorId: number) =>
            `/api/admin/fire_moderator?id=${adminId}&d_id=${departmentId}&m_id=${moderatorId}`,
        CHANGE_MODERATOR_SALARY: (adminId: number, departmentId: number, moderatorId: number, salary: string) =>
            `/api/admin/change_moderator_salary?id=${adminId}&d_id=${departmentId}&m_id=${moderatorId}&salary=${salary}`,
        CREATE_DEPARTMENT: (adminId: number, name: string) =>
            `/api/admin/create_department?id=${adminId}&name=${encodeURIComponent(name)}`
    },
    CHAT: {
        START_CR_CHAT: (clientId: number, repetitorId: number) =>
            `/api/chat/start_cr_chat?c_id=${clientId}&r_id=${repetitorId}`,
        START_CM_CHAT: (clientId: number, moderatorId: number) =>
            `/api/chat/start_cm_chat?c_id=${clientId}&m_id=${moderatorId}`,
        START_RM_CHAT: (repetitorId: number, moderatorId: number) =>
            `/api/chat/start_rm_chat?r_id=${repetitorId}&m_id=${moderatorId}`,
        GET_CHAT: (chatId: number) =>
            `/api/chat/get_chat?id=${chatId}`,
        SEND_MESSAGE: (senderId: number, chatId: number) =>
            `/api/chat/send_message?sender_id=${senderId}&chat_id=${chatId}`,
        GET_MESSAGES: (chatId: number, offset: number, limit: number) =>
            `/api/chat/get_messages?id=${chatId}&messages_offset=${offset}&messages_limit=${limit}`,
        GET_CHATS: (userId: number, offset: number, limit: number) =>
            `/api/chat/get_chats?user_id=${userId}&offset=${offset}&limit=${limit}`,
        GET_CLIENT_CHATS: (id: number, offset: number, limit: number) =>
            `/api/chat/get_client_chats?id=${id}&chats_offset=${offset}&chats_limit=${limit}`,
        GET_REPETITOR_CHATS: (id: number, offset: number, limit: number) =>
            `/api/chat/get_repetitor_chats?id=${id}&chats_offset=${offset}&chats_limit=${limit}`,
        GET_MODERATOR_CHATS: (id: number, offset: number, limit: number) =>
            `/api/chat/get_moderator_chats?id=${id}&chats_offset=${offset}&chats_limit=${limit}`
    },
    GUEST: {
        GET_REPETITORS: (offset: number, limit: number) =>
            `/api/guest/get_repetitors?repetitors_offset=${offset}&repetitors_limit=${limit}`
    }
}; 