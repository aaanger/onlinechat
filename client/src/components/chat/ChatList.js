import React from 'react';
import styled from 'styled-components';
import { MessageCircle, Users, Lock } from 'lucide-react';
import { useChat } from '../../contexts/ChatContext';
import { format } from 'date-fns';
import { ru } from 'date-fns/locale';

const ChatListContainer = styled.div`
  flex: 1;
  overflow-y: auto;
  padding: 0 8px;
`;

const ChatItem = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  margin: 2px 0;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  color: white;
  
  &:hover {
    background: rgba(255, 255, 255, 0.1);
  }
  
  ${props => props.active && `
    background: rgba(59, 130, 246, 0.3);
    
    &:hover {
      background: rgba(59, 130, 246, 0.4);
    }
  `}
`;

const ChatAvatar = styled.div`
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: ${props => props.private ? '#f59e0b' : '#3b82f6'};
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
`;

const ChatInfo = styled.div`
  flex: 1;
  min-width: 0;
`;

const ChatName = styled.div`
  font-weight: 600;
  font-size: 14px;
  margin-bottom: 4px;
  display: flex;
  align-items: center;
  gap: 6px;
`;

const ChatPreview = styled.div`
  font-size: 12px;
  color: rgba(255, 255, 255, 0.7);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
`;

const ChatMeta = styled.div`
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 4px;
  flex-shrink: 0;
`;

const LastMessageTime = styled.div`
  font-size: 11px;
  color: rgba(255, 255, 255, 0.6);
`;

const MemberCount = styled.div`
  font-size: 11px;
  color: rgba(255, 255, 255, 0.6);
  display: flex;
  align-items: center;
  gap: 2px;
`;

const EmptyState = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 200px;
  color: rgba(255, 255, 255, 0.6);
  text-align: center;
  padding: 20px;
  
  p {
    margin-top: 8px;
    font-size: 14px;
  }
`;

const LoadingState = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100px;
  color: rgba(255, 255, 255, 0.6);
`;

const ChatList = ({ searchTerm }) => {
  const { chats, currentChat, setCurrentChat, loading } = useChat();

  const filteredChats = chats.filter(chat =>
    chat.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    (chat.description && chat.description.toLowerCase().includes(searchTerm.toLowerCase()))
  );

  const handleChatSelect = (chat) => {
    setCurrentChat(chat);
  };

  const formatLastMessageTime = (timestamp) => {
    if (!timestamp) return '';
    
    try {
      const date = new Date(timestamp);
      const now = new Date();
      const diffInHours = (now - date) / (1000 * 60 * 60);
      
      if (diffInHours < 24) {
        return format(date, 'HH:mm', { locale: ru });
      } else if (diffInHours < 168) { // 7 days
        return format(date, 'EEE', { locale: ru });
      } else {
        return format(date, 'dd.MM', { locale: ru });
      }
    } catch {
      return '';
    }
  };

  if (loading) {
    return (
      <LoadingState>
        Загрузка чатов...
      </LoadingState>
    );
  }

  if (filteredChats.length === 0) {
    return (
      <EmptyState>
        <MessageCircle size={32} />
        <p>
          {searchTerm ? 'Чаты не найдены' : 'У вас пока нет чатов'}
        </p>
      </EmptyState>
    );
  }

  return (
    <ChatListContainer>
      {filteredChats.map((chat) => (
        <ChatItem
          key={chat.id}
          active={currentChat?.id === chat.id}
          onClick={() => handleChatSelect(chat)}
        >
          <ChatAvatar private={chat.is_private}>
            {chat.is_private ? <Lock size={20} /> : <Users size={20} />}
          </ChatAvatar>
          
          <ChatInfo>
            <ChatName>
              {chat.name}
              {chat.is_private && <Lock size={12} />}
            </ChatName>
            <ChatPreview>
              {chat.description || 'Нет описания'}
            </ChatPreview>
          </ChatInfo>
          
          <ChatMeta>
            <LastMessageTime>
              {formatLastMessageTime(chat.created_at)}
            </LastMessageTime>
            <MemberCount>
              <Users size={10} />
              {chat.current_members}/{chat.max_members}
            </MemberCount>
          </ChatMeta>
        </ChatItem>
      ))}
    </ChatListContainer>
  );
};

export default ChatList;
