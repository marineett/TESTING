import React from 'react';
import { useNavigate } from 'react-router-dom';
import TransactionList from '../components/TransactionList';
import ContractList from '../components/ContractList';
import ModeratorProfile from '../components/ModeratorProfile';
import ChatListWindow from '../components/ChatListWindow';
import { Button, Layout, Typography } from 'antd';
import { API_ENDPOINTS } from '../config';

const { Header, Content } = Layout;
const { Title } = Typography;

const ModeratorPage: React.FC = () => {
    const navigate = useNavigate();
    const moderatorId = localStorage.getItem('user_id');

    const handleLogout = () => {
        localStorage.removeItem('userType');
        localStorage.removeItem('user_id');
        navigate('/login');
    };

    return (
        <Layout style={{ minHeight: '100vh' }}>
            <Header style={{
                position: 'fixed',
                zIndex: 1,
                width: '100%',
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center'
            }}>
                <Title level={3} style={{ color: 'white', margin: 0 }}>Кабинет модератора</Title>
                <Button type="primary" onClick={handleLogout}>Выйти</Button>
            </Header>
            <Content style={{ padding: '24px', marginTop: 64 }}>
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                    <div style={{ padding: 24, background: '#fff', minHeight: 360, borderRadius: 10 }}>
                        <ChatListWindow userId={moderatorId ? Number(moderatorId) : 0} role="moderator" />
                        {moderatorId && <ModeratorProfile id={moderatorId} />}
                    </div>
                <div className="mb-8 mt-8">
                    <h2 className="text-2xl font-bold mb-6 text-gray-800">Транзакции для подтверждения</h2>
                    <TransactionList />
                </div>
                <div className="mb-8">
                    <h2 className="text-2xl font-bold text-gray-900 mb-4">Контракты</h2>
                    <ContractList />
                </div>
            </div>
            </Content>
        </Layout>
    );
};

export default ModeratorPage; 