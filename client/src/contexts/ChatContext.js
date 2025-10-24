import React, { createContext, useContext, useState, useEffect, useCallback } from 'react';
import api from '../services/api';
import { useWebSocket } from '../hooks/useWebSocket';

const ChatContext = createContext();

export const useChat = () => {
  const context = useContext(ChatContext);
  if (!context) {
    throw new Error('useChat must be used within a ChatProvider');
  }
  return context;
};

export const ChatProvider = ({ children }) => {
  const [chats, setChats] = useState([]);
  const [currentChat, setCurrentChat] = useState(null);
  const [messages, setMessages] = useState({});
  const [loading, setLoading] = useState(false);
  const [onlineUsers, setOnlineUsers] = useState(new Set());

  // Handle incoming messages
  const handleWebSocketMessage = useCallback((event) => {
    try {
      const message = JSON.parse(event.data);
      setMessages(prev => ({
        ...prev,
        [message.chat_id]: prev[message.chat_id] 
          ? [...prev[message.chat_id], message]
          : [message]
      }));
      
      // Update last message in chats list if needed
      setChats(prev => prev.map(chat => 
        chat.id === message.chat_id 
          ? { ...chat, last_message: message }
          : chat
      ));
    } catch (error) {
      console.error('Failed to parse message:', error);
    }
  }, []);

  const { 
    socket, 
    isConnected, 
    sendMessage: sendSocketMessage 
  } = useWebSocket(currentChat?.id, handleWebSocketMessage);

  // Load user's chats
  const loadChats = useCallback(async () => {
    try {
      setLoading(true);
      const response = await api.get('/chats/');
      setChats(response.data.chats || []);
    } catch (error) {
      console.error('Failed to load chats:', error);
    } finally {
      setLoading(false);
    }
  }, []);

  // Load messages for a specific chat
  const loadMessages = useCallback(async (chatId, limit = 50, offset = 0) => {
    try {
      const response = await api.get(`/chats/${chatId}/messages`, {
        params: { limit, offset }
      });
      
      setMessages(prev => ({
        ...prev,
        [chatId]: response.data.messages || []
      }));
      
      return response.data;
    } catch (error) {
      console.error('Failed to load messages:', error);
      return { messages: [], total: 0, has_more: false };
    }
  }, []);

  // Create a new chat
  const createChat = async (chatData) => {
    try {
      const response = await api.post('/chats/', chatData);
      const newChat = response.data.chat;
      
      setChats(prev => [newChat, ...prev]);
      setCurrentChat(newChat);
      
      return { success: true, chat: newChat };
    } catch (error) {
      const message = error.response?.data?.error || 'Не удалось создать чат';
      return { success: false, error: message };
    }
  };

  // Search public chats
  const searchPublicChats = useCallback(async (searchTerm = '', limit = 20, offset = 0) => {
    try {
      setLoading(true);
      const response = await api.get('/chats/search', {
        params: { search: searchTerm, limit, offset }
      });
      return { success: true, chats: response.data.chats || [], total: response.data.total || 0 };
    } catch (error) {
      console.error('Failed to search chats:', error);
      return { success: false, error: 'Не удалось найти чаты', chats: [], total: 0 };
    } finally {
      setLoading(false);
    }
  }, []);

  // Join a chat
  const joinChat = async (chatId) => {
    try {
      await api.post(`/chats/${chatId}/join`);
      
      // Reload chats to get updated member count
      await loadChats();
      
      return { success: true };
    } catch (error) {
      const message = error.response?.data?.error || 'Не удалось присоединиться к чату';
      return { success: false, error: message };
    }
  };

  // Leave a chat
  const leaveChat = async (chatId) => {
    try {
      await api.post(`/chats/${chatId}/leave`);
      
      // Remove chat from list and clear messages
      setChats(prev => prev.filter(chat => chat.id !== chatId));
      setMessages(prev => {
        const newMessages = { ...prev };
        delete newMessages[chatId];
        return newMessages;
      });
      
      if (currentChat?.id === chatId) {
        setCurrentChat(null);
      }
      
      return { success: true };
    } catch (error) {
      const message = error.response?.data?.error || 'Не удалось покинуть чат';
      return { success: false, error: message };
    }
  };

  // Send a message
  const sendMessage = async (content, messageType = 'text', replyToId = null) => {
    if (!currentChat?.id || !socket || !isConnected) {
      return { success: false, error: 'Нет соединения с чатом' };
    }

    try {
      const messageData = {
        content,
        message_type: messageType,
        reply_to_id: replyToId
      };

      sendSocketMessage(messageData);
      return { success: true };
    } catch (error) {
      return { success: false, error: 'Не удалось отправить сообщение' };
    }
  };


  // Load chats on mount
  useEffect(() => {
    loadChats();
  }, [loadChats]);

  // Load messages when current chat changes
  useEffect(() => {
    if (currentChat?.id && !messages[currentChat.id]) {
      loadMessages(currentChat.id);
    }
  }, [currentChat, loadMessages]);

  const value = {
    chats,
    currentChat,
    messages: messages[currentChat?.id] || [],
    loading,
    onlineUsers,
    isConnected,
    loadChats,
    loadMessages,
    createChat,
    searchPublicChats,
    joinChat,
    leaveChat,
    sendMessage,
    setCurrentChat,
  };

  return (
    <ChatContext.Provider value={value}>
      {children}
    </ChatContext.Provider>
  );
};
