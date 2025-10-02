import { API_ENDPOINTS } from '../config';

export interface Transaction {
    id: number;
    amount: number;
    status: number;
    created_at: string;
    updated_at: string;
    contract_id: number;
    user_id: number;
}

export const getTransactionsToApprove = async (): Promise<Transaction[]> => {
    try {
        const url = API_ENDPOINTS.MODERATOR.GET_TRANSACTION_TO_APPROVE;
        console.log('Fetching transactions to approve from:', url);

        const response = await fetch(url, {
            method: 'GET',
            headers: {
                'Accept': 'application/json',
            },
            credentials: 'include',
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`HTTP error! status: ${response.status}, message: ${errorText}`);
        }

        return await response.json();
    } catch (error) {
        console.error('Error fetching transactions to approve:', error);
        throw error;
    }
};

export const approveTransaction = async (transactionId: number): Promise<void> => {
    try {
        const url = `/api/moderator/approve_transaction?id=${transactionId}`;
        console.log('Approving transaction:', url);

        const response = await fetch(url, {
            method: 'GET',
            credentials: 'include',
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`HTTP error! status: ${response.status}, message: ${errorText}`);
        }
    } catch (error) {
        console.error('Error approving transaction:', error);
        throw error;
    }
};

export const rejectTransaction = async (transactionId: number): Promise<void> => {
    try {
        const url = API_ENDPOINTS.MODERATOR.REJECT_TRANSACTION(transactionId);
        console.log('Rejecting transaction:', url);

        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Accept': 'application/json',
            },
            credentials: 'include',
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`HTTP error! status: ${response.status}, message: ${errorText}`);
        }
    } catch (error) {
        console.error('Error rejecting transaction:', error);
        throw error;
    }
}; 