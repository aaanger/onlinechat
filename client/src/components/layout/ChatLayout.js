import React from 'react';
import styled from 'styled-components';
import Sidebar from './Sidebar';
import ChatRoom from '../chat/ChatRoom';
import { useChat } from '../../contexts/ChatContext';

const LayoutContainer = styled.div`
  display: flex;
  height: 100vh;
  background: var(--background-color);
`;

const MainContent = styled.main`
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
`;

const EmptyState = styled.div`
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  color: var(--text-secondary);
  
  h2 {
    font-size: 24px;
    margin-bottom: 8px;
  }
  
  p {
    font-size: 16px;
  }
`;

const ChatLayout = () => {
  const { currentChat } = useChat();

  return (
    <LayoutContainer>
      <Sidebar />
      <MainContent>
        {currentChat ? (
          <ChatRoom />
        ) : (
          <EmptyState>
            <h2>Добро пожаловать в Online Chat!</h2>
            <p>Выберите чат из списка или создайте новый</p>
          </EmptyState>
        )}
      </MainContent>
    </LayoutContainer>
  );
};

export default ChatLayout;
