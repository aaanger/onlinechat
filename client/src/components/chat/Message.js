import React, { useState } from 'react';
import styled from 'styled-components';
import { format } from 'date-fns';
import { ru } from 'date-fns/locale';
import { MoreVertical, Reply, Edit, Trash2 } from 'lucide-react';

const MessageContainer = styled.div`
  display: flex;
  gap: 12px;
  margin-bottom: 4px;
  position: relative;
  
  &:hover {
    .message-actions {
      opacity: 1;
    }
  }
`;

const MessageBubble = styled.div`
  display: flex;
  flex-direction: column;
  max-width: 70%;
  position: relative;
  
  ${props => props.isOwn ? `
    align-self: flex-end;
    align-items: flex-end;
  ` : `
    align-self: flex-start;
    align-items: flex-start;
  `}
`;

const MessageContent = styled.div`
  background: ${props => props.isOwn ? 'var(--message-self)' : 'var(--message-other)'};
  color: ${props => props.isOwn ? 'white' : 'var(--text-primary)'};
  padding: 12px 16px;
  border-radius: ${props => props.isOwn ? '18px 18px 4px 18px' : '18px 18px 18px 4px'};
  word-wrap: break-word;
  line-height: 1.4;
`;

const MessageHeader = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: ${props => props.showAvatar ? '4px' : '0'};
  
  ${props => props.isOwn ? 'flex-direction: row-reverse;' : ''}
`;

const Avatar = styled.div`
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: ${props => props.isOwn ? 'var(--primary-color)' : '#64748b'};
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 600;
  font-size: 12px;
  flex-shrink: 0;
`;

const Username = styled.span`
  font-weight: 600;
  font-size: 14px;
  color: var(--text-secondary);
`;

const MessageTime = styled.span`
  font-size: 11px;
  color: ${props => props.isOwn ? 'rgba(255, 255, 255, 0.7)' : 'var(--text-muted)'};
  margin-top: 4px;
  align-self: ${props => props.isOwn ? 'flex-end' : 'flex-start'};
`;

const MessageActions = styled.div`
  position: absolute;
  top: -8px;
  right: ${props => props.isOwn ? 'auto' : '8px'};
  left: ${props => props.isOwn ? '8px' : 'auto'};
  background: white;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.2s ease;
  z-index: 10;
`;

const ActionButton = styled.button`
  padding: 8px;
  background: none;
  border: none;
  color: var(--text-secondary);
  cursor: pointer;
  border-radius: 4px;
  
  &:hover {
    background: var(--background-color);
    color: var(--text-primary);
  }
`;

const ReplyPreview = styled.div`
  background: rgba(0, 0, 0, 0.05);
  border-left: 3px solid var(--primary-color);
  padding: 8px 12px;
  margin-bottom: 8px;
  border-radius: 4px;
  font-size: 12px;
  color: var(--text-secondary);
`;

const Message = ({ message, currentUser, showAvatar = true }) => {
  const [showActions, setShowActions] = useState(false);
  const isOwn = message.user_id === currentUser?.id;

  const formatTime = (timestamp) => {
    try {
      return format(new Date(timestamp), 'HH:mm', { locale: ru });
    } catch {
      return '';
    }
  };

  const handleReply = () => {
    // TODO: Implement reply functionality
    console.log('Reply to message:', message.id);
  };

  const handleEdit = () => {
    // TODO: Implement edit functionality
    console.log('Edit message:', message.id);
  };

  const handleDelete = () => {
    // TODO: Implement delete functionality
    console.log('Delete message:', message.id);
  };

  return (
    <MessageContainer>
      {!isOwn && showAvatar && (
        <Avatar>
          {message.username?.charAt(0)?.toUpperCase() || 'U'}
        </Avatar>
      )}
      
      {isOwn && showAvatar && (
        <div style={{ width: '32px' }} />
      )}

      <MessageBubble isOwn={isOwn}>
        {(showAvatar || isOwn) && (
          <MessageHeader isOwn={isOwn} showAvatar={showAvatar}>
            {!isOwn && <Username>{message.username}</Username>}
          </MessageHeader>
        )}

        <MessageContent isOwn={isOwn}>
          {message.reply_to_id && (
            <ReplyPreview>
              Ответ на сообщение #{message.reply_to_id}
            </ReplyPreview>
          )}
          {message.content}
        </MessageContent>

        <MessageTime isOwn={isOwn}>
          {formatTime(message.created_at)}
        </MessageTime>

        <MessageActions 
          className="message-actions" 
          isOwn={isOwn}
          onMouseEnter={() => setShowActions(true)}
          onMouseLeave={() => setShowActions(false)}
        >
          <ActionButton onClick={handleReply} title="Ответить">
            <Reply size={14} />
          </ActionButton>
          {isOwn && (
            <>
              <ActionButton onClick={handleEdit} title="Редактировать">
                <Edit size={14} />
              </ActionButton>
              <ActionButton onClick={handleDelete} title="Удалить">
                <Trash2 size={14} />
              </ActionButton>
            </>
          )}
        </MessageActions>
      </MessageBubble>

      {isOwn && showAvatar && (
        <Avatar>
          {currentUser?.username?.charAt(0)?.toUpperCase() || 'U'}
        </Avatar>
      )}
      
      {!isOwn && showAvatar && (
        <div style={{ width: '32px' }} />
      )}
    </MessageContainer>
  );
};

export default Message;
