import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { InitUserData, UserType } from '../types/user';
import { API_ENDPOINTS } from '../config';

const RegistrationForm: React.FC = () => {
  const [formData, setFormData] = useState<InitUserData>({
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

  const [salary, setSalary] = useState<string>('');
  const [message, setMessage] = useState<{ type: 'success' | 'error' | null; text: string }>({
    type: null,
    text: ''
  });
  const [isSubmitting, setIsSubmitting] = useState(false);

  const navigate = useNavigate();

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    
    if (name === 'user_type') {
      setFormData(prev => ({
        ...prev,
        user_type: parseInt(value) as UserType
      }));
    } else if (name === 'salary') {
      setSalary(value);
    } else if (name.startsWith('auth_')) {
      setFormData(prev => ({
        ...prev,
        auth_data: {
          ...prev.auth_data,
          [name.replace('auth_', '')]: value,
        },
      }));
    } else {
      setFormData(prev => ({
        ...prev,
        personal_data: {
          ...prev.personal_data,
          [name]: value,
        },
      }));
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setMessage({ type: null, text: '' });
    setIsSubmitting(true);
    
    try {
      let endpoint = '';
      let dataToSend = {};

      const formattedData = {
        ...formData,
        personal_data: {
          ...formData.personal_data,
          passport_date: formData.personal_data.passport_date ? 
            new Date(formData.personal_data.passport_date).toISOString() : 
            ''
        }
      };

      switch (formData.user_type) {
        case UserType.Client:
          endpoint = API_ENDPOINTS.AUTH.REGISTER.CLIENT;
          dataToSend = formattedData;
          break;
        case UserType.Repetitor:
          endpoint = API_ENDPOINTS.AUTH.REGISTER.REPETITOR;
          dataToSend = formattedData;
          break;
        case UserType.Moderator:
          endpoint = API_ENDPOINTS.AUTH.REGISTER.MODERATOR;
          dataToSend = formattedData;
          break;
        case UserType.Admin:
          endpoint = API_ENDPOINTS.AUTH.REGISTER.ADMIN;
          dataToSend = {
            ...formattedData,
            salary: parseInt(salary)
          };
          break;
      }

      const response = await fetch(endpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(formattedData),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Ошибка при регистрации');
      }

      setMessage({
        type: 'success',
        text: 'Регистрация успешна! Перенаправление на страницу входа...'
      });
      
      setTimeout(() => {
        navigate('/login');
      }, 2000);
      
    } catch (error) {
      console.error('Registration error:', error);
      setMessage({
        type: 'error',
        text: error instanceof Error ? error.message : 'Произошла ошибка при регистрации'
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="max-w-lg mx-auto mt-8">
      <div className="bg-blue-50 p-4 rounded-lg mb-6">
        <h1 className="text-xl font-semibold text-blue-800 mb-2">Добро пожаловать в систему!</h1>
        <p className="text-blue-600">Пожалуйста, заполните форму регистрации, чтобы создать новый аккаунт.</p>
      </div>
      
      <form onSubmit={handleSubmit} className="p-6 bg-white rounded shadow">
        <h2 className="text-2xl font-bold mb-6">Регистрация</h2>
        
        {/* Сообщение об успехе/ошибке */}
        {message.type && (
          <div className={`mb-4 p-4 rounded ${
            message.type === 'success' ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'
          }`}>
            {message.text}
          </div>
        )}

        {/* Выбор типа пользователя */}
        <div className="mb-6">
          <label className="block mb-2">Тип пользователя:</label>
          <select
            name="user_type"
            value={formData.user_type}
            onChange={handleChange}
            className="w-full p-2 border rounded"
            disabled={isSubmitting}
          >
            <option value={UserType.Client}>Клиент</option>
            <option value={UserType.Repetitor}>Репетитор</option>
            <option value={UserType.Moderator}>Модератор</option>
            <option value={UserType.Admin}>Администратор</option>
          </select>
        </div>
        
        {/* Данные авторизации */}
        <div className="mb-6">
          <h3 className="text-xl font-semibold mb-4">Данные для входа</h3>
          <div className="space-y-4">
            <div>
              <label className="block mb-2">Логин:</label>
              <input
                type="text"
                name="auth_login"
                value={formData.auth_data.login}
                onChange={handleChange}
                className="w-full p-2 border rounded"
                required
                disabled={isSubmitting}
              />
            </div>
            <div>
              <label className="block mb-2">Пароль:</label>
              <input
                type="password"
                name="auth_password"
                value={formData.auth_data.password}
                onChange={handleChange}
                className="w-full p-2 border rounded"
                required
                disabled={isSubmitting}
              />
            </div>
          </div>
        </div>

        {/* Персональные данные */}
        <div className="mb-6">
          <h3 className="text-xl font-semibold mb-4">Личные данные</h3>
          <div className="space-y-4">
            <div>
              <label className="block mb-2">Имя:</label>
              <input
                type="text"
                name="first_name"
                value={formData.personal_data.first_name}
                onChange={handleChange}
                className="w-full p-2 border rounded"
                required
                disabled={isSubmitting}
              />
            </div>
            <div>
              <label className="block mb-2">Фамилия:</label>
              <input
                type="text"
                name="last_name"
                value={formData.personal_data.last_name}
                onChange={handleChange}
                className="w-full p-2 border rounded"
                required
                disabled={isSubmitting}
              />
            </div>
            <div>
              <label className="block mb-2">Отчество:</label>
              <input
                type="text"
                name="middle_name"
                value={formData.personal_data.middle_name}
                onChange={handleChange}
                className="w-full p-2 border rounded"
                disabled={isSubmitting}
              />
            </div>
            <div>
              <label className="block mb-2">Email:</label>
              <input
                type="email"
                name="email"
                value={formData.personal_data.email}
                onChange={handleChange}
                className="w-full p-2 border rounded"
                required
                disabled={isSubmitting}
              />
            </div>
            <div>
              <label className="block mb-2">Телефон:</label>
              <input
                type="tel"
                name="telephone_number"
                value={formData.personal_data.telephone_number}
                onChange={handleChange}
                className="w-full p-2 border rounded"
                required
                disabled={isSubmitting}
              />
            </div>
          </div>
        </div>

        {/* Паспортные данные */}
        <div className="mb-6">
          <h3 className="text-xl font-semibold mb-4">Паспортные данные</h3>
          <div className="space-y-4">
            <div>
              <label className="block mb-2">Серия паспорта:</label>
              <input
                type="text"
                name="passport_series"
                value={formData.personal_data.passport_series}
                onChange={handleChange}
                className="w-full p-2 border rounded"
                required
                disabled={isSubmitting}
              />
            </div>
            <div>
              <label className="block mb-2">Номер паспорта:</label>
              <input
                type="text"
                name="passport_number"
                value={formData.personal_data.passport_number}
                onChange={handleChange}
                className="w-full p-2 border rounded"
                required
                disabled={isSubmitting}
              />
            </div>
            <div>
              <label className="block mb-2">Дата выдачи:</label>
              <input
                type="date"
                name="passport_date"
                value={formData.personal_data.passport_date}
                onChange={handleChange}
                className="w-full p-2 border rounded"
                required
                disabled={isSubmitting}
              />
            </div>
            <div>
              <label className="block mb-2">Кем выдан:</label>
              <input
                type="text"
                name="passport_issued_by"
                value={formData.personal_data.passport_issued_by}
                onChange={handleChange}
                className="w-full p-2 border rounded"
                required
                disabled={isSubmitting}
              />
            </div>
          </div>
        </div>

        {/* Зарплата для администратора */}
        {formData.user_type === UserType.Admin && (
          <div className="mb-6">
            <h3 className="text-xl font-semibold mb-4">Данные администратора</h3>
            <div className="space-y-4">
              <div>
                <label className="block mb-2">Зарплата:</label>
                <input
                  type="number"
                  name="salary"
                  value={salary}
                  onChange={handleChange}
                  className="w-full p-2 border rounded"
                  required
                  disabled={isSubmitting}
                />
              </div>
            </div>
          </div>
        )}

        <button
          type="submit"
          className={`w-full py-2 px-4 rounded text-white ${
            isSubmitting 
              ? 'bg-gray-400 cursor-not-allowed' 
              : 'bg-blue-500 hover:bg-blue-600'
          }`}
          disabled={isSubmitting}
        >
          {isSubmitting ? 'Регистрация...' : 'Зарегистрироваться'}
        </button>
      </form>
    </div>
  );
};

export default RegistrationForm; 