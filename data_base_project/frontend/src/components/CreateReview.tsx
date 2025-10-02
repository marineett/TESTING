import React, { useState } from 'react';
import { createReview, CreateReviewRequest } from '../services/reviewService';

interface CreateReviewProps {
    clientId: number;
    repetitorId: number;
    contractId: number;
    onSuccess?: () => void;
    onError?: (error: Error) => void;
}

export const CreateReview: React.FC<CreateReviewProps> = ({ 
    clientId, 
    repetitorId,
    contractId,
    onSuccess,
    onError 
}) => {
    const [rating, setRating] = useState<number>(0);
    const [comment, setComment] = useState<string>('');
    const [isSubmitting, setIsSubmitting] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsSubmitting(true);

        try {
            await createReview({
                client_id: clientId,
                repetitor_id: repetitorId,
                contract_id: contractId,
                rating,
                comment
            });
            setRating(0);
            setComment('');
            onSuccess?.();
        } catch (error) {
            onError?.(error as Error);
        } finally {
            setIsSubmitting(false);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-4">
            <div>
                <label className="block text-sm font-medium text-gray-700">Rating</label>
                <div className="flex space-x-2">
                    {[1, 2, 3, 4, 5].map((star) => (
                        <button
                            key={star}
                            type="button"
                            onClick={() => setRating(star)}
                            className={`text-2xl ${rating >= star ? 'text-yellow-400' : 'text-gray-300'}`}
                        >
                            ★
                        </button>
                    ))}
                </div>
            </div>
            <div>
                <label className="block text-sm font-medium text-gray-700">Comment</label>
                <textarea
                    value={comment}
                    onChange={(e) => setComment(e.target.value)}
                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
                    rows={3}
                />
            </div>
            <button
                type="submit"
                disabled={isSubmitting || rating === 0 || !comment.trim()}
                className="inline-flex justify-center rounded-md border border-transparent bg-indigo-600 py-2 px-4 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 disabled:opacity-50"
            >
                {isSubmitting ? 'Submitting...' : 'Submit Review'}
            </button>
        </form>
    );
}; 