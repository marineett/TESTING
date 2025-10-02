import React, { useEffect, useState } from 'react';
import { API_ENDPOINTS } from '../config';

interface Resume {
    id: number;
    repetitor_id: number;
    title: string;
    description: string;
    price: Record<string, number>;
    created_at: string;
    updated_at: string;
}

interface RepetitorProfile {
    id: number;
    first_name: string;
    last_name: string;
    middle_name: string;
    telephone_number: string;
    email: string;
    resume: Resume;
}

const RepetitorProfile: React.FC = () => {
    const [profile, setProfile] = useState<RepetitorProfile | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [isEditing, setIsEditing] = useState(false);
    const [editedResume, setEditedResume] = useState<Resume | null>(null);

    useEffect(() => {
        const fetchProfile = async () => {
            const repetitorId = localStorage.getItem('user_id');
            if (!repetitorId) {
                setError('Не найден id репетитора');
                setLoading(false);
                return;
            }
            try {
                const response = await fetch(API_ENDPOINTS.REPETITOR.GET_PROFILE(Number(repetitorId), 0, 5), {
                    headers: {
                        'Authorization': `Bearer ${localStorage.getItem('token')}`
                    }
                });
                if (!response.ok) throw new Error('Не удалось загрузить профиль');
                const data = await response.json();
                setProfile(data);
                setEditedResume(data.resume);
            } catch (err) {
                setError(err instanceof Error ? err.message : 'Произошла ошибка');
            } finally {
                setLoading(false);
            }
        };
        fetchProfile();
    }, []);

    const handleSaveResume = async () => {
        if (!profile || !editedResume) return;
        try {
            const response = await fetch(API_ENDPOINTS.REPETITOR.CHANGE_RESUME(profile.id), {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                },
                body: JSON.stringify(editedResume)
            });
            if (!response.ok) throw new Error('Не удалось обновить резюме');
            setProfile(prev => prev ? { ...prev, resume: editedResume } : null);
            setIsEditing(false);
        } catch (err) {
            alert('Ошибка при обновлении резюме');
        }
    };

    if (loading) {
        return <div className="py-8 text-center">Загрузка профиля...</div>;
    }

    if (error) {
        return <div className="text-red-600 py-4">{error}</div>;
    }

    if (!profile) {
        return <div className="py-4 text-gray-500">Профиль не найден</div>;
    }

    return (
        <div className="profile-root">
            <h2 className="profile-title">Профиль репетитора</h2>
            <div className="profile-info">
                <div><b>ФИО:</b> {profile.last_name} {profile.first_name} {profile.middle_name}</div>
                <div><b>Телефон:</b> {profile.telephone_number}</div>
                <div><b>Email:</b> {profile.email}</div>
                <div>
                    <b>Резюме:</b>
                    {isEditing ? (
                        <div className="resume-edit">
                            <div>
                                <label>Название:</label>
                                <input
                                    type="text"
                                    value={editedResume?.title || ''}
                                    onChange={(e) => editedResume && setEditedResume({ ...editedResume, title: e.target.value })}
                                    style={{ marginLeft: 8, padding: 4, borderRadius: 4, border: '1px solid #ccc' }}
                                />
                            </div>
                            <div style={{ marginTop: 8 }}>
                                <label>Описание:</label>
                                <textarea
                                    value={editedResume?.description || ''}
                                    onChange={(e) => editedResume && setEditedResume({ ...editedResume, description: e.target.value })}
                                    style={{ marginLeft: 8, padding: 4, borderRadius: 4, border: '1px solid #ccc', width: '100%', minHeight: 100 }}
                                />
                            </div>
                            <div style={{ marginTop: 8 }}>
                                <label>Цены:</label>
                                <textarea
                                    value={JSON.stringify(editedResume?.price || {}, null, 2)}
                                    onChange={(e) => {
                                        if (editedResume) {
                                            try {
                                                const parsed = JSON.parse(e.target.value);
                                                setEditedResume({ ...editedResume, price: parsed });
                                            } catch (err) {
                                                // Игнорировать ошибки парсинга
                                            }
                                        }
                                    }}
                                    style={{ marginLeft: 8, padding: 4, borderRadius: 4, border: '1px solid #ccc', width: '100%', minHeight: 100 }}
                                />
                            </div>
                            <div style={{ marginTop: 8, display: 'flex', gap: 8 }}>
                                <button
                                    onClick={handleSaveResume}
                                    style={{ padding: '4px 12px', borderRadius: 4, border: '1px solid #4caf50', background: '#e8f5e9', color: '#2e7d32', cursor: 'pointer' }}
                                >
                                    Сохранить
                                </button>
                                <button
                                    onClick={() => {
                                        setIsEditing(false);
                                        setEditedResume(profile.resume);
                                    }}
                                    style={{ padding: '4px 12px', borderRadius: 4, border: '1px solid #bbb', background: '#f5f5f5', cursor: 'pointer' }}
                                >
                                    Отмена
                                </button>
                            </div>
                        </div>
                    ) : (
                        <div className="resume-info">
                            <div><b>Название:</b> {profile.resume.title}</div>
                            <div><b>Описание:</b> {profile.resume.description}</div>
                            <div><b>Цены:</b>
                                <ul style={{ margin: 0, paddingLeft: 20 }}>
                                    {Object.entries(profile.resume.price).map(([key, value]) => (
                                        <li key={key}>{key}: {value} ₽</li>
                                    ))}
                                </ul>
                            </div>
                            <button
                                onClick={() => setIsEditing(true)}
                                style={{ marginTop: 8, padding: '4px 12px', borderRadius: 4, border: '1px solid #bbb', background: '#f5f5f5', cursor: 'pointer' }}
                            >
                                Редактировать резюме
                            </button>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
};

export default RepetitorProfile; 