import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { RegistrationInfo, UserType } from '../types/front_types';
import { API_ENDPOINTS } from '../config';

const RegistrationPage: React.FC = () => {
    const navigate = useNavigate();
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);
    const [isLoading, setIsLoading] = useState(false);
    const [initData, setInitData] = useState<RegistrationInfo>({
        personal_data: {
            telephone_number: '',
            email: '',
            first_name: '',
            last_name: '',
            middle_name: '',
            passport_number: '',
            passport_series: '',
            passport_date: '',
            passport_issued_by: ''
        },
        auth_data: {
            login: '',
            password: ''
        },
        user_type: UserType.Client
    });

    const getRegistrationEndpoint = (userType: UserType): string => {
        switch (userType) {
            case UserType.Client:
                return API_ENDPOINTS.AUTH.REGISTER.CLIENT;
            case UserType.Repetitor:
                return API_ENDPOINTS.AUTH.REGISTER.REPETITOR;
            case UserType.Moderator:
                return API_ENDPOINTS.AUTH.REGISTER.MODERATOR;
            case UserType.Admin:
                return API_ENDPOINTS.AUTH.REGISTER.ADMIN;
            default:
                return API_ENDPOINTS.AUTH.REGISTER.CLIENT;
        }
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsLoading(true);
        setError(null);
        setSuccess(null);

        try {
            const endpoint = getRegistrationEndpoint(initData.user_type);
            const response = await fetch(endpoint, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(initData),
            });

            if (response.ok) {
                setSuccess('Регистрация успешно завершена!');
                setTimeout(() => {
                    navigate('/login');
                }, 2000);
                return;
            }

            const text = await response.text();
            setError(text || 'Произошла ошибка при регистрации');
            
        } catch (error) {
            setError(error instanceof Error ? error.message : 'Unknown error');
        } finally {
            setIsLoading(false);
        }
    };

    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
        const { name, value } = e.target;
        if (name === 'user_type') {
            setInitData(prev => ({
                ...prev,
                user_type: Number(value) as UserType
            }));
        } else if (name.startsWith('personal_')) {
            const fieldName = name.replace('personal_', '');
            if (fieldName === 'passport_date') {
                const dateValue = value ? value + "T00:00:00Z" : "0001-01-01T00:00:00Z";
                setInitData(prev => ({
                    ...prev,
                    personal_data: {
                        ...prev.personal_data,
                        [fieldName]: dateValue
                    }
                }));
            } else {
                setInitData(prev => ({
                    ...prev,
                    personal_data: {
                        ...prev.personal_data,
                        [fieldName]: value
                    }
                }));
            }
        } else if (name.startsWith('auth_')) {
            setInitData(prev => ({
                ...prev,
                auth_data: {
                    ...prev.auth_data,
                    [name.replace('auth_', '')]: value
                }
            }));
        }
    };

    return (
        <div className="min-h-screen bg-gray-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
            <div className="sm:mx-auto sm:w-full sm:max-w-md">
                <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
                    Регистрация
                </h2>
            </div>

            <div className="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
                <div className="bg-white py-8 px-4 shadow sm:rounded-lg sm:px-10">
                    {/* Сообщение об успехе */}
                    {success && (
                        <div className="mb-4 p-4 rounded-md bg-green-50 text-green-700">
                            {success}
                        </div>
                    )}

                    {/* Сообщение об ошибке */}
                    {error && (
                        <div className="mb-4 p-4 rounded-md bg-red-50 text-red-700 font-mono">
                            {error}
                        </div>
                    )}

                    <form className="space-y-6" onSubmit={handleSubmit}>
                        <div>
                            <label className="block text-sm font-medium text-gray-700">
                                Тип пользователя
                            </label>
                            <select
                                name="user_type"
                                value={initData.user_type}
                                onChange={handleChange}
                                className="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md"
                                required
                                disabled={isLoading}
                            >
                                <option value={UserType.Client}>Клиент</option>
                                <option value={UserType.Repetitor}>Репетитор</option>
                                <option value={UserType.Moderator}>Модератор</option>
                                <option value={UserType.Admin}>Администратор</option>
                            </select>
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700">
                                Имя
                            </label>
                            <input
                                type="text"
                                name="personal_first_name"
                                value={initData.personal_data.first_name}
                                onChange={handleChange}
                                className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                                required
                                disabled={isLoading}
                            />
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700">
                                Фамилия
                            </label>
                            <input
                                type="text"
                                name="personal_last_name"
                                value={initData.personal_data.last_name}
                                onChange={handleChange}
                                className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                                required
                                disabled={isLoading}
                            />
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700">
                                Отчество
                            </label>
                            <input
                                type="text"
                                name="personal_middle_name"
                                value={initData.personal_data.middle_name}
                                onChange={handleChange}
                                className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                                disabled={isLoading}
                            />
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700">
                                Email
                            </label>
                            <input
                                type="email"
                                name="personal_email"
                                value={initData.personal_data.email}
                                onChange={handleChange}
                                className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                                required
                                disabled={isLoading}
                            />
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700">
                                Телефон
                            </label>
                            <input
                                type="tel"
                                name="personal_telephone_number"
                                value={initData.personal_data.telephone_number}
                                onChange={handleChange}
                                className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                                required
                                disabled={isLoading}
                            />
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700">
                                Дата выдачи паспорта
                            </label>
                            <input
                                type="date"
                                name="personal_passport_date"
                                value={initData.personal_data.passport_date ? initData.personal_data.passport_date.split('T')[0] : ''}
                                onChange={handleChange}
                                className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                                disabled={isLoading}
                            />
                            <p className="mt-1 text-sm text-gray-500">
                                Выберите дату выдачи паспорта (необязательно)
                            </p>
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700">
                                Логин
                            </label>
                            <input
                                type="text"
                                name="auth_login"
                                value={initData.auth_data.login}
                                onChange={handleChange}
                                className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                                required
                                disabled={isLoading}
                            />
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700">
                                Пароль
                            </label>
                            <input
                                type="password"
                                name="auth_password"
                                value={initData.auth_data.password}
                                onChange={handleChange}
                                className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                                required
                                disabled={isLoading}
                            />
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700">
                                Серия паспорта
                            </label>
                            <input
                                type="text"
                                name="personal_passport_series"
                                value={initData.personal_data.passport_series}
                                onChange={handleChange}
                                className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                                required
                                disabled={isLoading}
                            />
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700">
                                Номер паспорта
                            </label>
                            <input
                                type="text"
                                name="personal_passport_number"
                                value={initData.personal_data.passport_number}
                                onChange={handleChange}
                                className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                                required
                                disabled={isLoading}
                            />
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700">
                                Кем выдан паспорт
                            </label>
                            <input
                                type="text"
                                name="personal_passport_issued_by"
                                value={initData.personal_data.passport_issued_by}
                                onChange={handleChange}
                                className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                                required
                                disabled={isLoading}
                            />
                        </div>

                        <div>
                            <button
                                type="submit"
                                disabled={isLoading}
                                className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
                            >
                                {isLoading ? 'Регистрация...' : 'Зарегистрироваться'}
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    );
};

export default RegistrationPage; 