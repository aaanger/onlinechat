import React, { useState } from 'react';
import styled from 'styled-components';
import { 
  Plus, 
  Search, 
  MessageCircle, 
  Users, 
  Settings, 
  LogOut,
  User 
} from 'lucide-react';
import { useAuth } from '../../contexts/AuthContext';
import { useChat } from '../../contexts/ChatContext';
import ChatList from '../chat/ChatList';
import CreateChatModal from '../chat/CreateChatModal';
import SearchChatModal from '../chat/SearchChatModal';

const SidebarContainer = styled.div`
  width: 320px;
  background: var(--sidebar-bg);
  display: flex;
  flex-direction: column;
  border-right: 1px solid var(--border-color);
`;

const Header = styled.div`
  padding: 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
`;

const UserInfo = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
  color: white;
`;

const Avatar = styled.div`
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: var(--primary-color);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  color: white;
`;

const UserDetails = styled.div`
  flex: 1;
`;

const Username = styled.div`
  font-weight: 600;
  font-size: 14px;
`;

const UserStatus = styled.div`
  font-size: 12px;
  color: rgba(255, 255, 255, 0.7);
`;

const HeaderActions = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 16px;
`;

const ActionButton = styled.button`
  padding: 8px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.1);
  color: white;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
  
  &:hover {
    background: rgba(255, 255, 255, 0.2);
  }
`;

const SearchBar = styled.div`
  position: relative;
  margin-top: 16px;
`;

const SearchInput = styled.input`
  width: 100%;
  padding: 12px 12px 12px 40px;
  background: rgba(255, 255, 255, 0.1);
  border: none;
  border-radius: 8px;
  color: white;
  font-size: 14px;
  
  &::placeholder {
    color: rgba(255, 255, 255, 0.6);
  }
  
  &:focus {
    outline: none;
    background: rgba(255, 255, 255, 0.15);
  }
`;

const SearchIcon = styled.div`
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: rgba(255, 255, 255, 0.6);
`;

const ChatSection = styled.div`
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
`;

const SectionHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px 8px;
  color: rgba(255, 255, 255, 0.8);
  font-size: 14px;
  font-weight: 500;
`;

const AddButton = styled.button`
  background: none;
  border: none;
  color: rgba(255, 255, 255, 0.6);
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  
  &:hover {
    color: white;
    background: rgba(255, 255, 255, 0.1);
  }
`;

const Footer = styled.div`
  padding: 16px 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
`;

const FooterButton = styled.button`
  width: 100%;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: none;
  border: none;
  color: rgba(255, 255, 255, 0.8);
  cursor: pointer;
  border-radius: 8px;
  font-size: 14px;
  transition: all 0.2s ease;
  
  &:hover {
    background: rgba(255, 255, 255, 0.1);
    color: white;
  }
`;

const LogoutButton = styled(FooterButton)`
  color: #ef4444;
  
  &:hover {
    background: rgba(239, 68, 68, 0.1);
    color: #ef4444;
  }
`;

const Sidebar = () => {
  const { user, logout } = useAuth();
  const { currentChat } = useChat();
  const [searchTerm, setSearchTerm] = useState('');
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showSearchModal, setShowSearchModal] = useState(false);

  const handleLogout = () => {
    if (window.confirm('Вы уверены, что хотите выйти?')) {
      logout();
    }
  };

  return (
    <>
      <SidebarContainer>
        <Header>
          <UserInfo>
            <Avatar>
              <User size={20} />
            </Avatar>
            <UserDetails>
              <Username>{user?.username}</Username>
              <UserStatus>В сети</UserStatus>
            </UserDetails>
          </UserInfo>
          
          <HeaderActions>
            <ActionButton title="Настройки">
              <Settings size={16} />
            </ActionButton>
          </HeaderActions>
          
          <SearchBar>
            <SearchIcon>
              <Search size={16} />
            </SearchIcon>
            <SearchInput
              type="text"
              placeholder="Поиск чатов..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
            />
          </SearchBar>
        </Header>

        <ChatSection>
          <SectionHeader>
            <span>Чаты</span>
            <div style={{ display: 'flex', gap: '8px' }}>
              <AddButton onClick={() => setShowSearchModal(true)} title="Найти чаты">
                <Search size={16} />
              </AddButton>
              <AddButton onClick={() => setShowCreateModal(true)} title="Создать чат">
                <Plus size={16} />
              </AddButton>
            </div>
          </SectionHeader>
          
          <ChatList searchTerm={searchTerm} />
        </ChatSection>

        <Footer>
          <FooterButton onClick={() => {/* Settings */}}>
            <Settings size={16} />
            Настройки
          </FooterButton>
          <LogoutButton onClick={handleLogout}>
            <LogOut size={16} />
            Выйти
          </LogoutButton>
        </Footer>
      </SidebarContainer>

      {showCreateModal && (
        <CreateChatModal onClose={() => setShowCreateModal(false)} />
      )}
      
      {showSearchModal && (
        <SearchChatModal onClose={() => setShowSearchModal(false)} />
      )}
    </>
  );
};

export default Sidebar;
