import { API_ENDPOINTS } from '../config';
import { Review, CreateReviewData } from '../types/review';

export interface CreateReviewRequest extends CreateReviewData {
    contract_id: number;
}

export const createReview = async (review: CreateReviewRequest): Promise<Review> => {
    try {
        const { contract_id, ...reviewData } = review;
        const url = `${API_ENDPOINTS.CLIENT.MAKE_REVIEW}?contract_id=${contract_id}`;
        
        console.log('Making review request to:', url);
        console.log('Review data:', reviewData);

        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            },
            credentials: 'include',
            body: JSON.stringify(reviewData),
        });

        console.log('Review response status:', response.status);
        const responseText = await response.text();
        console.log('Review response text:', responseText);

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}, message: ${responseText}`);
        }

        return JSON.parse(responseText);
    } catch (error) {
        console.error('Error creating review:', error);
        throw error;
    }
};

export const createRepetitorReview = async (review: CreateReviewRequest): Promise<Review> => {
    try {
        const { contract_id, ...reviewData } = review;
        const url = `${API_ENDPOINTS.REPETITOR.MAKE_REVIEW}?contract_id=${contract_id}`;
        
        console.log('Making repetitor review request to:', url);
        console.log('Review data:', reviewData);

        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            },
            credentials: 'include',
            body: JSON.stringify(reviewData),
        });

        console.log('Review response status:', response.status);
        const responseText = await response.text();
        console.log('Review response text:', responseText);

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}, message: ${responseText}`);
        }

        return JSON.parse(responseText);
    } catch (error) {
        console.error('Error creating repetitor review:', error);
        throw error;
    }
};

export const getRepetitorReviews = async (reviewId: number): Promise<Review[]> => {
    try {
        const url = API_ENDPOINTS.CLIENT.GET_REVIEWS(reviewId);
        console.log('Fetching reviews from URL:', url);

        const response = await fetch(url, {
            method: 'GET',
            headers: {
                'Accept': 'application/json',
            },
            credentials: 'include',
        });

        console.log('Reviews response status:', response.status);
        console.log('Reviews response headers:', Object.fromEntries(response.headers.entries()));

        if (!response.ok) {
            const errorText = await response.text();
            console.error('Error response text:', errorText);
            throw new Error(`HTTP error! status: ${response.status}, message: ${errorText}`);
        }

        const data = await response.json();
        console.log('Received reviews data:', data);
        return data;
    } catch (error) {
        console.error('Error fetching reviews:', {
            reviewId,
            error: error instanceof Error ? error.message : 'Unknown error',
            errorObject: error
        });
        throw error;
    }
};

export const getReviewById = async (reviewId: number): Promise<Review> => {
    try {
        const url = `${API_ENDPOINTS.CLIENT.GET_REVIEW}?review_id=${reviewId}`;
        console.log('Fetching review from URL:', url);

        const response = await fetch(url, {
            method: 'GET',
            headers: {
                'Accept': 'application/json',
            },
            credentials: 'include',
        });

        console.log('Review response status:', response.status);
        console.log('Review response headers:', Object.fromEntries(response.headers.entries()));
        
        const responseText = await response.text();
        console.log('Review response text:', responseText);

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}, message: ${responseText}`);
        }

        const parsedResponse = JSON.parse(responseText);
        console.log('Parsed review response:', parsedResponse);
        return parsedResponse;
    } catch (error) {
        console.error('Detailed error in getReviewById:', {
            reviewId,
            error: error instanceof Error ? error.message : 'Unknown error',
            errorObject: error
        });
        throw error;
    }
}; 

export const payForContract = async (contractId: number, amount: number, userId: number): Promise<void> => {
    try {
        const url = `/api/repetitor/pay_for_contract?contract_id=${contractId}&amount=${amount}&user_id=${userId}`;
        
        console.log('Making payment request to:', url);

        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Accept': 'application/json',
            },
            credentials: 'include',
        });

        console.log('Payment response status:', response.status);
        const responseText = await response.text();
        console.log('Payment response text:', responseText);

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}, message: ${responseText}`);
        }
    } catch (error) {
        console.error('Error paying for contract:', error);
        throw error;
    }
}; 