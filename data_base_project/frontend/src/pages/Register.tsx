import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Form, Input, Button, Card, Typography, Alert, Select } from 'antd';
import { UserOutlined, LockOutlined, MailOutlined } from '@ant-design/icons';
import { UserType } from '../types/front_types';
import { API_ENDPOINTS } from '../config';

const { Title } = Typography;
const { Option } = Select;

const Register: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const navigate = useNavigate();

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

    const onFinish = async (values: any) => {
        setLoading(true);
        setError(null);
        try {
            const endpoint = getRegistrationEndpoint(values.userType);
            const response = await fetch(endpoint, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    login: values.username,
                    password: values.password,
                    email: values.email,
                }),
            });

            if (!response.ok) {
                const errorText = await response.text();
                setError(errorText || 'Ошибка регистрации');
                setLoading(false);
                return;
            }

            const data = await response.json();
            if (data.success) {
                navigate('/login');
            } else {
                setError(data.message || 'Ошибка регистрации');
            }
        } catch (err) {
            setError('Ошибка соединения с сервером');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div style={{ padding: '24px', maxWidth: '400px', margin: '0 auto' }}>
            <Card>
                <Title level={2} style={{ textAlign: 'center', marginBottom: '24px' }}>
                    Регистрация
                </Title>
                {error && <Alert message={error} type="error" showIcon style={{ marginBottom: 16 }} />}
                <Form
                    name="register"
                    initialValues={{ remember: true }}
                    onFinish={onFinish}
                    layout="vertical"
                >
                    <Form.Item
                        name="userType"
                        rules={[{ required: true, message: 'Пожалуйста, выберите тип пользователя!' }]}
                    >
                        <Select placeholder="Тип пользователя" disabled={loading}>
                            <Option value={UserType.Client}>Клиент</Option>
                            <Option value={UserType.Repetitor}>Репетитор</Option>
                            <Option value={UserType.Moderator}>Модератор</Option>
                            <Option value={UserType.Admin}>Администратор</Option>
                        </Select>
                    </Form.Item>

                    <Form.Item
                        name="username"
                        rules={[{ required: true, message: 'Пожалуйста, введите логин!' }]}
                    >
                        <Input
                            prefix={<UserOutlined />}
                            placeholder="Логин"
                            size="large"
                            disabled={loading}
                        />
                    </Form.Item>

                    <Form.Item
                        name="email"
                        rules={[
                            { required: true, message: 'Пожалуйста, введите email!' },
                            { type: 'email', message: 'Введите корректный email!' }
                        ]}
                    >
                        <Input
                            prefix={<MailOutlined />}
                            placeholder="Email"
                            size="large"
                            disabled={loading}
                        />
                    </Form.Item>

                    <Form.Item
                        name="password"
                        rules={[{ required: true, message: 'Пожалуйста, введите пароль!' }]}
                    >
                        <Input.Password
                            prefix={<LockOutlined />}
                            placeholder="Пароль"
                            size="large"
                            disabled={loading}
                        />
                    </Form.Item>

                    <Form.Item>
                        <Button type="primary" htmlType="submit" block size="large" loading={loading}>
                            Зарегистрироваться
                        </Button>
                    </Form.Item>
                </Form>
            </Card>
        </div>
    );
};

export default Register; 