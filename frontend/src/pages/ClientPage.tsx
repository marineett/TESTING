import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Layout, Button, Typography } from 'antd';
import ContractFeed from '../components/ContractFeed';
import CreateContractForm from '../components/CreateContractForm';
import ChatListWindow from '../components/ChatListWindow';
import { API_ENDPOINTS } from '../config';

const { Header, Content } = Layout;
const { Title } = Typography;

const ClientPage: React.FC = () => {
    const navigate = useNavigate();
    const [clientId] = useState<number>(Number(localStorage.getItem('user_id')) || 1);
    const [showCreateForm, setShowCreateForm] = useState(false);
    const [feedKey, setFeedKey] = useState(0);
    const [profile, setProfile] = useState<any>(null);
    const [reviewBatch, setReviewBatch] = useState(1);
    const REVIEWS_PER_BATCH = 5;

    useEffect(() => {
        const fetchProfile = async () => {
            if (!clientId) return;
            const offset = Math.max(0, (reviewBatch - 1) * REVIEWS_PER_BATCH);
            const response = await fetch(API_ENDPOINTS.CLIENT.GET_PROFILE(clientId, offset, REVIEWS_PER_BATCH));
            if (response.ok) {
                setProfile(await response.json());
            }
        };
        fetchProfile();
    }, [clientId, reviewBatch]);

    const handleLogout = () => {
        localStorage.removeItem('userType');
        localStorage.removeItem('user_id');
        navigate('/login');
    };

    const handleContractCreated = () => {
        setShowCreateForm(false);
        setFeedKey(k => k + 1);
    };

    const handleShowMoreReviews = () => {
        setReviewBatch((prev) => prev + 1);
    };

    return (
        <Layout style={{ minHeight: '100vh' }}>
            <Header style={{ position: 'fixed', zIndex: 1, width: '100%', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <Title level={3} style={{ color: 'white', margin: 0 }}>Кабинет клиента</Title>
                <Button type="primary" onClick={handleLogout}>
                    Выйти
                </Button>
            </Header>
            <Content style={{ padding: '24px', marginTop: 64 }}>
                <div style={{ padding: 24, background: '#fff', minHeight: 360 }}>
                    <ChatListWindow userId={clientId} role="client" />
                    {profile && (
                        <div style={{ marginBottom: 24, padding: 16, background: '#f5f5f5', borderRadius: 8 }}>
                            <div style={{ fontSize: 20, fontWeight: 600, marginBottom: 8 }}>
                                {profile.last_name} {profile.first_name} {profile.middle_name}
                            </div>
                            <div style={{ color: '#888', marginBottom: 4 }}>Email: {profile.email}</div>
                            <div style={{ color: '#888', marginBottom: 4 }}>Телефон: {profile.telephone_number}</div>
                            <div style={{ color: '#888', marginBottom: 8 }}>Средний рейтинг: {profile.mean_rating?.toFixed(2) ?? '—'}</div>
                            <div style={{ fontWeight: 500, marginBottom: 4 }}>Отзывы:</div>
                            {profile.reviews && profile.reviews.length > 0 ? (
                                <>
                                    {profile.reviews.map((review: any, idx: number) => (
                                        <div key={idx} style={{ background: '#fff', borderRadius: 6, padding: 10, marginBottom: 8, boxShadow: '0 1px 4px #eee' }}>
                                            <div style={{ fontWeight: 500 }}>Оценка: {review.rating}</div>
                                            <div style={{ color: '#555' }}>{review.comment}</div>
                                            <div style={{ fontSize: 12, color: '#aaa' }}>{review.created_at ? `Создано: ${new Date(review.created_at).toLocaleString()}` : ''}</div>
                                        </div>
                                    ))}
                                    {profile.reviews.length === REVIEWS_PER_BATCH && (
                                        <button onClick={handleShowMoreReviews} style={{ marginTop: 8, padding: '6px 16px', borderRadius: 4, background: '#e5e7eb', border: 'none', cursor: 'pointer' }}>
                                            Показать ещё
                                        </button>
                                    )}
                                </>
                            ) : (
                                <div style={{ color: '#aaa' }}>Нет отзывов</div>
                            )}
                        </div>
                    )}
                    <Title level={2}>Добро пожаловать в кабинет клиента</Title>
                    <p>Здесь вы можете управлять своими заказами и взаимодействовать с репетиторами.</p>
                    <Button type="primary" onClick={() => setShowCreateForm(f => !f)} style={{ marginBottom: 16 }}>
                        {showCreateForm ? 'Скрыть форму' : 'Создать контракт'}
                    </Button>
                    {showCreateForm && (
                        <CreateContractForm clientId={clientId} onContractCreated={handleContractCreated} />
                    )}
                    <ContractFeed mode="client" id={clientId} key={feedKey} />
                </div>
            </Content>
        </Layout>
    );
};

export default ClientPage; 