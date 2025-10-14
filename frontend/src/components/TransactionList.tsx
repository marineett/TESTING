import React, { useState, useEffect } from 'react';
import { Transaction } from '../services/transactionService';
import { getTransactionsToApprove, approveTransaction, rejectTransaction } from '../services/transactionService';
import { CheckCircleOutlined, CloseCircleOutlined, DollarOutlined, LoadingOutlined } from '@ant-design/icons';

const TransactionList: React.FC = () => {
    const [transactions, setTransactions] = useState<Transaction[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [processingTransaction, setProcessingTransaction] = useState<number | null>(null);

    const fetchTransactions = async () => {
        try {
            setLoading(true);
            setError(null);
            const data = await getTransactionsToApprove();
            let txs: Transaction[] = [];
            if (Array.isArray(data)) {
                txs = data;
            } else if (data && typeof data === 'object') {
                txs = [data];
            }
            setTransactions(txs);
        } catch (err) {
            setError('Не удалось загрузить транзакции');
            console.error('Error fetching transactions:', err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchTransactions();
    }, []);

    const handleApprove = async (transactionId: number) => {
        try {
            setProcessingTransaction(transactionId);
            await approveTransaction(transactionId);
            await fetchTransactions();
        } catch (err) {
            console.error('Error approving transaction:', err);
        } finally {
            setProcessingTransaction(null);
        }
    };

    const handleReject = async (transactionId: number) => {
        try {
            setProcessingTransaction(transactionId);
            await rejectTransaction(transactionId);
            await fetchTransactions();
        } catch (err) {
            console.error('Error rejecting transaction:', err);
        } finally {
            setProcessingTransaction(null);
        }
    };

    if (loading) {
        return (
            <div className="flex justify-center items-center h-64">
                <LoadingOutlined className="text-4xl text-blue-500" />
            </div>
        );
    }

    if (error) {
        return (
            <div className="text-red-500 text-center p-4">
                {error}
            </div>
        );
    }

    if (transactions.length === 0) {
        return (
            <div className="text-gray-500 text-center p-4">
                Нет транзакций для подтверждения
            </div>
        );
    }

    return (
        <div className="space-y-4">
            {transactions.map((transaction) => (
                <div key={transaction.id} className="bg-white rounded-lg shadow-md p-6">
                    <div className="flex justify-between items-start mb-4">
                        <div>
                            <h3 className="text-lg font-semibold">Транзакция #{transaction.id}</h3>
                            <p className="text-sm text-gray-500">
                                Создана: {new Date(transaction.created_at).toLocaleString()}
                            </p>
                        </div>
                        <div className="flex items-center gap-2">
                            <span className="text-gray-500 flex items-center gap-1">
                                <DollarOutlined />
                                {transaction.amount} ₽
                            </span>
                        </div>
                    </div>

                    <div className="grid grid-cols-2 gap-4 mb-4">
                        <div>
                            <p className="text-sm text-gray-500">ID контракта</p>
                            <p className="font-medium">{transaction.contract_id}</p>
                        </div>
                        <div>
                            <p className="text-sm text-gray-500">ID пользователя</p>
                            <p className="font-medium">{transaction.user_id}</p>
                        </div>
                    </div>

                    <div className="flex justify-end gap-2">
                        <button
                            onClick={() => handleReject(transaction.id)}
                            disabled={processingTransaction === transaction.id}
                            className="px-4 py-2 bg-red-100 text-red-600 rounded-lg hover:bg-red-200 transition-colors flex items-center gap-2"
                        >
                            {processingTransaction === transaction.id ? (
                                <LoadingOutlined />
                            ) : (
                                <CloseCircleOutlined />
                            )}
                            Отклонить
                        </button>
                        <button
                            onClick={() => handleApprove(transaction.id)}
                            disabled={processingTransaction === transaction.id}
                            className="px-4 py-2 bg-green-100 text-green-600 rounded-lg hover:bg-green-200 transition-colors flex items-center gap-2"
                        >
                            {processingTransaction === transaction.id ? (
                                <LoadingOutlined />
                            ) : (
                                <CheckCircleOutlined />
                            )}
                            Одобрить
                        </button>
                    </div>
                </div>
            ))}
        </div>
    );
};

export default TransactionList; 