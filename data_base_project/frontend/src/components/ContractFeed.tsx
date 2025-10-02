import React, { useEffect, useState } from 'react';
import { API_ENDPOINTS } from '../config';
import { Contract } from '../types/contract';
import ContractCard from './ContractCard';

interface ContractFeedProps {
    mode: 'client' | 'repetitor' | 'repetitor-status' | 'available';
    id?: number;
    onContractSelect?: (contract: Contract) => void;
}

const STATUS_TABS = [
    { key: 'pending', label: 'На рассмотрении' },
    { key: 'active', label: 'Активные' },
    { key: 'completed', label: 'Завершённые' },
    { key: 'cancelled', label: 'Отменённые' },
    { key: 'banned', label: 'Заблокированные' },
] as const;
type StatusTabKey = typeof STATUS_TABS[number]['key'];

const CONTRACT_CATEGORIES = [
    { value: 1, label: 'Перевод' },
    { value: 2, label: 'Написание' },
    { value: 3, label: 'Дизайн' },
    { value: 4, label: 'Программирование' },
    { value: 5, label: 'Другое' },
];
const CONTRACT_SUBCATEGORIES = [
    { value: 2, label: 'Перевод' },
    { value: 3, label: 'Написание' },
    { value: 4, label: 'Дизайн' },
    { value: 5, label: 'Программирование' },
    { value: 6, label: 'Другое' },
];

