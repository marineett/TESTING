import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { API_ENDPOINTS } from '../config';

interface Contract {
    id: number;
    title: string;
    description: string;
    status: string;
    created_at: string;
    contract_type?: string;
    owner_name?: string;
    owner_email?: string;
    total_amount?: number;
    currency?: string;
    last_modified?: string;
    version?: string;
    client_id?: number;
    repetitor_id?: number;
}

const BATCH_SIZE = 10;

const ContractList: React.FC = () => {
    const [contracts, setContracts] = useState<Contract[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [currentPage, setCurrentPage] = useState(0);
    const [hasMore, setHasMore] = useState(true);
    const navigate = useNavigate();

    const fetchContracts = async (from: number) => {
        try {
            const response = await fetch(API_ENDPOINTS.MODERATOR.GET_CONTRACTS(from, BATCH_SIZE), {
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                }
            });
            
            if (!response.ok) {
                throw new Error('Не удалось загрузить контракты');
            }

            const data = await response.json();
            setContracts(data);
            setHasMore(data.length === BATCH_SIZE);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Произошла ошибка');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        setLoading(true);
        fetchContracts(currentPage * BATCH_SIZE);
    }, [currentPage]);

    const handlePreviousPage = () => {
        if (currentPage > 0) {
            setCurrentPage(prev => prev - 1);
        }
    };

    const handleNextPage = () => {
        if (hasMore) {
            setCurrentPage(prev => prev + 1);
        }
    };

    const getStatusColor = (status: any) => {
        const statusStr = String(status).toLowerCase();
        switch (statusStr) {
            case 'active':
                return 'bg-green-100 text-green-800';
            case 'pending':
                return 'bg-yellow-100 text-yellow-800';
            case 'expired':
                return 'bg-red-100 text-red-800';
            case 'draft':
                return 'bg-gray-100 text-gray-800';
            default:
                return 'bg-blue-100 text-blue-800';
        }
    };

    const handleBanContract = async (contractId: number) => {
        if (!window.confirm('Вы уверены, что хотите забанить этот контракт?')) return;
        try {
            const response = await fetch(API_ENDPOINTS.MODERATOR.BAN_CONTRACT(contractId), {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                }
            });
            if (!response.ok) {
                throw new Error('Ошибка при бане контракта');
            }
            fetchContracts(currentPage * BATCH_SIZE);
        } catch (err) {
            alert('Ошибка при бане контракта');
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

    if (loading) {
        return (
            <div className="flex justify-center items-center py-8">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="bg-red-50 border-l-4 border-red-400 p-4">
                <div className="flex">
                    <div className="flex-shrink-0">
                        <svg className="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
                            <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
                        </svg>
                    </div>
                    <div className="ml-3">
                        <p className="text-sm text-red-700">{error}</p>
                    </div>
                </div>
            </div>
        );
    }

    return (
        <div>
            <div className="bg-white shadow overflow-hidden sm:rounded-md">
                <ul className="divide-y divide-gray-200">
                    {contracts.map((contract) => (
                        <li key={contract.id} className="px-6 py-6 hover:bg-gray-50 transition-colors duration-150">
                            <div className="flex flex-col space-y-4">
                                <div className="flex items-center justify-between">
                                    <div className="flex items-center space-x-3">
                                        <h3 className="text-xl font-semibold text-gray-900">{contract.title}</h3>
                                        <span className={`px-3 py-1 inline-flex text-sm leading-5 font-semibold rounded-full ${getStatusColor(contract.status)}`}>
                                            {contract.status}
                                        </span>
                                    </div>
                                    <div className="text-sm text-gray-500 flex items-center gap-2">
                                        ID: #{contract.id}
                                        {/* Кнопка бан */}
                                        {([1, 2].includes(Number(contract.status))) && (
                                            <button
                                                onClick={() => handleBanContract(contract.id)}
                                                className="ml-2 px-3 py-1 bg-red-600 text-white text-xs rounded hover:bg-red-700 transition-colors"
                                            >
                                                Забанить контракт
                                            </button>
                                        )}
                                        {/* Кнопки чатов для модератора */}
                                        <>
                                            <button
                                                onClick={async () => {
                                                    const moderatorId = Number(localStorage.getItem('user_id'));
                                                    if (contract.client_id && moderatorId) {
                                                        try {
                                                            const res = await fetch(API_ENDPOINTS.CHAT.START_CM_CHAT(contract.client_id, moderatorId), {
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
                                                    } else {
                                                        alert('client_id или m_id не найден');
                                                    }
                                                }}
                                                className="ml-2 px-3 py-1 bg-blue-600 text-white text-xs rounded hover:bg-blue-700 transition-colors"
                                            >
                                                Чат с клиентом
                                            </button>
                                            <button
                                                onClick={async () => {
                                                    const moderatorId = Number(localStorage.getItem('user_id'));
                                                    if (contract.repetitor_id && moderatorId) {
                                                        try {
                                                            const res = await fetch(API_ENDPOINTS.CHAT.START_RM_CHAT(contract.repetitor_id, moderatorId), {
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
                                                    } else {
                                                        alert('repetitor_id или m_id не найден');
                                                    }
                                                }}
                                                className="ml-2 px-3 py-1 bg-green-600 text-white text-xs rounded hover:bg-green-700 transition-colors"
                                            >
                                                Чат с репетитором
                                            </button>
                                        </>
                                    </div>
                                </div>

                                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                    <div className="space-y-2">
                                        <p className="text-gray-700">{contract.description}</p>
                                        {contract.contract_type && (
                                            <div className="flex items-center text-sm text-gray-600">
                                                <svg className="h-4 w-4 mr-1.5" width={16} height={16} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                                                </svg>
                                                Type: {contract.contract_type}
                                            </div>
                                        )}
                                    </div>

                                    {(contract.owner_name || contract.owner_email) && (
                                        <div className="space-y-2">
                                            <h4 className="text-sm font-medium text-gray-900">Owner Information</h4>
                                            {contract.owner_name && (
                                                <p className="text-sm text-gray-600">Name: {contract.owner_name}</p>
                                            )}
                                            {contract.owner_email && (
                                                <p className="text-sm text-gray-600">Email: {contract.owner_email}</p>
                                            )}
                                        </div>
                                    )}
                                </div>

                                {contract.total_amount && (
                                    <div className="flex items-center text-sm text-gray-600">
                                        <svg className="h-4 w-4 mr-1.5" width={16} height={16} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                                        </svg>
                                        Amount: {contract.total_amount} {contract.currency || 'USD'}
                                    </div>
                                )}

                                <div className="flex flex-wrap gap-4 text-sm text-gray-500">
                                    <div className="flex items-center">
                                        <svg className="h-4 w-4 mr-1.5" width={16} height={16} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                                        </svg>
                                        Создано: {new Date(contract.created_at).toLocaleDateString()}
                                    </div>
                                    {contract.last_modified && (
                                        <div className="flex items-center">
                                            <svg className="h-4 w-4 mr-1.5" width={16} height={16} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                                            </svg>
                                            Modified: {new Date(contract.last_modified).toLocaleDateString()}
                                        </div>
                                    )}
                                    {contract.version && (
                                        <div className="flex items-center">
                                            <svg className="h-4 w-4 mr-1.5" width={16} height={16} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
                                            </svg>
                                            Version: {contract.version}
                                        </div>
                                    )}
                                </div>
                            </div>
                        </li>
                    ))}
                </ul>
            </div>

            {/* Пагинация */}
            <div style={{ background: '#232b36', color: '#fff', fontWeight: 500 }} className="px-4 py-3 flex items-center justify-center border-t border-gray-200 sm:px-6 mt-4">
                <nav className="relative z-0 inline-flex rounded-md shadow-sm -space-x-px" aria-label="Pagination">
                    <button
                        onClick={handlePreviousPage}
                        disabled={currentPage === 0}
                        style={{ background: '#232b36', color: '#fff', borderColor: '#232b36', fontWeight: 500, ...(currentPage === 0 ? { opacity: 0.5, cursor: 'not-allowed' } : {}) }}
                        className={`relative inline-flex items-center px-4 py-2 rounded-l-md border text-sm font-medium`}
                    >
                        Назад
                    </button>
                    <span style={{ color: '#fff', fontWeight: 500 }} className="mx-2">Страница {currentPage + 1}</span>
                    <button
                        onClick={handleNextPage}
                        disabled={!hasMore}
                        style={{ background: '#232b36', color: '#fff', borderColor: '#232b36', fontWeight: 500, ...(!hasMore ? { opacity: 0.5, cursor: 'not-allowed' } : {}) }}
                        className={`relative inline-flex items-center px-4 py-2 rounded-r-md border text-sm font-medium`}
                    >
                        Вперед
                    </button>
                </nav>
            </div>
        </div>
    );
};

export default ContractList; 