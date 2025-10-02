export interface Contract {
    id: number;
    title: string;
    description: string;
    price: number;
    status: number;
    created_at: string;
    updated_at: string;
    client_id: number;
    repetitor_id: number;
    subject: string;
    payment_status: number;
    commission: number;
    start_date: string;
    end_date: string;
    contract_category: number;
    review_client_id?: number;
    review_repetitor_id?: number;
}

export interface ContractInitInfo {
    title: string;
    description: string;
    price: number;
    client_id: number;
    start_date?: string;
    commission: number;
    duration: number;
    contract_category: number;
    contract_subcategories: number[];
}

export interface ContractListResponse {
    contracts: Contract[];
    total: number;
} 