import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { API_ENDPOINTS } from '../config';

interface Chat {
  id: number;
  created_at: string;
}

interface ChatListWindowProps {
  userId: number;
  chatsOffset?: number;
  chatsLimit?: number;
  role: 'client' | 'repetitor' | 'moderator';
}

const ChatListWindow: React.FC<ChatListWindowProps> = ({ userId, chatsOffset = 0, chatsLimit = 20, role }) => {
  const [chats, setChats] = useState<Chat[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchChats = async () => {
      setLoading(true);
      setError(null);
      try {
        let url = '';
        if (role === 'client') {
          url = API_ENDPOINTS.CHAT.GET_CLIENT_CHATS(userId, chatsOffset, chatsLimit);
        } else if (role === 'repetitor') {
          url = API_ENDPOINTS.CHAT.GET_REPETITOR_CHATS(userId, chatsOffset, chatsLimit);
        } else if (role === 'moderator') {
          url = API_ENDPOINTS.CHAT.GET_MODERATOR_CHATS(userId, chatsOffset, chatsLimit);
        } else {
          url = API_ENDPOINTS.CHAT.GET_CHATS(userId, chatsOffset, chatsLimit);
        }
        const res = await fetch(url, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        });
        if (!res.ok) throw new Error('Ошибка загрузки чатов');
        const data = await res.json();
        console.log('Полученные данные для чатов:', data);
        let chatsArray: Chat[] = [];
        if (Array.isArray(data)) {
          chatsArray = data.filter((chat: any) => chat && typeof chat.id === 'number');
        } else if (Array.isArray(data.chats)) {
          chatsArray = data.chats.filter((chat: any) => chat && typeof chat.id === 'number');
        }
        setChats(chatsArray);
      } catch (e: any) {
        setError(e.message);
      } finally {
        setLoading(false);
      }
    };
    fetchChats();
  }, [userId, chatsOffset, chatsLimit, role]);

  return (
    <div style={{ background: '#f9f9f9', borderRadius: 8, padding: 16, marginBottom: 24 }}>
      <h3 style={{ marginBottom: 12 }}>Ваши чаты</h3>
      {loading && <div>Загрузка...</div>}
      {error && <div style={{ color: 'red' }}>{error}</div>}
      {!loading && !error && chats.length === 0 && <div>Нет чатов</div>}
      <ul style={{ listStyle: 'none', padding: 0 }}>
        {chats.map(chat => (
          <li key={chat.id} style={{ marginBottom: 8, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <span>Чат #{chat.id}</span>
            <button
              onClick={() => navigate(`/chat?id=${chat.id}`)}
              style={{ padding: '4px 12px', borderRadius: 4, background: '#2563eb', color: '#fff', border: 'none', cursor: 'pointer' }}
            >
              Открыть
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default ChatListWindow; 