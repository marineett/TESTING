export enum UserType {
  Unauthorized = 0,
  Client = 1,
  Repetitor = 2,
  Moderator = 3,
  Admin = 4
}

export interface PassportData {
  passport_number: string;
  passport_series: string;
  passport_date: string; // будет преобразовано в Date на бэкенде
  passport_issued_by: string;
}

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

export interface RepetitorData {
  education: string;
  experience: string;
  subjects: string[];
  achievements: string;
  teaching_methods: string;
  price_per_hour: number;
}

export interface ModeratorData {
  department: string;
  responsibilities: string[];
  work_schedule: string;
}

export interface AdminData {
  access_level: string;
  department: string;
  responsibilities: string[];
}

export interface InitUserData {
  personal_data: PersonalData;
  auth_data: AuthData;
  user_type: UserType;
  repetitor_data?: RepetitorData;
  moderator_data?: ModeratorData;
  admin_data?: AdminData;
}

export interface UserData {
  id: number;
  registration_date: Date;
  last_login_date: Date;
  personal_data_id: number;
}

export interface AuthVerdict {
  user_id: number;
  user_type: UserType;
}

export interface AuthInfo {
  id: number;
  user_id: number;
  user_type: UserType;
  login: string;
  password: string;
}

export interface RegistrationData {
  username: string;
  password: string;
  user_type: UserType;
  registration_date: Date;
} 