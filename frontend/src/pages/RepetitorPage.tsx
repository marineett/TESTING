import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import ContractFeed from '../components/ContractFeed';
import { ReviewModal } from '../components/ReviewModal';
import { Contract } from '../types/contract';
import { createRepetitorReview, CreateReviewRequest } from '../services/reviewService';
import { Button, Layout, Typography } from 'antd';
import ChatListWindow from '../components/ChatListWindow';
import { API_ENDPOINTS } from '../config';

const { Header, Content } = Layout;
const { Title } = Typography;

const REVIEWS_PER_BATCH = 5;

const RepetitorPage: React.FC = () => {
    const navigate = useNavigate();
    const [repetitorId] = React.useState<number>(Number(localStorage.getItem('repetitor_id')) || 1);
    const [isReviewModalOpen, setIsReviewModalOpen] = useState(false);
    const [selectedContract, setSelectedContract] = useState<Contract | null>(null);
    const [profile, setProfile] = useState<any>(null);
    const [reviewBatch, setReviewBatch] = useState(1);

    useEffect(() => {
        const fetchProfile = async () => {
            if (!repetitorId) return;
            const offset = Math.max(0, (reviewBatch - 1) * REVIEWS_PER_BATCH);
            const response = await fetch(API_ENDPOINTS.REPETITOR.GET_PROFILE(repetitorId, offset, REVIEWS_PER_BATCH));
            if (response.ok) {
                setProfile(await response.json());
            }
        };
        fetchProfile();
    }, [repetitorId, reviewBatch]);

    const handleShowMoreReviews = () => {
        setReviewBatch((prev) => prev + 1);
    };

    const handleLogout = () => {
        localStorage.removeItem('userType');
        localStorage.removeItem('user_id');
        localStorage.removeItem('repetitor_id');
        navigate('/login');
    };

    const handleReviewSubmit = async (review: CreateReviewRequest) => {
        try {
            console.log('Submitting repetitor review:', review);
            await createRepetitorReview(review);
            setIsReviewModalOpen(false);
            setSelectedContract(null);
        } catch (error) {
            console.error('Error submitting repetitor review:', error);
        }
    };

    const handleContractSelect = (contract: Contract) => {
        setSelectedContract(contract);
        setIsReviewModalOpen(true);
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
                <Title level={3} style={{ color: 'white', margin: 0 }}>Кабинет репетитора</Title>
                <Button type="primary" onClick={handleLogout}>Выйти</Button>
            </Header>
            <Content style={{ padding: '24px', marginTop: 64 }}>
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                    <div style={{ padding: 24, background: '#fff', minHeight: 360, borderRadius: 10 }}>
                        <ChatListWindow userId={repetitorId} role="repetitor" />
                        {profile && (
                            <div style={{ marginBottom: 24, padding: 16, background: '#f5f5f5', borderRadius: 8 }}>
                                <div style={{ fontSize: 20, fontWeight: 600, marginBottom: 8 }}>
                                    {profile.last_name} {profile.first_name} {profile.middle_name}
                                </div>
                                <div style={{ color: '#888', marginBottom: 4 }}>Email: {profile.email}</div>
                                <div style={{ color: '#888', marginBottom: 4 }}>Телефон: {profile.telephone_number}</div>
                                <div style={{ color: '#888', marginBottom: 8 }}>Средний рейтинг: {profile.mean_rating?.toFixed(2) ?? '—'}</div>
                                <div style={{ fontWeight: 500, marginBottom: 4 }}>Резюме:</div>
                                <div style={{ color: '#555', marginBottom: 8 }}>{profile.resume_title}</div>
                                <div style={{ color: '#555', marginBottom: 8 }}>{profile.resume_description}</div>
                                {profile.resume_prices && Object.keys(profile.resume_prices).length > 0 && (
                                    <div style={{ marginBottom: 8 }}>
                                        <div style={{ fontWeight: 500, marginBottom: 4 }}>Цены:</div>
                                        <ul style={{ color: '#555', marginBottom: 8 }}>
                                            {Object.entries(profile.resume_prices).map(([service, price]) => (
                                                <li key={service}>{service}: {Number(price)} ₽</li>
                                            ))}
                                        </ul>
                                    </div>
                                )}
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
                    </div>
                    <div className="mb-8">
                        <h2 className="text-2xl font-bold mb-6 text-gray-800">Доступные контракты</h2>
                        <ContractFeed mode="available" id={repetitorId} />
                    </div>
                    <h2 className="text-2xl font-bold mb-4 text-gray-800">Мои заказы</h2>
                    <ContractFeed 
                        mode="repetitor" 
                        id={repetitorId} 
                        onContractSelect={handleContractSelect}
                    />
                </div>
            </Content>
        </Layout>
    );
};

export default RepetitorPage; 