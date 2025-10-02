import React, { useEffect, useState, useRef } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { API_ENDPOINTS } from '../config';
import { ArrowLeftOutlined, SendOutlined } from '@ant-design/icons';
import '../App.css';

interface Chat {
    id: number;
    client_id: number;
    repetitor_id: number;
    created_at: string;
}

interface Message {
    id: number;
    chat_id: number;
    sender_id: number;
    content: string;
    created_at: string;
}

const MESSAGES_PER_PAGE = 20;

const ChatPage: React.FC = () => {
    const navigate = useNavigate();
    const [searchParams] = useSearchParams();
    const [chat, setChat] = useState<Chat | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [message, setMessage] = useState('');
    const [isSending, setIsSending] = useState(false);
    const [messages, setMessages] = useState<Message[]>([]);
    const [messagesOffset, setMessagesOffset] = useState(0);
    const [isLoadingMessages, setIsLoadingMessages] = useState(false);
    const [hasMoreMessages, setHasMoreMessages] = useState(true);
    const messagesEndRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        const chatId = searchParams.get('id');
        if (!chatId) {
            setError('ID чата не указан');
            return;
        }

        const fetchChat = async () => {
            try {
                const response = await fetch(API_ENDPOINTS.CHAT.GET_CHAT(Number(chatId)), {
                    method: 'GET',
                    headers: {
                        'Accept': 'application/json',
                    },
                    credentials: 'include',
                });

                if (!response.ok) {
                    throw new Error('Failed to fetch chat');
                }

                const chatData = await response.json();
                setChat(chatData);
            } catch (error) {
                console.error('Error fetching chat:', error);
                setError('Не удалось загрузить чат');
            }
        };

        fetchChat();
    }, [searchParams]);

    const fetchMessages = async (offset: number) => {
        if (!chat) return;

        setIsLoadingMessages(true);
        try {
            const response = await fetch(API_ENDPOINTS.CHAT.GET_MESSAGES(chat.id, offset, MESSAGES_PER_PAGE), {
                method: 'GET',
                headers: {
                    'Accept': 'application/json',
                },
                credentials: 'include',
            });

            if (!response.ok) {
                throw new Error('Failed to fetch messages');
            }

            const newMessages = await response.json();
            setMessages(prev => {
                const existingIds = new Set(prev.map((m: Message) => m.id));
                const filteredNew = newMessages.filter((m: Message) => !existingIds.has(m.id));
                return [...filteredNew, ...prev];
            });
            setHasMoreMessages(newMessages.length === MESSAGES_PER_PAGE);
        } catch (error) {
            console.error('Error fetching messages:', error);
            setError('Не удалось загрузить сообщения');
        } finally {
            setIsLoadingMessages(false);
        }
    };

    useEffect(() => {
        if (chat) {
            fetchMessages(0);
        }
    }, [chat]);

    const scrollToBottom = () => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    };

    useEffect(() => {
        scrollToBottom();
    }, [messages]);

    const handleLoadMore = () => {
        if (!isLoadingMessages && hasMoreMessages) {
            const newOffset = messagesOffset + MESSAGES_PER_PAGE;
            setMessagesOffset(newOffset);
            fetchMessages(newOffset);
        }
    };

    const handleBack = () => {
        navigate(-1);
    };

    const handleSendMessage = async () => {
        if (!message.trim() || !chat) return;

        setIsSending(true);
        try {
            const userId = localStorage.getItem('user_id');
            if (!userId) {
                throw new Error('User ID not found');
            }

            const response = await fetch(API_ENDPOINTS.CHAT.SEND_MESSAGE(Number(userId), chat.id), {
                method: 'POST',
                headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
                body: JSON.stringify(message),
            });

            if (!response.ok) {
                throw new Error('Failed to send message');
            }

            setMessage('');
            let sentMessage: Message | null = null;
            try {
                const text = await response.text();
                sentMessage = text ? JSON.parse(text) : null;
            } catch {
                sentMessage = null;
            }
            if (sentMessage && sentMessage.id) {
                setMessages(prev => [...prev, sentMessage as Message]);
            } else {
                fetchMessages(messagesOffset);
            }
        } catch (error) {
            console.error('Error sending message:', error);
            alert('Не удалось отправить сообщение');
        } finally {
            setIsSending(false);
        }
    };

    const handleKeyPress = (e: React.KeyboardEvent) => {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
            handleSendMessage();
        }
    };

    const formatMessageTime = (dateString: string) => {
        return new Date(dateString).toLocaleTimeString('ru-RU', {
            hour: '2-digit',
            minute: '2-digit'
        });
    };

    function groupMessagesBySequence(messages: Message[]) {
        if (messages.length === 0) return [];
        const groups: { sender_id: number; messages: Message[] }[] = [];
        let currentGroup = { sender_id: messages[0].sender_id, messages: [messages[0]] };
        for (let i = 1; i < messages.length; i++) {
            const msg = messages[i];
            if (msg.sender_id === currentGroup.sender_id) {
                currentGroup.messages.push(msg);
            } else {
                groups.push(currentGroup);
                currentGroup = { sender_id: msg.sender_id, messages: [msg] };
            }
        }
        groups.push(currentGroup);
        return groups;
    }

    return (
        <div className="page-container">
            <div className="content-container">
                <div className="flex items-center mb-6">
                    <button 
                        onClick={handleBack}
                        className="flex items-center gap-2 px-4 py-2 bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-lg transition-colors mr-4"
                    >
                        <ArrowLeftOutlined />
                        Назад
                    </button>
                    <h1 className="text-2xl font-semibold">Чат #{chat?.id}</h1>
                </div>

                {error ? (
                    <div className="text-red-500">{error}</div>
                ) : !chat ? (
                    <div>Загрузка...</div>
                ) : (
                    <div className="flex flex-col h-[calc(100vh-200px)]">
                        <div
                            className="flex-1 overflow-y-auto mb-4 p-4 bg-gray-50 rounded-lg"
                            style={{ width: '100%' }}
                        >
                            <div className="space-y-6">
                                {hasMoreMessages && !isLoadingMessages && (
                                    <button
                                        onClick={handleLoadMore}
                                        className="w-full py-2 text-blue-500 hover:text-blue-600"
                                    >
                                        Загрузить еще
                                    </button>
                                )}
                                {isLoadingMessages && (
                                    <div className="text-center text-gray-500">Загрузка сообщений...</div>
                                )}
                                {groupMessagesBySequence(
                                    messages.sort((a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime())
                                ).map((group, idx) => (
                                    <div key={idx} className="mb-4 w-full">
                                        <div className="space-y-1 flex flex-col w-full">
                                            {group.messages.map((msg) => {
                                                const isMe = msg.sender_id === Number(localStorage.getItem('user_id'));
                                                return (
                                                    <div
                                                        key={msg.id}
                                                        className={`p-3 rounded-lg ${isMe ? 'bg-blue-500 text-white' : 'bg-gray-200 text-gray-800'}`}
                                                        style={{
                                                            maxWidth: 480,
                                                            wordBreak: 'break-word',
                                                            overflowWrap: 'break-word',
                                                            whiteSpace: 'pre-line',
                                                            marginRight: isMe ? 'auto' : undefined,
                                                            marginLeft: !isMe ? 'auto' : undefined,
                                                        }}
                                                    >
                                                        <div className="text-sm">{msg.content}</div>
                                                        <div className={`text-xs mt-1 ${isMe ? 'text-blue-100' : 'text-gray-500'}`}>{formatMessageTime(msg.created_at)}</div>
                                                    </div>
                                                );
                                            })}
                                        </div>
                                    </div>
                                ))}
                                <div ref={messagesEndRef} />
                            </div>
                        </div>
                        
                        <div className="flex gap-2">
                            <textarea
                                value={message}
                                onChange={(e) => setMessage(e.target.value)}
                                onKeyPress={handleKeyPress}
                                placeholder="Введите сообщение..."
                                className="flex-1 p-3 border border-gray-300 rounded-lg focus:outline-none focus:border-blue-500 resize-none"
                                rows={3}
                            />
                            <button
                                onClick={handleSendMessage}
                                disabled={isSending || !message.trim()}
                                className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
                            >
                                <SendOutlined />
                                Отправить
                            </button>
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
};

export default ChatPage; 