import React, { useState } from 'react';
import styled from 'styled-components';
import { X, MessageCircle, Lock, Users } from 'lucide-react';
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
  max-width: 480px;
  max-height: 90vh;
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

const Form = styled.form`
  display: flex;
  flex-direction: column;
  gap: 20px;
`;

const InputGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const Label = styled.label`
  font-weight: 500;
  color: var(--text-primary);
  font-size: 14px;
`;

const Input = styled.input`
  padding: 12px 16px;
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

const TextArea = styled.textarea`
  padding: 12px 16px;
  border: 2px solid var(--border-color);
  border-radius: 8px;
  font-size: 16px;
  font-family: inherit;
  resize: vertical;
  min-height: 80px;
  transition: all 0.2s ease;
  
  &:focus {
    outline: none;
    border-color: var(--primary-color);
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
  }
`;

const CheckboxGroup = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  border: 2px solid var(--border-color);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  
  &:hover {
    border-color: var(--primary-color);
  }
`;

const Checkbox = styled.input`
  width: 18px;
  height: 18px;
  accent-color: var(--primary-color);
`;

const CheckboxLabel = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
  color: var(--text-primary);
`;

const MemberLimitGroup = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
`;

const MemberLimitInput = styled(Input)`
  width: 120px;
`;

const ButtonGroup = styled.div`
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  margin-top: 8px;
`;

const Button = styled.button`
  padding: 12px 24px;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  
  ${props => props.primary ? `
    background: var(--primary-color);
    color: white;
    border: none;
    
    &:hover {
      background: var(--primary-hover);
    }
    
    &:disabled {
      background: var(--text-muted);
      cursor: not-allowed;
    }
  ` : `
    background: white;
    color: var(--text-secondary);
    border: 2px solid var(--border-color);
    
    &:hover {
      border-color: var(--text-secondary);
      color: var(--text-primary);
    }
  `}
`;

const CreateChatModal = ({ onClose }) => {
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    is_private: false,
    max_members: 100
  });
  const [loading, setLoading] = useState(false);
  const { createChat } = useChat();

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!formData.name.trim()) {
      toast.error('Название чата не может быть пустым');
      return;
    }

    setLoading(true);

    const result = await createChat({
      name: formData.name.trim(),
      description: formData.description.trim(),
      is_private: formData.is_private,
      max_members: parseInt(formData.max_members) || 100
    });

    if (result.success) {
      toast.success('Чат создан успешно!');
      onClose();
    } else {
      toast.error(result.error);
    }

    setLoading(false);
  };

  return (
    <ModalOverlay onClick={onClose}>
      <Modal onClick={(e) => e.stopPropagation()}>
        <ModalHeader>
          <Title>
            <MessageCircle size={20} />
            Создать новый чат
          </Title>
          <CloseButton onClick={onClose}>
            <X size={20} />
          </CloseButton>
        </ModalHeader>

        <Form onSubmit={handleSubmit}>
          <InputGroup>
            <Label htmlFor="name">Название чата</Label>
            <Input
              id="name"
              name="name"
              type="text"
              value={formData.name}
              onChange={handleChange}
              placeholder="Введите название чата"
              maxLength={100}
              required
            />
          </InputGroup>

          <InputGroup>
            <Label htmlFor="description">Описание (необязательно)</Label>
            <TextArea
              id="description"
              name="description"
              value={formData.description}
              onChange={handleChange}
              placeholder="Краткое описание чата"
              maxLength={500}
            />
          </InputGroup>

          <CheckboxGroup onClick={() => setFormData(prev => ({ ...prev, is_private: !prev.is_private }))}>
            <Checkbox
              name="is_private"
              checked={formData.is_private}
              onChange={handleChange}
              type="checkbox"
            />
            <CheckboxLabel>
              <Lock size={16} />
              Приватный чат
            </CheckboxLabel>
          </CheckboxGroup>

          <InputGroup>
            <Label htmlFor="max_members">Максимум участников</Label>
            <MemberLimitGroup>
              <MemberLimitInput
                id="max_members"
                name="max_members"
                type="number"
                value={formData.max_members}
                onChange={handleChange}
                min="2"
                max="1000"
              />
              <Users size={16} color="var(--text-muted)" />
            </MemberLimitGroup>
          </InputGroup>

          <ButtonGroup>
            <Button type="button" onClick={onClose}>
              Отменить
            </Button>
            <Button primary type="submit" disabled={loading}>
              {loading ? 'Создание...' : 'Создать чат'}
            </Button>
          </ButtonGroup>
        </Form>
      </Modal>
    </ModalOverlay>
  );
};

export default CreateChatModal;
