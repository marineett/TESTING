export interface Review {
    id: number;
    client_id: number;
    repetitor_id: number;
    rating: number;
    comment: string;
    created_at: string;
}

export interface CreateReviewData {
    client_id: number;
    repetitor_id: number;
    rating: number;
    comment: string;
} 