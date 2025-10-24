import React, { useState } from 'react';
import styled from 'styled-components';
import { Users, Lock, MoreVertical, Settings, LogOut } from 'lucide-react';
import { useChat } from '../../contexts/ChatContext';
import { useAuth } from '../../contexts/AuthContext';

const Header = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 24px;
  background: white;
  border-bottom: 1px solid var(--border-color);
`;

const ChatInfo = styled.div`
  display: flex;
  align-items: center;
  gap: 16px;
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
`;

const ChatDetails = styled.div`
  display: flex;
  flex-direction: column;
`;

const ChatName = styled.div`
  font-weight: 600;
  font-size: 16px;
  color: var(--text-primary);
  display: flex;
  align-items: center;
  gap: 8px;
`;

const ChatMeta = styled.div`
  font-size: 14px;
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  gap: 16px;
  margin-top: 2px;
`;

const MemberCount = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
`;

const HeaderActions = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
`;

const ActionButton = styled.button`
  padding: 8px;
  border-radius: 8px;
  background: none;
  border: none;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
  
  &:hover {
    background: var(--background-color);
    color: var(--text-primary);
  }
`;

const Dropdown = styled.div`
  position: relative;
`;

const DropdownMenu = styled.div`
  position: absolute;
  right: 0;
  top: 100%;
  background: white;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  z-index: 100;
  min-width: 160px;
  overflow: hidden;
`;

const DropdownItem = styled.button`
  width: 100%;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: none;
  border: none;
  color: var(--text-primary);
  text-align: left;
  cursor: pointer;
  transition: all 0.2s ease;
  
  &:hover {
    background: var(--background-color);
  }
  
  ${props => props.danger && `
    color: #ef4444;
    
    &:hover {
      background: #fef2f2;
    }
  `}
`;

const ChatHeader = ({ chat }) => {
  const { leaveChat } = useChat();
  const { user } = useAuth();
  const [showDropdown, setShowDropdown] = useState(false);

  const handleLeaveChat = async () => {
    if (window.confirm(`Вы уверены, что хотите покинуть чат "${chat.name}"?`)) {
      const result = await leaveChat(chat.id);
      if (!result.success) {
        alert(result.error);
      }
    }
    setShowDropdown(false);
  };

  return (
    <Header>
      <ChatInfo>
        <ChatAvatar private={chat.is_private}>
          {chat.is_private ? <Lock size={24} /> : <Users size={24} />}
        </ChatAvatar>
        
        <ChatDetails>
          <ChatName>
            {chat.name}
            {chat.is_private && <Lock size={16} />}
          </ChatName>
          
          <ChatMeta>
            <MemberCount>
              <Users size={14} />
              {chat.current_members}/{chat.max_members} участников
            </MemberCount>
            {chat.description && (
              <span>{chat.description}</span>
            )}
          </ChatMeta>
        </ChatDetails>
      </ChatInfo>

      <HeaderActions>
        <Dropdown>
          <ActionButton onClick={() => setShowDropdown(!showDropdown)}>
            <MoreVertical size={20} />
          </ActionButton>
          
          {showDropdown && (
            <>
              <div 
                style={{ 
                  position: 'fixed', 
                  top: 0, 
                  left: 0, 
                  right: 0, 
                  bottom: 0, 
                  zIndex: 99 
                }} 
                onClick={() => setShowDropdown(false)}
              />
              <DropdownMenu>
                <DropdownItem>
                  <Settings size={16} />
                  Настройки чата
                </DropdownItem>
                <DropdownItem danger onClick={handleLeaveChat}>
                  <LogOut size={16} />
                  Покинуть чат
                </DropdownItem>
              </DropdownMenu>
            </>
          )}
        </Dropdown>
      </HeaderActions>
    </Header>
  );
};

export default ChatHeader;
