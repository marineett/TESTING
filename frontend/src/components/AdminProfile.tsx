import React, { useEffect, useState } from 'react';
import { API_ENDPOINTS } from '../config';

interface AdminProfile {
    first_name: string;
    last_name: string;
    middle_name: string;
    telephone_number: string;
    email: string;
    salary: number;
}

interface AdminProfileProps {
    id: string;
}

const AdminProfile: React.FC<AdminProfileProps> = ({ id }) => {
    const [profile, setProfile] = useState<AdminProfile | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchProfile = async () => {
            try {
                const response = await fetch(API_ENDPOINTS.ADMIN.GET_PROFILE(Number(id)), {
                    headers: {
                        'Authorization': `Bearer ${localStorage.getItem('token')}`
                    }
                });
                if (!response.ok) throw new Error('Не удалось загрузить профиль администратора');
                const data = await response.json();
                setProfile(data);
            } catch (err) {
                setError(err instanceof Error ? err.message : 'Произошла ошибка');
            } finally {
                setLoading(false);
            }
        };
        fetchProfile();
    }, [id]);

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

    if (!profile) return null;

    return (
        <div style={{ marginBottom: 24, padding: 16, background: '#f5f5f5', borderRadius: 8 }}>
            <div style={{ fontSize: 20, fontWeight: 600, marginBottom: 8 }}>
                    Профиль администратора
            </div>
            <div style={{ color: '#222', marginBottom: 4 }}><b>ФИО</b></div>
            <div style={{ marginBottom: 4 }}>{profile.last_name} {profile.first_name} {profile.middle_name}</div>
            <div style={{ color: '#222', marginBottom: 4 }}><b>Email</b></div>
            <div style={{ marginBottom: 4 }}>{profile.email}</div>
            <div style={{ color: '#222', marginBottom: 4 }}><b>Телефон</b></div>
            <div style={{ marginBottom: 4 }}>{profile.telephone_number}</div>
            <div style={{ color: '#222', marginBottom: 4 }}><b>Зарплата</b></div>
            <div>{profile.salary.toLocaleString()} ₽</div>
        </div>
    );
};

export default AdminProfile; 