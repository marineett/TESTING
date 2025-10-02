import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Form, Input, Button, Card, Typography, Alert } from 'antd';
import { UserOutlined, LockOutlined } from '@ant-design/icons';
import { UserType } from '../types/front_types';
import { API_ENDPOINTS } from '../config';

const { Title } = Typography;

const Login: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    const userId = localStorage.getItem('user_id');
    if (userId) {
      console.log('user_id from localStorage:', userId);
    }
  }, []);

  const onFinish = async (values: any) => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(API_ENDPOINTS.AUTH.LOGIN, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          login: values.username,
          password: values.password,
        }),
      });
      if (!response.ok) {
        const errorText = await response.text();
        setError(errorText || 'Ошибка авторизации');
        setLoading(false);
        return;
      }
      const verdict = await response.json();
      if (
        verdict.user_type === undefined ||
        verdict.user_type === UserType.Unauthorized
      ) {
        setError('Ошибка авторизации');
        setLoading(false);
        return;
      }
      localStorage.setItem('user_id', verdict.user_id);
      if (verdict.user_type === UserType.Repetitor) {
        localStorage.setItem('repetitor_id', verdict.user_id);
      }
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
          navigate('/dashboard');
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
          Login
        </Title>
        {error && <Alert message={error} type="error" showIcon style={{ marginBottom: 16 }} />}
        <Form
          name="login"
          initialValues={{ remember: true }}
          onFinish={onFinish}
          layout="vertical"
        >
          <Form.Item
            name="username"
            rules={[{ required: true, message: 'Please input your username!' }]}
          >
            <Input
              prefix={<UserOutlined />}
              placeholder="Username"
              size="large"
              disabled={loading}
            />
          </Form.Item>

          <Form.Item
            name="password"
            rules={[{ required: true, message: 'Please input your password!' }]}
          >
            <Input.Password
              prefix={<LockOutlined />}
              placeholder="Password"
              size="large"
              disabled={loading}
            />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" block size="large" loading={loading}>
              Log in
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
};

export default Login; 