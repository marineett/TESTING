import React, { useEffect, useState } from 'react';
import { API_ENDPOINTS } from '../config';

interface ModeratorProfile {
    id: number;
    first_name: string;
    last_name: string;
    middle_name: string;
    telephone_number: string;
    email: string;
    salary: number;
    departments: string[];
}

interface CompleteDepartmentInfo {
    id: number;
    name: string;
    head_id: number;
    Moderators: ModeratorProfile[];
}

const DepartmentList: React.FC = () => {
    const [departments, setDepartments] = useState<CompleteDepartmentInfo[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [editingModeratorId, setEditingModeratorId] = useState<number | null>(null);
    const [newSalary, setNewSalary] = useState<string>('');

    useEffect(() => {
        const fetchDepartments = async () => {
            const adminId = localStorage.getItem('user_id');
            if (!adminId) {
                setError('Не найден id администратора');
                setLoading(false);
                return;
            }
            try {
                const response = await fetch(API_ENDPOINTS.ADMIN.GET_DEPARTMENTS(Number(adminId)), {
                    headers: {
                        'Authorization': `Bearer ${localStorage.getItem('token')}`
                    }
                });
                if (!response.ok) throw new Error('Не удалось загрузить отделы');
                const data = await response.json();
                setDepartments(data);
            } catch (err) {
                setError(err instanceof Error ? err.message : 'Произошла ошибка');
            } finally {
                setLoading(false);
            }
        };
        fetchDepartments();
    }, []);

    const handleFire = async (departmentId: number, moderatorId: number) => {
        try {
            const adminId = localStorage.getItem('user_id');
            if (!adminId || !moderatorId) return;
            const response = await fetch(API_ENDPOINTS.ADMIN.FIRE_MODERATOR(Number(adminId), departmentId, moderatorId), {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                }
            });
            if (!response.ok) throw new Error('Не удалось уволить модератора');
            setDepartments(prev => prev.map(dep =>
                dep.id === departmentId
                    ? { ...dep, Moderators: dep.Moderators.filter(mod => mod.id !== moderatorId) }
                    : dep
            ));
        } catch (err) {
            alert('Ошибка при увольнении модератора');
        }
    };

    const handleChangeSalary = async (departmentId: number, moderatorId: number) => {
        try {
            const adminId = localStorage.getItem('user_id');
            if (!adminId || !moderatorId) return;
            const response = await fetch(API_ENDPOINTS.ADMIN.CHANGE_MODERATOR_SALARY(Number(adminId), departmentId, moderatorId, newSalary), {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                }
            });
            if (!response.ok) throw new Error('Не удалось изменить зарплату');
            setDepartments(prev => prev.map(dep =>
                dep.id === departmentId
                    ? {
                        ...dep,
                        Moderators: dep.Moderators.map(mod =>
                            mod.id === moderatorId ? { ...mod, salary: Number(newSalary) } : mod
                        )
                    }
                    : dep
            ));
            setEditingModeratorId(null);
            setNewSalary('');
        } catch (err) {
            alert('Ошибка при изменении зарплаты');
        }
    };

    if (loading) {
        return <div className="py-8 text-center">Загрузка отделов...</div>;
    }

    if (error) {
        return <div className="text-red-600 py-4">{error}</div>;
    }

    if (departments.length === 0) {
        return <div className="py-4 text-gray-500">Нет отделов</div>;
    }

    return (
        <div className="departments-root">
            <h2 className="departments-title">Список отделов</h2>
            <div className="departments-row">
                {departments.map(dep => (
                    <div key={dep.id} className="department-card">
                        <div className="department-name">{dep.name}</div>
                        <div className="department-id">ID отдела: {dep.id}</div>
                        <div>
                            <div className="department-mods-label">Модераторы:</div>
                            {(dep.Moderators.length === 0) ? (
                                <div className="department-no-mods">Нет модераторов</div>
                            ) : (
                                <ul className="department-mod-list">
                                    {dep.Moderators.map((mod, idx) => (
                                        <li key={mod.email + idx} style={{
                                            display: 'block',
                                            marginBottom: 16,
                                            padding: 8,
                                            border: '1px solid #f0f0f0',
                                            borderRadius: 6,
                                            background: '#fafbfc'
                                        }}>
                                            <div><b>{mod.last_name} {mod.first_name} {mod.middle_name}</b></div>
                                            <div>Телефон: {mod.telephone_number}</div>
                                            <div>Email: {mod.email}</div>
                                            <div>Зарплата: {mod.salary}</div>
                                            <div style={{ display: 'flex', gap: 8, marginTop: 8 }}>
                                                {editingModeratorId === mod.id ? (
                                                    <div style={{ display: 'flex', gap: 8 }}>
                                                        <input
                                                            type="number"
                                                            value={newSalary}
                                                            onChange={(e) => setNewSalary(e.target.value)}
                                                            placeholder="Новая зарплата"
                                                            style={{ padding: '2px 6px', borderRadius: 4, border: '1px solid #ccc', fontSize: 13, height: 28 }}
                                                        />
                                                        <button
                                                            style={{ padding: '2px 10px', borderRadius: 4, border: '1px solid #4caf50', background: '#e8f5e9', color: '#2e7d32', cursor: 'pointer', fontSize: 13, height: 28 }}
                                                            onClick={() => handleChangeSalary(dep.id, mod.id)}
                                                        >
                                                            Сохранить
                                                        </button>
                                                        <button
                                                            style={{ padding: '2px 10px', borderRadius: 4, border: '1px solid #bbb', background: '#f5f5f5', cursor: 'pointer', fontSize: 13, height: 28 }}
                                                            onClick={() => {
                                                                setEditingModeratorId(null);
                                                                setNewSalary('');
                                                            }}
                                                        >
                                                            Отмена
                                                        </button>
                                                    </div>
                                                ) : (
                                                    <button
                                                        style={{ padding: '2px 10px', borderRadius: 4, border: '1px solid #bbb', background: '#f5f5f5', cursor: 'pointer', fontSize: 13, height: 28 }}
                                                        onClick={() => setEditingModeratorId(mod.id)}
                                                    >
                                                        Поменять зарплату
                                                    </button>
                                                )}
                                                <button style={{ padding: '2px 10px', borderRadius: 4, border: '1px solid #e57373', background: '#ffeaea', color: '#c62828', cursor: 'pointer', fontSize: 13, height: 28 }}
                                                    onClick={() => handleFire(dep.id, mod.id)}
                                                >
                                                    Уволить
                                                </button>
                                            </div>
                                        </li>
                                    ))}
                                </ul>
                            )}
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default DepartmentList; 