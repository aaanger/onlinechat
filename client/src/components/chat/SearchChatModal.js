import React, { useState, useEffect } from 'react';
import styled from 'styled-components';
import { X, Search, Users, Lock, Plus } from 'lucide-react';
import { useChat } from '../../contexts/ChatContext';
import toast from 'react-hot-toast';

const ModalOverlay = styled.div`
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
`;

const Modal = styled.div`
  background: white;
  border-radius: 16px;
  padding: 24px;
  width: 100%;
  max-width: 600px;
  max-height: 80vh;
  overflow-y: auto;
`;

const ModalHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
`;

const Title = styled.h2`
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
  display: flex;
  align-items: center;
  gap: 8px;
`;

const CloseButton = styled.button`
  background: none;
  border: none;
  color: var(--text-muted);
  cursor: pointer;
  padding: 8px;
  border-radius: 8px;
  
  &:hover {
    background: var(--background-color);
    color: var(--text-secondary);
  }
`;

const SearchContainer = styled.div`
  position: relative;
  margin-bottom: 20px;
`;

const SearchInput = styled.input`
  width: 100%;
  padding: 12px 16px 12px 48px;
  border: 2px solid var(--border-color);
  border-radius: 8px;
  font-size: 16px;
  transition: all 0.2s ease;
  
  &:focus {
    outline: none;
    border-color: var(--primary-color);
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
  }
`;

const SearchIcon = styled.div`
  position: absolute;
  left: 16px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--text-muted);
`;

const ChatList = styled.div`
  max-height: 400px;
  overflow-y: auto;
`;

const ChatItem = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  margin: 8px 0;
  border: 2px solid var(--border-color);
  border-radius: 12px;
  transition: all 0.2s ease;
  cursor: pointer;
  
  &:hover {
    border-color: var(--primary-color);
    background: rgba(59, 130, 246, 0.05);
  }
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
  font-size: 16px;
  margin-bottom: 4px;
  display: flex;
  align-items: center;
  gap: 6px;
  color: var(--text-primary);
`;

const ChatDescription = styled.div`
  font-size: 14px;
  color: var(--text-secondary);
  margin-bottom: 8px;
  line-height: 1.4;
`;

const ChatMeta = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 12px;
  color: var(--text-muted);
`;

const JoinButton = styled.button`
  background: var(--primary-color);
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 6px;
  
  &:hover {
    background: var(--primary-hover);
  }
  
  &:disabled {
    background: var(--text-muted);
    cursor: not-allowed;
  }
`;

const EmptyState = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 200px;
  color: var(--text-muted);
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
  color: var(--text-muted);
`;

const SearchChatModal = ({ onClose }) => {
  const [searchTerm, setSearchTerm] = useState('');
  const [searchResults, setSearchResults] = useState([]);
  const [loading, setLoading] = useState(false);
  const [joiningChats, setJoiningChats] = useState(new Set());
  const { searchPublicChats, joinChat } = useChat();

  const handleSearch = async (term) => {
    if (term.trim() === '') {
      setSearchResults([]);
      return;
    }

    setLoading(true);
    const result = await searchPublicChats(term, 20, 0);
    
    if (result.success) {
      setSearchResults(result.chats);
    } else {
      toast.error(result.error);
      setSearchResults([]);
    }
    setLoading(false);
  };

  const handleJoinChat = async (chatId) => {
    setJoiningChats(prev => new Set([...prev, chatId]));
    
    const result = await joinChat(chatId);
    
    if (result.success) {
      toast.success('Успешно присоединились к чату!');
      onClose();
    } else {
      toast.error(result.error);
    }
    
    setJoiningChats(prev => {
      const newSet = new Set(prev);
      newSet.delete(chatId);
      return newSet;
    });
  };

  useEffect(() => {
    const timeoutId = setTimeout(() => {
      handleSearch(searchTerm);
    }, 300);

    return () => clearTimeout(timeoutId);
  }, [searchTerm]);

  return (
    <ModalOverlay onClick={onClose}>
      <Modal onClick={(e) => e.stopPropagation()}>
        <ModalHeader>
          <Title>
            <Search size={20} />
            Поиск чатов
          </Title>
          <CloseButton onClick={onClose}>
            <X size={20} />
          </CloseButton>
        </ModalHeader>

        <SearchContainer>
          <SearchIcon>
            <Search size={20} />
          </SearchIcon>
          <SearchInput
            type="text"
            placeholder="Поиск по названию или описанию чата..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </SearchContainer>

        {loading ? (
          <LoadingState>
            Поиск чатов...
          </LoadingState>
        ) : searchResults.length === 0 && searchTerm ? (
          <EmptyState>
            <Search size={32} />
            <p>Чаты не найдены</p>
          </EmptyState>
        ) : searchResults.length === 0 ? (
          <EmptyState>
            <Users size={32} />
            <p>Введите поисковый запрос для поиска чатов</p>
          </EmptyState>
        ) : (
          <ChatList>
            {searchResults.map((chat) => (
              <ChatItem key={chat.id}>
                <ChatAvatar private={chat.is_private}>
                  {chat.is_private ? <Lock size={20} /> : <Users size={20} />}
                </ChatAvatar>
                
                <ChatInfo>
                  <ChatName>
                    {chat.name}
                    {chat.is_private && <Lock size={14} />}
                  </ChatName>
                  <ChatDescription>
                    {chat.description || 'Нет описания'}
                  </ChatDescription>
                  <ChatMeta>
                    <Users size={12} />
                    {chat.current_members}/{chat.max_members} участников
                  </ChatMeta>
                </ChatInfo>
                
                <JoinButton
                  onClick={() => handleJoinChat(chat.id)}
                  disabled={joiningChats.has(chat.id)}
                >
                  <Plus size={14} />
                  {joiningChats.has(chat.id) ? 'Присоединение...' : 'Присоединиться'}
                </JoinButton>
              </ChatItem>
            ))}
          </ChatList>
        )}
      </Modal>
    </ModalOverlay>
  );
};

export default SearchChatModal;
