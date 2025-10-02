import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Layout, Typography, Button, List, Avatar, Rate, Pagination } from 'antd';
import { UserOutlined } from '@ant-design/icons';
import { API_ENDPOINTS } from '../config';

const { Header, Content } = Layout;
const { Title } = Typography;

interface Repetitor {
    first_name: string;
    mean_rating: number;
}

const HomePage: React.FC = () => {
    console.log('HomePage mounted');
    const navigate = useNavigate();
    const [repetitors, setRepetitors] = useState<Repetitor[]>([]);
    const [loading, setLoading] = useState(false);
    const [total, setTotal] = useState(0);
    const [currentPage, setCurrentPage] = useState(1);
    const pageSize = 10;

    useEffect(() => {
        console.log('useEffect called, currentPage:', currentPage);
        fetchRepetitors();
    }, [currentPage]);

    const fetchRepetitors = async () => {
        setLoading(true);
        try {
            const offset = (currentPage - 1) * pageSize;
            const url = API_ENDPOINTS.GUEST.GET_REPETITORS(offset, pageSize);
            console.log('Fetching repetitors from:', url);
            const response = await fetch(url);
            console.log('Fetch response:', response);
            if (!response.ok) throw new Error('Failed to fetch repetitors');
            const data = await response.json();
            console.log('Received data:', data);
            if (Array.isArray(data)) {
                setRepetitors(data);
                setTotal(data.length);
            } else {
                setRepetitors(data.repetitors || []);
                setTotal(data.total || 0);
            }
        } catch (error) {
            console.error('Error fetching repetitors:', error);
        } finally {
            setLoading(false);
        }
    };

    const handleLogout = () => {
        localStorage.removeItem('userType');
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
                alignItems: 'center',
                background: '#0a1929',
                boxShadow: '0 1px 2px 0 rgba(0,0,0,0.05)'
            }}>
                <Title level={3} style={{ color: 'white', margin: 0 }}>Главная страница</Title>
                <Button type="primary" onClick={handleLogout} style={{ fontWeight: 500, fontSize: 16 }}>Выйти</Button>
            </Header>
            <Content style={{ padding: '24px', marginTop: 64 }}>
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                    <div style={{ padding: 24, background: '#fff', minHeight: 360, borderRadius: 10 }}>
                        <Title level={4}>Список репетиторов</Title>
                        <List
                            loading={loading}
                            itemLayout="horizontal"
                            dataSource={repetitors}
                            renderItem={item => (
                                <List.Item>
                                    <List.Item.Meta
                                        avatar={<Avatar icon={<UserOutlined />} />}
                                        title={item.first_name}
                                        description={<Rate disabled defaultValue={item.mean_rating} />}
                                    />
                                </List.Item>
                            )}
                        />
                        <div style={{ textAlign: 'center', marginTop: 16 }}>
                            <Pagination
                                current={currentPage}
                                total={total}
                                pageSize={pageSize}
                                onChange={setCurrentPage}
                            />
                        </div>
                    </div>
                </div>
            </Content>
        </Layout>
    );
};

export default HomePage; 