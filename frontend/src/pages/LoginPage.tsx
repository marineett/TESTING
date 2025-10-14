import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { AuthData, UserType } from '../types/front_types';
import { API_ENDPOINTS } from '../config';

const LoginPage: React.FC = () => {
    const navigate = useNavigate();
    const [authData, setAuthData] = useState<AuthData>({
        login: '',
        password: ''
    });
    const [error, setError] = useState<string>('');
    const [isLoading, setIsLoading] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');
        setIsLoading(true);

        try {
            const response = await fetch(API_ENDPOINTS.AUTH.LOGIN, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(authData),
            });

            if (!response.ok) {
                const errorText = await response.text();
                setError(errorText || 'Ошибка авторизации');
                return;
            }

            const verdict = await response.json();
            
            if (verdict && typeof verdict.user_type === 'number') {
                if (verdict.user_type >= UserType.Client && verdict.user_type <= UserType.Admin) {
                    localStorage.setItem('userType', verdict.user_type.toString());
                    switch (verdict.user_type) {
                        case UserType.Client:
                            navigate('/client');
                            break;
                        case UserType.Repetitor:
                            navigate('/repetitor');
                            break;
                        case UserType.Moderator:
                            navigate('/moderator');
                            break;
                        case UserType.Admin:
                            navigate('/admin');
                            break;
                        default:
                            setError('Неизвестный тип пользователя');
                    }
                } else {
                    setError('Недопустимый тип пользователя');
                }
            } else {
                setError('Некорректный ответ от сервера');
            }
        } catch (error) {
            console.error('Ошибка авторизации:', error);
            setError('Произошла ошибка при попытке авторизации');
        } finally {
            setIsLoading(false);
        }
    };

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setAuthData(prev => ({
            ...prev,
            [name]: value
        }));
    };

    return (
        <div className="min-h-screen bg-gray-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
            <div className="sm:mx-auto sm:w-full sm:max-w-md">
                <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
                    Вход в систему
                </h2>
            </div>

            <div className="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
                <div className="bg-white py-8 px-4 shadow sm:rounded-lg sm:px-10">
                    {error && (
                        <div className="mb-4 p-4 rounded-md bg-red-50 text-red-700 font-mono">
                            {error}
                        </div>
                    )}

                    <form className="space-y-6" onSubmit={handleSubmit}>
                        <div>
                            <label className="block text-sm font-medium text-gray-700">
                                Логин
                            </label>
                            <input
                                type="text"
                                name="login"
                                value={authData.login}
                                onChange={handleChange}
                                required
                                disabled={isLoading}
                                className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                            />
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700">
                                Пароль
                            </label>
                            <input
                                type="password"
                                name="password"
                                value={authData.password}
                                onChange={handleChange}
                                required
                                disabled={isLoading}
                                className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                            />
                        </div>

                        <div>
                            <button
                                type="submit"
                                disabled={isLoading}
                                className={`w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white ${
                                    isLoading
                                        ? 'bg-indigo-400 cursor-not-allowed'
                                        : 'bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500'
                                }`}
                            >
                                {isLoading ? 'Вход...' : 'Войти'}
                            </button>
                        </div>
                    </form>

                    <div className="mt-6 text-center">
                        <p className="text-sm text-gray-600">
                            Нет аккаунта?{' '}
                            <a href="/register" className="font-medium text-indigo-600 hover:text-indigo-500">
                                Зарегистрироваться
                            </a>
                        </p>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default LoginPage; 