import React, { useState } from 'react';
import { StarFilled } from '@ant-design/icons';
import { Review } from '../types/review';
import { getReviewById } from '../services/reviewService';
import { Modal } from 'antd';

interface ReviewListProps {
    reviews: Review[];
    loading?: boolean;
    error?: string;
}

export const ReviewList: React.FC<ReviewListProps> = ({ reviews, loading, error }) => {
    const [selectedReview, setSelectedReview] = useState<Review | null>(null);
    const [isModalVisible, setIsModalVisible] = useState(false);

    const handleViewReview = async (reviewId: number) => {
        try {
            const review = await getReviewById(reviewId);
            setSelectedReview(review);
            setIsModalVisible(true);
        } catch (error) {
            console.error('Error fetching review:', error);
        }
    };

    const handleCloseModal = () => {
        setIsModalVisible(false);
        setSelectedReview(null);
    };

    if (loading) {
        return (
            <div className="flex justify-center items-center py-8">
                <div className="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-blue-500"></div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="text-center py-4">
                <p className="text-red-500">{error}</p>
            </div>
        );
    }

    if (reviews.length === 0) {
        return (
            <div className="text-center py-8">
                <p className="text-gray-500">Пока нет отзывов</p>
            </div>
        );
    }

    return (
        <div className="space-y-4">
            {reviews.map((review) => (
                <div key={review.id} className="bg-white rounded-lg shadow p-4">
                    <div className="flex justify-between items-start mb-2">
                        <div className="flex items-center space-x-1 text-yellow-400">
                            {[...Array(5)].map((_, index) => (
                                <StarFilled
                                    key={index}
                                    className={index < review.rating ? 'text-yellow-400' : 'text-gray-300'}
                                />
                            ))}
                            <span className="ml-2 text-gray-600">
                                {review.rating} из 5
                            </span>
                        </div>
                        <div className="text-sm text-gray-500">
                            {new Date(review.created_at).toLocaleDateString('ru-RU', {
                                year: 'numeric',
                                month: 'long',
                                day: 'numeric'
                            })}
                        </div>
                        <button
                            onClick={() => handleViewReview(review.id)}
                            className="text-blue-500 hover:text-blue-700"
                        >
                            View Details
                        </button>
                    </div>
                    <p className="text-gray-700 mt-2">{review.comment}</p>
                </div>
            ))}

            <Modal
                title="Review Details"
                open={isModalVisible}
                onCancel={handleCloseModal}
                footer={null}
            >
                {selectedReview && (
                    <div className="space-y-4">
                        <p><strong>Rating:</strong> {selectedReview.rating}/5</p>
                        <p><strong>Comment:</strong> {selectedReview.comment}</p>
                        <p><strong>Создано:</strong> {new Date(selectedReview.created_at).toLocaleString()}</p>
                    </div>
                )}
            </Modal>
        </div>
    );
}; 