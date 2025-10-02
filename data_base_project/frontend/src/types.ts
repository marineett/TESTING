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
} 