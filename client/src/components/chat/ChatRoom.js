import React, { useEffect, useRef } from 'react';
import styled from 'styled-components';
import { Users, Wifi, WifiOff, Settings } from 'lucide-react';
import { useChat } from '../../contexts/ChatContext';
import { useAuth } from '../../contexts/AuthContext';
import MessageList from './MessageList';
import MessageInput from './MessageInput';
import ChatHeader from './ChatHeader';

const ChatRoomContainer = styled.div`
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: var(--chat-bg);
`;

const MessagesContainer = styled.div`
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
`;

const ConnectionStatus = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 8px 16px;
  background: ${props => props.connected ? '#10b981' : '#ef4444'};
  color: white;
  font-size: 12px;
  font-weight: 500;
`;

const ChatRoom = () => {
  const { currentChat, messages, isConnected } = useChat();
  const { user } = useAuth();

  if (!currentChat) {
    return null;
  }

  return (
    <ChatRoomContainer>
      <ChatHeader chat={currentChat} />
      
      {!isConnected && (
        <ConnectionStatus connected={false}>
          <WifiOff size={14} />
          Нет соединения
        </ConnectionStatus>
      )}
      
      {isConnected && (
        <ConnectionStatus connected={true}>
          <Wifi size={14} />
          Подключено
        </ConnectionStatus>
      )}

      <MessagesContainer>
        <MessageList messages={messages} currentUser={user} />
      </MessagesContainer>

      <MessageInput 
        disabled={!isConnected}
        placeholder={
          !isConnected 
            ? 'Нет соединения...' 
            : currentChat.is_private 
              ? 'Напишите сообщение...' 
              : 'Отправить сообщение в чат...'
        }
      />
    </ChatRoomContainer>
  );
};

export default ChatRoom;
