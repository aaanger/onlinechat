import React, { useState, useRef, useEffect } from 'react';
import styled from 'styled-components';
import { Send, Smile, Paperclip, X } from 'lucide-react';
import { useChat } from '../../contexts/ChatContext';

const InputContainer = styled.div`
  padding: 16px 24px;
  background: white;
  border-top: 1px solid var(--border-color);
  display: flex;
  align-items: flex-end;
  gap: 12px;
`;

const InputWrapper = styled.div`
  flex: 1;
  position: relative;
  background: var(--background-color);
  border: 2px solid var(--border-color);
  border-radius: 24px;
  display: flex;
  align-items: flex-end;
  min-height: 48px;
  max-height: 120px;
  transition: all 0.2s ease;
  
  &:focus-within {
    border-color: var(--primary-color);
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
  }
`;

const TextArea = styled.textarea`
  flex: 1;
  padding: 12px 16px;
  background: none;
  border: none;
  outline: none;
  resize: none;
  font-family: inherit;
  font-size: 16px;
  color: var(--text-primary);
  line-height: 1.4;
  
  &::placeholder {
    color: var(--text-muted);
  }
  
  &::-webkit-scrollbar {
    width: 4px;
  }
  
  &::-webkit-scrollbar-track {
    background: transparent;
  }
  
  &::-webkit-scrollbar-thumb {
    background: var(--border-color);
    border-radius: 2px;
  }
`;

const InputActions = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px;
`;

const ActionButton = styled.button`
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: none;
  border: none;
  color: var(--text-muted);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  
  &:hover {
    background: var(--border-color);
    color: var(--text-secondary);
  }
  
  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
`;

const SendButton = styled(ActionButton)`
  background: var(--primary-color);
  color: white;
  
  &:hover:not(:disabled) {
    background: var(--primary-hover);
  }
`;

const ReplyBar = styled.div`
  position: absolute;
  bottom: 100%;
  left: 0;
  right: 0;
  background: var(--background-color);
  border: 1px solid var(--border-color);
  border-bottom: none;
  border-radius: 12px 12px 0 0;
  padding: 12px 16px;
  display: flex;
  align-items: center;
  gap: 12px;
`;

const ReplyContent = styled.div`
  flex: 1;
  font-size: 14px;
  color: var(--text-secondary);
`;

const ReplyText = styled.div`
  color: var(--text-primary);
  margin-top: 2px;
`;

const MessageInput = ({ disabled = false, placeholder = "Напишите сообщение..." }) => {
  const [message, setMessage] = useState('');
  const [replyTo, setReplyTo] = useState(null);
  const textareaRef = useRef(null);
  const { sendMessage, currentChat } = useChat();

  useEffect(() => {
    if (textareaRef.current) {
      textareaRef.current.style.height = 'auto';
      textareaRef.current.style.height = `${Math.min(textareaRef.current.scrollHeight, 120)}px`;
    }
  }, [message]);

  const handleSend = async () => {
    if (!message.trim() || disabled || !currentChat) return;

    const result = await sendMessage(message.trim(), 'text', replyTo?.id);
    
    if (result.success) {
      setMessage('');
      setReplyTo(null);
    }
  };

  const handleKeyDown = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  const handleFileUpload = () => {
    // TODO: Implement file upload
    console.log('File upload not implemented yet');
  };

  const handleEmojiClick = () => {
    // TODO: Implement emoji picker
    console.log('Emoji picker not implemented yet');
  };

  const canSend = message.trim() && !disabled && currentChat;

  return (
    <InputContainer>
      <InputWrapper>
        {replyTo && (
          <ReplyBar>
            <div>
              <div>Ответ на сообщение от {replyTo.username}:</div>
              <ReplyText>{replyTo.content}</ReplyText>
            </div>
            <ActionButton onClick={() => setReplyTo(null)}>
              <X size={16} />
            </ActionButton>
          </ReplyBar>
        )}
        
        <TextArea
          ref={textareaRef}
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder={placeholder}
          disabled={disabled}
          rows={1}
        />
        
        <InputActions>
          <ActionButton onClick={handleEmojiClick} title="Эмодзи">
            <Smile size={18} />
          </ActionButton>
          
          <ActionButton onClick={handleFileUpload} title="Приложить файл">
            <Paperclip size={18} />
          </ActionButton>
        </InputActions>
      </InputWrapper>

      <SendButton 
        onClick={handleSend} 
        disabled={!canSend}
        title="Отправить сообщение"
      >
        <Send size={18} />
      </SendButton>
    </InputContainer>
  );
};

export default MessageInput;
