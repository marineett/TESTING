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

interface Department {
    id: number;
    name: string;
}

const ModeratorList: React.FC = () => {
    const [moderators, setModerators] = useState<ModeratorProfile[]>([]);
    const [departments, setDepartments] = useState<Department[]>([]);
    const [loading, setLoading] = useState(true);
    const [hireError, setHireError] = useState<{ [modId: number]: string | null }>({});
    const [selectedDepartment, setSelectedDepartment] = useState<number | null>(null);
    const [hiringModerator, setHiringModerator] = useState<number | null>(null);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const adminId = localStorage.getItem('user_id');
                if (!adminId) {
                    setLoading(false);
                    return;
                }

                const moderatorsResponse = await fetch(API_ENDPOINTS.ADMIN.GET_MODERATORS(Number(adminId)), {
                    headers: {
                        'Authorization': `Bearer ${localStorage.getItem('token')}`
                    }
                });
                if (!moderatorsResponse.ok) throw new Error('Не удалось загрузить модераторов');
                const moderatorsData = await moderatorsResponse.json();
                setModerators(moderatorsData);

                const departmentsResponse = await fetch(API_ENDPOINTS.ADMIN.GET_DEPARTMENTS(Number(adminId)), {
                    headers: {
                        'Authorization': `Bearer ${localStorage.getItem('token')}`
                    }
                });
                if (!departmentsResponse.ok) throw new Error('Не удалось загрузить отделы');
                const departmentsData = await departmentsResponse.json();
                setDepartments(departmentsData);
            } catch (err) {
            } finally {
                setLoading(false);
            }
        };
        fetchData();
    }, []);

    const handleHire = async (moderatorId: number, departmentId: number) => {
        try {
            const adminId = localStorage.getItem('user_id');
            if (!adminId) return;

            const response = await fetch(API_ENDPOINTS.ADMIN.HIRE_MODERATOR(Number(adminId), departmentId, moderatorId), {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                }
            });

            if (!response.ok) throw new Error('Не удалось нанять модератора');
            
            const updatedModerators = moderators.map(mod => 
                mod.id === moderatorId 
                    ? { ...mod, departments: [...mod.departments, departments.find(d => d.id === departmentId)?.name || ''] }
                    : mod
            );
            setModerators(updatedModerators);
            setHiringModerator(null);
            setSelectedDepartment(null);
            setHireError(prev => ({ ...prev, [moderatorId]: null }));
        } catch (err) {
            setHireError(prev => ({ ...prev, [moderatorId]: err instanceof Error ? err.message : 'Произошла ошибка при найме модератора' }));
        }
    };

    if (loading) {
        return <div>Загрузка модераторов...</div>;
    }

    if (moderators.length === 0) {
        return <div>Нет доступных модераторов</div>;
    }

    return (
      <div style={{ marginTop: 32 }}>
        <h2>Список модераторов</h2>
        <ul>
          {moderators.map(mod => (
            <li key={mod.id} style={{ marginBottom: 16, borderBottom: '1px solid #eee', paddingBottom: 12 }}>
              <div><b>{mod.last_name} {mod.first_name} {mod.middle_name}</b></div>
              <div>Телефон: {mod.telephone_number}</div>
              <div>Email: {mod.email}</div>
              <div>Зарплата: {mod.salary}</div>
              <div>Отделы: {mod.departments.filter(d => d && d.trim()).length > 0 ? (
                <ul style={{ margin: 0, padding: 0, listStyle: 'none' }}>
                  {mod.departments.filter(d => d && d.trim()).map((d, idx) => (
                    <li key={idx}>{d}</li>
                  ))}
                </ul>
              ) : 'Нет отделов'}</div>
              <div style={{ marginTop: 8 }}>
                {hiringModerator === mod.id ? (
                  <>
                    <select
                      value={selectedDepartment || (departments[0]?.id || '')}
                      onChange={e => setSelectedDepartment(Number(e.target.value))}
                      style={{ marginRight: 8 }}
                    >
                      {departments.map(dept => (
                        <option key={dept.id} value={dept.id}>{dept.name}</option>
                      ))}
                    </select>
                    <button
                      onClick={() => selectedDepartment && handleHire(mod.id, selectedDepartment)}
                      disabled={!selectedDepartment}
                      style={{ marginRight: 8 }}
                    >
                      Нанять
                    </button>
                    <button
                      onClick={() => {
                        setHiringModerator(null);
                        setSelectedDepartment(null);
                        setHireError(prev => ({ ...prev, [mod.id]: null }));
                      }}
                    >
                      Отмена
                    </button>
                  </>
                ) : (
                  <button onClick={() => {
                    setHiringModerator(mod.id);
                    if (departments.length > 0) setSelectedDepartment(departments[0].id);
                    setHireError(prev => ({ ...prev, [mod.id]: null }));
                  }}>
                    Нанять в отдел
                  </button>
                )}
                {hireError[mod.id] && (
                  <div style={{ color: 'red', marginTop: 8 }}>{hireError[mod.id]}</div>
                )}
              </div>
            </li>
          ))}
        </ul>
      </div>
    );
};

export default ModeratorList; 