const ContractFeed: React.FC<ContractFeedProps> = ({ mode, id, onContractSelect }) => {
    const [contracts, setContracts] = useState<{
        pending: Contract[];
        active: Contract[];
        completed: Contract[];
        cancelled: Contract[];
        banned: Contract[];
    }>({
        pending: [],
        active: [],
        completed: [],
        cancelled: [],
        banned: []
    });
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [activeStatus, setActiveStatus] = useState<StatusTabKey>('pending');

    const fetchContractsClient = async (status: number) => {
        if (typeof id !== 'number') throw new Error('id is required for client mode');
        try {
            const url = API_ENDPOINTS.CLIENT.CONTRACTS(id, 0, 100, status);
            console.log('Fetching contracts from:', url);
            
            const response = await fetch(url, {
                method: 'GET',
                headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/json',
                },
                credentials: 'include'
            });

            if (!response.ok) {
                const errorText = await response.text();
                console.error('Error response:', {
                    status: response.status,
                    statusText: response.statusText,
                    headers: Object.fromEntries(response.headers.entries()),
                    body: errorText
                });
                throw new Error(`HTTP error! status: ${response.status}, message: ${errorText}`);
            }

            const data = await response.json();
            console.log('Received contracts:', data);
            return data;
        } catch (err) {
            console.error('Error fetching contracts:', err);
            throw err;
        }
    };

    const fetchContractsRepetitor = async (status: number) => {
        if (typeof id !== 'number') throw new Error('id is required for repetitor mode');
        try {
            const url = API_ENDPOINTS.REPETITOR.GET_CONTRACTS(id, 0, 100, status);
            console.log('Fetching repetitor contracts from:', url);
            const response = await fetch(url, {
                method: 'GET',
                headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/json',
                },
                credentials: 'include'
            });
            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(`HTTP error! status: ${response.status}, message: ${errorText}`);
            }
            const data = await response.json();
            return data;
        } catch (err) {
            console.error('Error fetching repetitor contracts:', err);
            throw err;
        }
    };

    const fetchContractsRepetitorStatus = async (status: number) => {
        if (typeof id !== 'number') throw new Error('id is required for repetitor-status mode');
        try {
            const url = API_ENDPOINTS.REPETITOR.GET_CONTRACTS(id, 0, 100, status);
            console.log('Fetching repetitor contracts by status from:', url);
            const response = await fetch(url, {
                method: 'GET',
                headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/json',
                },
                credentials: 'include'
            });
            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(`HTTP error! status: ${response.status}, message: ${errorText}`);
            }
            const data = await response.json();
            return data;
        } catch (err) {
            console.error('Error fetching repetitor contracts by status:', err);
            throw err;
        }
    };

    const fetchAvailableContracts = async (status: number) => {
        try {
            const url = API_ENDPOINTS.REPETITOR.GET_AVAILABLE_CONTRACTS(0, 100, status);
            console.log('Fetching available contracts from:', url);
            const response = await fetch(url, {
                method: 'GET',
                headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/json',
                },
                credentials: 'include'
            });
            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(`HTTP error! status: ${response.status}, message: ${errorText}`);
            }
            const data = await response.json();
            return data;
        } catch (err) {
            console.error('Error fetching available contracts:', err);
            throw err;
        }
    };

    const loadContracts = async () => {
        try {
            setLoading(true);
            setError(null);
            if (mode === 'available') {
                const statusMap: Record<StatusTabKey, number> = {
                    pending: 1,
                    active: 2,
                    completed: 3,
                    cancelled: 4,
                    banned: 5
                };
                const allContracts = await fetchAvailableContracts(statusMap[activeStatus]).catch(() => []);
                setContracts({
                    pending: activeStatus === 'pending' ? allContracts : [],
                    active: activeStatus === 'active' ? allContracts : [],
                    completed: activeStatus === 'completed' ? allContracts : [],
                    cancelled: activeStatus === 'cancelled' ? allContracts : [],
                    banned: activeStatus === 'banned' ? allContracts : [],
                });
            } else if (mode === 'client') {
                const [pending, active, completed, cancelled, banned] = await Promise.all([
                    fetchContractsClient(1).catch(() => []),
                    fetchContractsClient(2).catch(() => []),
                    fetchContractsClient(3).catch(() => []),
                    fetchContractsClient(4).catch(() => []),
                    fetchContractsClient(5).catch(() => [])
                ]);
                setContracts({ pending, active, completed, cancelled, banned });
            } else if (mode === 'repetitor-status') {
                const [pending, active, completed, cancelled, banned] = await Promise.all([
                    fetchContractsRepetitorStatus(1).catch(() => []),
                    fetchContractsRepetitorStatus(2).catch(() => []),
                    fetchContractsRepetitorStatus(3).catch(() => []),
                    fetchContractsRepetitorStatus(4).catch(() => []),
                    fetchContractsRepetitorStatus(5).catch(() => [])
                ]);
                setContracts({ pending, active, completed, cancelled, banned });
            } else {
                const [pending, active, completed, cancelled, banned] = await Promise.all([
                    fetchContractsRepetitor(1).catch(() => []),
                    fetchContractsRepetitor(2).catch(() => []),
                    fetchContractsRepetitor(3).catch(() => []),
                    fetchContractsRepetitor(4).catch(() => []),
                    fetchContractsRepetitor(5).catch(() => [])
                ]);
                setContracts({ pending, active, completed, cancelled, banned });
            }
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Произошла ошибка при загрузке контрактов');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        loadContracts();
    }, [id, mode, activeStatus]);

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
        return new Date(dateString).toLocaleDateString('ru-RU', {
            year: 'numeric',
            month: 'long',
            day: 'numeric'
        });
    };

    const sortContractsByDate = (contracts: Contract[]): Contract[] => {
        return [...contracts].sort((a, b) => 
            new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
        );
    };

    const renderContractSection = (title: string, contracts: Contract[]) => {
        if (contracts.length === 0) return null;
        return (
            <div className="mb-8">
                <h2 className="text-xl font-semibold mb-4">{title}</h2>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {contracts.map((contract) => (
                        <ContractCard
                            key={contract.id}
                            contract={contract}
                            mode={mode === 'repetitor-status' ? 'repetitor' : mode}
                            onRespond={loadContracts}
                            repetitorId={id}
                            onSelect={onContractSelect}
                        />
                    ))}
                </div>
            </div>
        );
    };

    const renderStatusTabs = () => (
        <div className="mb-6">
            <div className="flex space-x-2 overflow-x-auto pb-2">
                {STATUS_TABS.map(tab => (
                    <button
                        key={tab.key}
                        onClick={() => setActiveStatus(tab.key)}
                        className={`px-4 py-2 rounded-lg font-medium transition-colors ${
                            activeStatus === tab.key
                                ? 'bg-blue-600 text-white'
                                : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                        }`}
                    >
                        {tab.label}
                    </button>
                ))}
            </div>
        </div>
    );

    if (loading) {
        return (
            <div className="flex justify-center items-center h-64">
                <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="text-center py-8">
                <div className="text-red-500 mb-2">Ошибка загрузки контрактов</div>
                <button 
                    onClick={() => window.location.reload()}
                    className="text-blue-500 hover:text-blue-700"
                >
                    Попробовать снова
                </button>
            </div>
        );
    }

    if (mode === 'available') {
        return (
            <div className="container mx-auto px-4 py-8">
                {renderContractSection('Контракты на рассмотрении', contracts.pending)}
            </div>
        );
    }

    return (
        <div className="container mx-auto px-4 py-8">
            {renderStatusTabs()}
            {activeStatus === 'pending' && renderContractSection('Контракты на рассмотрении', contracts.pending)}
            {activeStatus === 'active' && renderContractSection('Активные контракты', contracts.active)}
            {activeStatus === 'completed' && renderContractSection('Завершённые контракты', contracts.completed)}
            {activeStatus === 'cancelled' && renderContractSection('Отменённые контракты', contracts.cancelled)}
            {activeStatus === 'banned' && renderContractSection('Заблокированные контракты', contracts.banned)}
        </div>
    );
};

export default ContractFeed; 