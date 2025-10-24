import React, { useEffect, useRef } from 'react';
import styled from 'styled-components';
import { format, formatDistanceToNow } from 'date-fns';
import { ru } from 'date-fns/locale';
import Message from './Message';

const MessagesContainer = styled.div`
  flex: 1;
  overflow-y: auto;
  padding: 16px 24px;
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const MessageGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 4px;
`;

const DateSeparator = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 16px 0;
  
  &::before,
  &::after {
    content: '';
    flex: 1;
    height: 1px;
    background: var(--border-color);
  }
  
  span {
    padding: 8px 16px;
    background: var(--background-color);
    border: 1px solid var(--border-color);
    border-radius: 16px;
    font-size: 12px;
    color: var(--text-secondary);
    font-weight: 500;
  }
`;

const EmptyState = styled.div`
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  color: var(--text-secondary);
  
  h3 {
    font-size: 18px;
    margin-bottom: 8px;
  }
  
  p {
    font-size: 14px;
  }
`;

const LoadingState = styled.div`
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-secondary);
`;

const MessageList = ({ messages, currentUser, loading = false }) => {
  const messagesEndRef = useRef(null);
  const messagesContainerRef = useRef(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const groupMessagesByDate = (messages) => {
    if (!messages || messages.length === 0) return [];

    const groups = [];
    let currentGroup = {
      date: null,
      messages: []
    };

    messages.forEach((message, index) => {
      const messageDate = new Date(message.created_at);
      const messageDateStr = format(messageDate, 'yyyy-MM-dd');

      if (currentGroup.date !== messageDateStr) {
        if (currentGroup.messages.length > 0) {
          groups.push(currentGroup);
        }
        currentGroup = {
          date: messageDateStr,
          messages: [message]
        };
      } else {
        currentGroup.messages.push(message);
      }

      if (index === messages.length - 1 && currentGroup.messages.length > 0) {
        groups.push(currentGroup);
      }
    });

    return groups;
  };

  const formatGroupDate = (dateStr) => {
    const date = new Date(dateStr);
    const today = new Date();
    const yesterday = new Date(today);
    yesterday.setDate(yesterday.getDate() - 1);

    if (format(date, 'yyyy-MM-dd') === format(today, 'yyyy-MM-dd')) {
      return 'Сегодня';
    } else if (format(date, 'yyyy-MM-dd') === format(yesterday, 'yyyy-MM-dd')) {
      return 'Вчера';
    } else {
      return format(date, 'd MMMM yyyy', { locale: ru });
    }
  };

  const messageGroups = groupMessagesByDate(messages);

  if (loading) {
    return (
      <MessagesContainer ref={messagesContainerRef}>
        <LoadingState>
          Загрузка сообщений...
        </LoadingState>
      </MessagesContainer>
    );
  }

  if (!messages || messages.length === 0) {
    return (
      <MessagesContainer ref={messagesContainerRef}>
        <EmptyState>
          <h3>Начните общение!</h3>
          <p>Отправьте первое сообщение в этот чат</p>
        </EmptyState>
        <div ref={messagesEndRef} />
      </MessagesContainer>
    );
  }

  return (
    <MessagesContainer ref={messagesContainerRef}>
      {messageGroups.map((group, groupIndex) => (
        <React.Fragment key={group.date || groupIndex}>
          <DateSeparator>
            <span>{formatGroupDate(group.date)}</span>
          </DateSeparator>
          
          <MessageGroup>
            {group.messages.map((message, messageIndex) => (
              <Message
                key={message.id}
                message={message}
                currentUser={currentUser}
                showAvatar={messageIndex === 0 || 
                  group.messages[messageIndex - 1].user_id !== message.user_id ||
                  (new Date(message.created_at) - new Date(group.messages[messageIndex - 1].created_at)) > 300000 // 5 minutes
                }
              />
            ))}
          </MessageGroup>
        </React.Fragment>
      ))}
      <div ref={messagesEndRef} />
    </MessagesContainer>
  );
};

export default MessageList;
