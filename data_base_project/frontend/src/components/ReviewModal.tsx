import React, { useState } from 'react';
import { Contract } from '../types/contract';
import { Review } from '../types/review';
import { CreateReviewRequest } from '../services/reviewService';

interface ReviewModalProps {
    isOpen: boolean;
    onClose: () => void;
    onSubmit: (review: CreateReviewRequest) => void;
    contract: Contract;
}

export const ReviewModal: React.FC<ReviewModalProps> = ({ isOpen, onClose, onSubmit, contract }) => {
    console.log('[ReviewModal] rendered, isOpen:', isOpen, 'contract.id:', contract.id);
    const [rating, setRating] = useState(5);
    const [comment, setComment] = useState('');
    const [isSubmitting, setIsSubmitting] = useState(false);

    const handleSubmit = async () => {
        try {
            setIsSubmitting(true);
            const reviewData: CreateReviewRequest = {
                client_id: Number(contract.client_id),
                repetitor_id: Number(contract.repetitor_id),
                contract_id: Number(contract.id),
                rating,
                comment
            };
            
            console.log('Submitting review data:', reviewData);
            await onSubmit(reviewData);
            setComment('');
            setRating(5);
            onClose();
        } catch (error) {
            console.error('Error submitting review:', error);
        } finally {
            setIsSubmitting(false);
        }
    };

    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white p-6 rounded-lg w-96">
                <h2 className="text-xl font-bold mb-4">Оставить отзыв</h2>
                <div className="mb-4">
                    <label className="block text-sm font-medium text-gray-700 mb-2">Оценка</label>
                    <select
                        value={rating}
                        onChange={(e) => setRating(Number(e.target.value))}
                        className="w-full p-2 border rounded"
                    >
                        {[1, 2, 3, 4, 5].map((value) => (
                            <option key={value} value={value}>{value}</option>
                        ))}
                    </select>
                </div>
                <div className="mb-4">
                    <label className="block text-sm font-medium text-gray-700 mb-2">Комментарий</label>
                    <textarea
                        value={comment}
                        onChange={(e) => setComment(e.target.value)}
                        className="w-full p-2 border rounded"
                        rows={4}
                        placeholder="Напишите ваш отзыв..."
                        required
                    />
                </div>
                <div className="flex justify-end gap-2">
                    <button
                        onClick={onClose}
                        className="px-4 py-2 bg-gray-200 rounded hover:bg-gray-300"
                        disabled={isSubmitting}
                    >
                        Отмена
                    </button>
                    <button
                        onClick={handleSubmit}
                        className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:opacity-50"
                        disabled={isSubmitting || !comment.trim()}
                    >
                        {isSubmitting ? 'Отправка...' : 'Отправить'}
                    </button>
                </div>
            </div>
        </div>
    );
}; 