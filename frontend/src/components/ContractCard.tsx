import React, { useState, useEffect } from 'react';
import { Contract } from '../types/contract';
import { ClockCircleOutlined, DollarOutlined, CalendarOutlined, CheckCircleOutlined, CloseCircleOutlined, WarningOutlined, UserOutlined, StarOutlined, MessageOutlined } from '@ant-design/icons';
import { API_ENDPOINTS } from '../config';
import { ReviewModal } from './ReviewModal';
import { ReviewList } from './ReviewList';
import { Review } from '../types/review';
import { createReview, createRepetitorReview, CreateReviewRequest } from '../services/reviewService';
import { payForContract } from '../services/reviewService';
import { useNavigate } from 'react-router-dom';
import { Form, Input, Button, message, InputNumber, List } from 'antd';
import { Lesson } from '../types/lesson';

interface ContractCardProps {
    contract: Contract;
    mode?: 'available' | 'client' | 'repetitor';
    onRespond?: () => void;
    repetitorId?: number;
    onSelect?: (contract: Contract) => void;
}

const getStatusColor = (status: number): string => {
    switch (status) {
        case 1: return 'bg-yellow-100 text-yellow-800'; // Pending
        case 2: return 'bg-green-100 text-green-800';   // Active
        case 3: return 'bg-blue-100 text-blue-800';     // Completed
        case 4: return 'bg-red-100 text-red-800';       // Cancelled
        case 5: return 'bg-gray-100 text-gray-800';     // Banned
        default: return 'bg-gray-100 text-gray-800';
    }
};

const getStatusText = (status: number): string => {
    switch (status) {
        case 1: return 'На рассмотрении';
        case 2: return 'Активный';
        case 3: return 'Завершен';
        case 4: return 'Отменен';
        case 5: return 'Заблокирован';
        default: return 'Неизвестно';
    }
};

const getStatusIcon = (status: number) => {
    switch (status) {
        case 1: return <WarningOutlined className="text-yellow-500" />;
        case 2: return <CheckCircleOutlined className="text-green-500" />;
        case 3: return <CheckCircleOutlined className="text-blue-500" />;
        case 4: return <CloseCircleOutlined className="text-red-500" />;
        case 5: return <CloseCircleOutlined className="text-gray-500" />;
        default: return null;
    }
};

const getPaymentStatusText = (status: number): string => {
    switch (status) {
        case 1: return 'Ожидает оплаты';
        case 2: return 'Оплачен';
        case 3: return 'Отказано';
        case 4: return 'Возвращен';
        default: return 'Не определен';
    }
};

const getContractTypeText = (category: number): string => {
    switch (category) {
        case 1: return 'Перевод';
        case 2: return 'Написание';
        case 3: return 'Дизайн';
        case 4: return 'Программирование';
        case 5: return 'Другое';
        default: return 'Не определен';
    }
};

const formatDate = (dateString: string): string => {
    const date = new Date(dateString);
    return date.toLocaleDateString('ru-RU', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });
};

