import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import AdminProfile from '../components/AdminProfile';
import DepartmentList from '../components/DepartmentList';
import ModeratorList from '../components/ModeratorList';
import { Layout, Typography, Button } from 'antd';
import { API_ENDPOINTS } from '../config';

const { Header, Content } = Layout;
const { Title } = Typography;

const AdminPage: React.FC = () => {
    const navigate = useNavigate();
    const adminId = localStorage.getItem('user_id');
    const [departmentName, setDepartmentName] = useState('');
    const [message, setMessage] = useState<string | null>(null);
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        console.log('AdminPage: ModeratorList должен отображаться');
    }, []);

    const handleLogout = () => {
        localStorage.removeItem('userType');
        navigate('/login');
    };

    const handleCreateDepartment = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!departmentName.trim() || !adminId) return;
        setLoading(true);
        setMessage(null);
        try {
            const response = await fetch(
                API_ENDPOINTS.ADMIN.CREATE_DEPARTMENT(Number(adminId), departmentName),
                {
                    method: 'GET',
                    headers: {
                        'Authorization': `Bearer ${localStorage.getItem('token')}`
                    }
                }
            );
            if (!response.ok) throw new Error('Ошибка при создании отдела');
            setMessage('Отдел успешно создан!');
            setDepartmentName('');
        } catch (err) {
            setMessage('Ошибка при создании отдела');
        } finally {
            setLoading(false);
        }
    };

    return (
        <Layout style={{ minHeight: '100vh' }}>
            <Header style={{
                position: 'fixed',
                zIndex: 1,
                width: '100%',
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                background: '#0a1929',
                boxShadow: '0 1px 2px 0 rgba(0,0,0,0.05)'
            }}>
                <Title level={3} style={{ color: 'white', margin: 0 }}>Панель администратора</Title>
                <Button type="primary" onClick={handleLogout} style={{ fontWeight: 500, fontSize: 16 }}>Выйти</Button>
            </Header>
            <Content style={{ padding: '24px', marginTop: 64 }}>
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                    <div style={{ padding: 24, background: '#fff', minHeight: 360, borderRadius: 10 }}>
                {adminId && <AdminProfile id={adminId} />}
                <div className="mt-8 max-w-md bg-white p-6 rounded shadow">
                    <h2 className="text-lg font-semibold mb-4">Создать отдел</h2>
                    <form onSubmit={handleCreateDepartment} className="flex flex-col gap-4">
                        <input
                            type="text"
                            className="border rounded px-3 py-2"
                            placeholder="Название отдела"
                            value={departmentName}
                            onChange={e => setDepartmentName(e.target.value)}
                            required
                        />
                        <button
                            type="submit"
                            className="bg-indigo-600 text-white px-4 py-2 rounded hover:bg-indigo-700 disabled:opacity-50"
                            disabled={loading}
                        >
                            {loading ? 'Создание...' : 'Создать'}
                        </button>
                        {message && (
                            <div className={`text-sm ${message.includes('успешно') ? 'text-green-600' : 'text-red-600'}`}>
                                {message}
                            </div>
                        )}
                    </form>
                </div>
                <ModeratorList />
                <DepartmentList />
            </div>
        </div>
            </Content>
        </Layout>
    );
};

export default AdminPage; 