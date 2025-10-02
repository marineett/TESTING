export interface PersonalData {
    telephone_number: string;
    email: string;
    first_name: string;
    last_name: string;
    middle_name: string;
    passport_number: string;
    passport_series: string;
    passport_date: string;
    passport_issued_by: string;
}

export interface AuthData {
    login: string;
    password: string;
}

export enum UserType {
    Unauthorized = 0,
    Client = 1,
    Repetitor = 2,
    Moderator = 3,
    Admin = 4
}

export interface RegistrationInfo {
    personal_data: PersonalData;
    auth_data: AuthData;
    user_type: UserType;
}

export interface AuthVerdict {
    success: boolean;
    user_type: UserType;
    message?: string;
} 