const ContractCard: React.FC<ContractCardProps> = ({ contract, mode, onRespond, repetitorId, onSelect }) => {
    const navigate = useNavigate();
    const [isReviewModalOpen, setIsReviewModalOpen] = useState(false);
    const [loadingReviews, setLoadingReviews] = useState(false);
    const [reviewsError, setReviewsError] = useState<string | null>(null);
    const [showReviews, setShowReviews] = useState(false);
    const [currentReview, setCurrentReview] = useState<Review | null>(null);
    const [isPaying, setIsPaying] = useState(false);
    const [form] = Form.useForm();
    const [isAddingLesson, setIsAddingLesson] = useState(false);
    const [lessons, setLessons] = useState<Lesson[]>([]);
    const [loadingLessons, setLoadingLessons] = useState(false);
    const [lessonsOffset, setLessonsOffset] = useState(0);
    const [lessonsSize] = useState(5);

    useEffect(() => {
        fetchLessons();
    }, [contract.id, lessonsOffset]);

    const fetchLessons = async () => {
        if (!contract.id) return;
        
        setLoadingLessons(true);
        try {
            const response = await fetch(API_ENDPOINTS.REPETITOR.GET_LESSONS(contract.id, lessonsOffset, lessonsSize), {
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                }
            });
            if (!response.ok) throw new Error('Failed to fetch lessons');
            
            const data = await response.json();
            setLessons(data);
        } catch (error) {
            console.error('Error fetching lessons:', error);
            message.error('Не удалось загрузить уроки');
        } finally {
            setLoadingLessons(false);
        }
    };

    const fetchReview = async (reviewId: number) => {
        if (!reviewId) return;
        
        console.log('Starting to fetch single review with ID:', reviewId);
        setLoadingReviews(true);
        setReviewsError(null);
        try {
            const url = `${API_ENDPOINTS.CLIENT.GET_REVIEW}?review_id=${reviewId}`;
            console.log('Making request to URL:', url);
            
            const response = await fetch(url, {
                method: 'GET',
                headers: {
                    'Accept': 'application/json',
                },
                credentials: 'include',
            });

            console.log('Review response status:', response.status);
            const responseText = await response.text();
            console.log('Review response text:', responseText);

            if (!response.ok) {
                throw new Error(`Failed to fetch review: ${responseText}`);
            }

            const data = JSON.parse(responseText);
            console.log('Received review data:', data);
            setCurrentReview(data);
            setShowReviews(true);
        } catch (error) {
            console.error('Detailed error when fetching review:', {
                reviewId,
                error: error instanceof Error ? error.message : 'Unknown error',
                errorObject: error
            });
            setReviewsError('Не удалось загрузить отзыв');
        } finally {
            setLoadingReviews(false);
        }
    };

    const handleReviewSubmit = async (review: CreateReviewRequest) => {
        try {
            console.log('Submitting review:', review);
            if (mode === 'repetitor') {
                await createRepetitorReview(review);
            } else {
                await createReview(review);
            }
            setIsReviewModalOpen(false);
        } catch (error) {
            console.error('Error submitting review:', error);
        }
    };

    const handleRespond = async () => {
        if (!repetitorId) {
            console.error('Repetitor ID is required to accept a contract');
            return;
        }

        try {
            const url = API_ENDPOINTS.REPETITOR.ACCEPT_CONTRACT(repetitorId, contract.id);
            console.log('Making request to:', url);
            const response = await fetch(url, {
                method: 'POST',
                headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/json',
                },
                credentials: 'include'
            });

            if (!response.ok) {
                const errorText = await response.text();
                console.error('Response not OK:', response.status, errorText);
                throw new Error(`Failed to accept contract: ${errorText}`);
            }

            console.log('Contract accepted successfully');
            onRespond?.();
        } catch (error) {
            console.error('Error accepting contract:', error);
        }
    };

    const handlePay = async () => {
        if (!contract.id || !repetitorId) {
            console.error('Contract ID and repetitor ID are required to make a payment');
            return;
        }

        setIsPaying(true);
        try {
            await payForContract(contract.id, contract.price, repetitorId);
            onRespond?.();
        } catch (error) {
            console.error('Error paying for contract:', error);
        } finally {
            setIsPaying(false);
        }
    };

    const handleChatClick = async (e: React.MouseEvent) => {
        e.stopPropagation();
        
        try {
            const createChatResponse = await fetch(API_ENDPOINTS.CHAT.START_CR_CHAT(contract.client_id, contract.repetitor_id), {
                method: 'GET',
                headers: {
                    'Accept': 'application/json',
                },
                credentials: 'include',
            });

            if (!createChatResponse.ok) {
                throw new Error('Failed to create chat');
            }

            const chatId = await createChatResponse.json();
        
            const getChatResponse = await fetch(API_ENDPOINTS.CHAT.GET_CHAT(chatId), {
                method: 'GET',
                headers: {
                    'Accept': 'application/json',
                },
                credentials: 'include',
            });

            if (!getChatResponse.ok) {
                throw new Error('Failed to get chat details');
            }

            navigate(`/chat?id=${chatId}`);
        } catch (error) {
            console.error('Error creating/accessing chat:', error);
            alert('Не удалось создать или получить доступ к чату');
        }
    };

    const handleAddLesson = async (values: { duration: number }) => {
        setIsAddingLesson(true);
        try {
            const response = await fetch(API_ENDPOINTS.REPETITOR.ADD_LESSON, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                },
                body: JSON.stringify({
                    contract_id: contract.id,
                    duration: values.duration
                })
            });

            if (!response.ok) throw new Error('Failed to add lesson');
            
            message.success('Урок успешно добавлен');
            form.resetFields();
            onRespond?.(); // Обновляем список контрактов
        } catch (error) {
            console.error('Error adding lesson:', error);
            message.error('Не удалось добавить урок');
        } finally {
            setIsAddingLesson(false);
        }
    };

    const handleOpenCMChat = async (clientId: number) => {
        try {
            const moderatorId = Number(localStorage.getItem('user_id'));
            const res = await fetch(API_ENDPOINTS.CHAT.START_CM_CHAT(clientId, moderatorId), {
                method: 'GET',
                headers: { 'Accept': 'application/json' },
                credentials: 'include',
            });
            if (!res.ok) throw new Error('Ошибка создания чата');
            const chatId = await res.json();
            navigate(`/chat?id=${chatId}`);
        } catch (e) {
            alert('Не удалось открыть чат с клиентом');
        }
    };

    const handleOpenRMChat = async (repetitorId: number) => {
        try {
            const moderatorId = Number(localStorage.getItem('user_id'));
            const res = await fetch(API_ENDPOINTS.CHAT.START_RM_CHAT(repetitorId, moderatorId), {
                method: 'GET',
                headers: { 'Accept': 'application/json' },
                credentials: 'include',
            });
            if (!res.ok) throw new Error('Ошибка создания чата');
            const chatId = await res.json();
            navigate(`/chat?id=${chatId}`);
        } catch (e) {
            alert('Не удалось открыть чат с репетитором');
        }
    };

    const handleCompleteContract = async (repetitorId: number, contractId: number) => {
        try {
            const response = await fetch(API_ENDPOINTS.REPETITOR.COMPLETE_CONTRACT(repetitorId, contractId), {
                method: 'GET',
                credentials: 'include',
            });
            if (!response.ok) throw new Error('Ошибка при завершении контракта');
            onRespond?.();
        } catch (err) {
            alert('Ошибка при завершении контракта');
        }
    };

    return (
        <div 
            className="bg-white rounded-xl shadow-lg p-6 hover:shadow-xl transition-all duration-300 transform hover:-translate-y-1 border border-gray-100 cursor-pointer" 
            onClick={() => onSelect?.(contract)}
        >
            <div className="flex justify-between items-start mb-4">
                <div>
                    <h3 className="text-xl font-semibold text-gray-800">Контракт #{contract.id}</h3>
                    <p className="text-sm text-gray-500 mt-1">{getContractTypeText(contract.contract_category)}</p>
                </div>
                <span className={`px-3 py-1 rounded-full text-sm font-medium flex items-center gap-2 ${getStatusColor(contract.status)}`}>
                    {getStatusIcon(contract.status)}
                    {getStatusText(contract.status)}
                </span>
            </div>
            
            <p className="text-gray-600 mb-6 line-clamp-3">{contract.description}</p>
            
            <div className="space-y-3 text-sm">
                <div className="flex justify-between items-center">
                    <span className="text-gray-500 flex items-center gap-2">
                        <DollarOutlined />
                        Стоимость
                    </span>
                    <span className="font-semibold text-gray-800">{contract.price.toLocaleString()} ₽</span>
                </div>
                <div className="flex justify-between items-center">
                    <span className="text-gray-500 flex items-center gap-2">
                        <DollarOutlined />
                        Комиссия
                    </span>
                    <span className="font-semibold text-gray-800">{contract.commission.toLocaleString()} ₽</span>
                </div>
                <div className="flex justify-between items-center">
                    <span className="text-gray-500 flex items-center gap-2">
                        <ClockCircleOutlined />
                        Статус оплаты
                    </span>
                    <span className="font-semibold text-gray-800">{getPaymentStatusText(contract.payment_status)}</span>
                </div>
                <div className="flex justify-between items-center">
                    <span className="text-gray-500 flex items-center gap-2">
                        <CalendarOutlined />
                        Начало
                    </span>
                    <span className="font-semibold text-gray-800">{formatDate(contract.start_date)}</span>
                </div>
                <div className="flex justify-between items-center">
                    <span className="text-gray-500 flex items-center gap-2">
                        <CalendarOutlined />
                        Окончание
                    </span>
                    <span className="font-semibold text-gray-800">{formatDate(contract.end_date)}</span>
                </div>
                {typeof contract.repetitor_id !== 'undefined' && (
                    <div className="flex justify-between items-center">
                        <span className="text-gray-500 flex items-center gap-2">
                            <UserOutlined />
                            ID репетитора
                        </span>
                        <span className="font-semibold text-gray-800">{contract.repetitor_id}</span>
                    </div>
                )}
            </div>

            {mode === 'available' && (
                <div className="mt-6">
                    <button
                        onClick={handleRespond}
                        className="w-full bg-blue-600 text-white py-2 px-4 rounded-lg hover:bg-blue-700 transition-colors duration-200"
                    >
                        Откликнуться
                    </button>
                </div>
            )}

            {mode === 'client' && contract.status === 2 && (
                <div className="mt-6">
                    <button
                        onClick={() => setIsReviewModalOpen(true)}
                        disabled={!!contract.review_client_id}
                        className={`w-full bg-green-600 text-white py-2 px-4 rounded-lg flex items-center justify-center gap-2 transition-colors duration-200 ${
                            contract.review_client_id
                                ? 'opacity-50 cursor-not-allowed'
                                : 'hover:bg-green-700'
                        }`}
                    >
                        <StarOutlined />
                        Оставить отзыв
                    </button>
                </div>
            )}

            {(contract.status === 2 || contract.status === 3) && (
                <div className="mt-4 space-y-2">
                    <button
                        onClick={() => {
                            console.log('Client review button clicked. Review ID:', contract.review_client_id);
                            contract.review_client_id && fetchReview(contract.review_client_id);
                        }}
                        disabled={!contract.review_client_id}
                        className={`w-full text-sm font-medium flex items-center gap-1 justify-center py-2 border rounded-lg ${
                            contract.review_client_id 
                            ? 'text-blue-600 hover:text-blue-800 border-blue-200 hover:bg-blue-50' 
                            : 'text-gray-400 border-gray-200 cursor-not-allowed'
                        }`}
                    >
                        <StarOutlined />
                        {contract.review_client_id ? 'Показать отзыв клиента' : 'Отзыв клиента отсутствует'}
                    </button>
                    <button
                        onClick={() => {
                            console.log('Repetitor review button clicked. Review ID:', contract.review_repetitor_id);
                            contract.review_repetitor_id && fetchReview(contract.review_repetitor_id);
                        }}
                        disabled={!contract.review_repetitor_id}
                        className={`w-full text-sm font-medium flex items-center gap-1 justify-center py-2 border rounded-lg ${
                            contract.review_repetitor_id 
                            ? 'text-blue-600 hover:text-blue-800 border-blue-200 hover:bg-blue-50' 
                            : 'text-gray-400 border-gray-200 cursor-not-allowed'
                        }`}
                    >
                        <StarOutlined />
                        {contract.review_repetitor_id ? 'Показать отзыв репетитора' : 'Отзыв репетитора отсутствует'}
                    </button>
                </div>
            )}

            {mode === 'repetitor' && contract.status === 2 && contract.payment_status !== 2 && (
                <div className="mt-6">
                    <button
                        onClick={handlePay}
                        disabled={isPaying}
                        className="w-full bg-green-600 text-white py-2 px-4 rounded-lg hover:bg-green-700 transition-colors duration-200 flex items-center justify-center gap-2"
                    >
                        <DollarOutlined />
                        {isPaying ? 'Оплата...' : 'Оплатить контракт'}
                    </button>
                </div>
            )}

            {mode === 'repetitor' && contract.status === 2 && (
                <div className="mt-6 flex flex-col gap-2">
                    <button
                        onClick={async (e) => {
                            e.stopPropagation();
                            if (!repetitorId) return;
                            if (!window.confirm('Вы уверены, что хотите отменить этот контракт?')) return;
                            try {
                                const response = await fetch(API_ENDPOINTS.REPETITOR.CANCEL_CONTRACT(repetitorId, contract.id), {
                                    method: 'GET',
                                    credentials: 'include',
                                });
                                if (!response.ok) throw new Error('Ошибка при отмене контракта');
                                onRespond?.();
                            } catch (err) {
                                alert('Ошибка при отмене контракта');
                            }
                        }}
                        className="w-full bg-red-600 text-white py-2 px-4 rounded-lg hover:bg-red-700 transition-colors duration-200 flex items-center justify-center gap-2"
                    >
                        <CloseCircleOutlined />
                        Отменить контракт
                    </button>
                    <button
                        onClick={async (e) => {
                            e.stopPropagation();
                            if (!repetitorId) return;
                            if (!window.confirm('Вы уверены, что хотите завершить этот контракт?')) return;
                            try {
                                const url = `/api/repetitor/complete_contract?id=${repetitorId}&c_id=${contract.id}`;
                                const response = await fetch(url, {
                                    method: 'GET',
                                    credentials: 'include',
                                });
                                if (!response.ok) throw new Error('Ошибка при завершении контракта');
                                onRespond?.();
                            } catch (err) {
                                alert('Ошибка при завершении контракта');
                            }
                        }}
                        className="w-full bg-blue-600 text-white py-2 px-4 rounded-lg hover:bg-blue-700 transition-colors duration-200 flex items-center justify-center gap-2"
                    >
                        <CheckCircleOutlined />
                        Завершить контракт
                    </button>
                </div>
            )}

            {showReviews && currentReview && (
                <div className="mt-4">
                    <ReviewList
                        reviews={[currentReview]}
                        loading={loadingReviews}
                        error={reviewsError || undefined}
                    />
                </div>
            )}

            {mode === 'client' && (
            <ReviewModal
                isOpen={isReviewModalOpen}
                onClose={() => setIsReviewModalOpen(false)}
                onSubmit={handleReviewSubmit}
                contract={contract}
            />
            )}

            <div className="mt-4 flex justify-end">
                {contract.client_id !== 0 && contract.repetitor_id !== 0 && (
                    <button
                        onClick={handleChatClick}
                        className="flex items-center gap-2 px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
                    >
                        <MessageOutlined />
                        Перейти в чат
                    </button>
                )}
            </div>

            {contract.status === 2 && ( // Если контракт активен
                <div className="mt-4 pt-4 border-t border-gray-100" onClick={e => e.stopPropagation()}>
                    <Form
                        form={form}
                        onFinish={handleAddLesson}
                        layout="inline"
                        size="small"
                    >
                        <Form.Item
                            name="duration"
                            rules={[
                                { required: true, message: 'Введите длительность' },
                                { type: 'number', min: 1, message: 'Длительность должна быть положительным числом' }
                            ]}
                        >
                            <InputNumber
                                min={1}
                                placeholder="Длительность (мин)"
                                style={{ width: '120px' }}
                            />
                        </Form.Item>
                        <Form.Item>
                            <Button 
                                type="primary" 
                                htmlType="submit" 
                                loading={isAddingLesson}
                                onClick={e => e.stopPropagation()}
                            >
                                Добавить урок
                            </Button>
                        </Form.Item>
                    </Form>

                    <div className="mt-4">
                        <h4 className="text-lg font-semibold mb-2">История уроков</h4>
                        <List
                            loading={loadingLessons}
                            dataSource={lessons}
                            renderItem={(lesson, idx) => (
                                <List.Item style={{ borderBottom: '1px solid #f0f0f0', padding: '12px 0' }}>
                                    <div className="w-full">
                                        <div className="flex justify-between items-center">
                                            <span className="text-gray-600 font-medium">
                                                #{lessonsOffset + idx + 1}. Длительность: {lesson.duration} мин
                                            </span>
                                        </div>
                                        <div className="text-gray-400 text-xs mt-1" style={{ whiteSpace: 'pre-line' }}>
                                            {new Date(lesson.created_at).toLocaleString('ru-RU', {
                                                year: 'numeric',
                                                month: '2-digit',
                                                day: '2-digit',
                                                hour: '2-digit',
                                                minute: '2-digit',
                                                second: '2-digit',
                                            })}
                                        </div>
                                    </div>
                                </List.Item>
                            )}
                            pagination={{
                                current: Math.floor(lessonsOffset / lessonsSize) + 1,
                                pageSize: lessonsSize,
                                onChange: (page) => setLessonsOffset((page - 1) * lessonsSize),
                                showSizeChanger: false,
                            }}
                        />
                    </div>
                </div>
            )}
        </div>
    );
};

export default ContractCard; 