import React, { useEffect, useState } from 'react';
import { API_ENDPOINTS } from '../config';

interface Contract {
  id: number;
  title: string;
  status: string;
  // Добавьте другие нужные поля по вашему контракту
}

const Client: React.FC = () => {
  const [contracts, setContracts] = useState<Contract[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const userId = localStorage.getItem('user_id');
    if (!userId) {
      setError('Не найден user_id');
      setLoading(false);
      return;
    }

    fetch(API_ENDPOINTS.CLIENT.CONTRACTS(Number(userId), 0, 10, 0))
      .then(res => {
        if (!res.ok) throw new Error('Ошибка загрузки фида');
        return res.json();
      })
      .then(data => setContracts(data))
      .catch(err => setError(err.message))
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <div>Загрузка...</div>;
  if (error) return <div>Ошибка: {error}</div>;

  return (
    <div>
      <h2>Ваши контракты</h2>
      {contracts.length === 0 ? (
        <div>Нет контрактов</div>
      ) : (
        <ul>
          {contracts.map(contract => (
            <li key={contract.id}>
              <strong>{contract.title}</strong> — статус: {contract.status}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

export default Client; 