import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import styled from 'styled-components';
import { Mail, Lock, User, Eye, EyeOff } from 'lucide-react';
import { useAuth } from '../../contexts/AuthContext';
import toast from 'react-hot-toast';

const RegisterContainer = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
`;

const RegisterForm = styled.form`
  background: white;
  padding: 40px;
  border-radius: 16px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 400px;
`;

const Title = styled.h1`
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary);
  text-align: center;
  margin-bottom: 8px;
`;

const Subtitle = styled.p`
  color: var(--text-secondary);
  text-align: center;
  margin-bottom: 32px;
`;

const InputGroup = styled.div`
  position: relative;
  margin-bottom: 24px;
`;

const Input = styled.input`
  width: 100%;
  padding: 16px 16px 16px 48px;
  border: 2px solid var(--border-color);
  border-radius: 12px;
  font-size: 16px;
  transition: all 0.2s ease;
  
  &:focus {
    border-color: var(--primary-color);
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
  }
  
  &::placeholder {
    color: var(--text-muted);
  }
`;

const InputIcon = styled.div`
  position: absolute;
  left: 16px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--text-muted);
`;

const PasswordToggle = styled.button`
  position: absolute;
  right: 16px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--text-muted);
  background: none;
  border: none;
  
  &:hover {
    color: var(--text-secondary);
  }
`;

const Button = styled.button`
  width: 100%;
  padding: 16px;
  background: var(--primary-color);
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-bottom: 24px;
  
  &:hover {
    background: var(--primary-hover);
    transform: translateY(-1px);
  }
  
  &:disabled {
    background: var(--text-muted);
    cursor: not-allowed;
    transform: none;
  }
`;

const LinkText = styled.div`
  text-align: center;
  color: var(--text-secondary);
  
  a {
    color: var(--primary-color);
    text-decoration: none;
    font-weight: 500;
    
    &:hover {
      text-decoration: underline;
    }
  }
`;

const PasswordRequirements = styled.div`
  font-size: 12px;
  color: var(--text-muted);
  margin-top: 8px;
  line-height: 1.4;
`;

const Register = () => {
  const [formData, setFormData] = useState({
    email: '',
    username: '',
    password: '',
  });
  const [showPassword, setShowPassword] = useState(false);
  const [loading, setLoading] = useState(false);
  const { register } = useAuth();

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (formData.password.length < 8) {
      toast.error('Пароль должен содержать минимум 8 символов');
      return;
    }

    if (formData.username.length < 3) {
      toast.error('Имя пользователя должно содержать минимум 3 символа');
      return;
    }

    setLoading(true);

    const result = await register(formData.email, formData.username, formData.password);
    
    if (result.success) {
      toast.success('Регистрация успешна!');
    } else {
      toast.error(result.error);
    }
    
    setLoading(false);
  };

  return (
    <RegisterContainer>
      <RegisterForm onSubmit={handleSubmit}>
        <Title>Регистрация</Title>
        <Subtitle>Создайте новый аккаунт</Subtitle>
        
        <InputGroup>
          <InputIcon>
            <Mail size={20} />
          </InputIcon>
          <Input
            type="email"
            name="email"
            placeholder="Email"
            value={formData.email}
            onChange={handleChange}
            required
          />
        </InputGroup>

        <InputGroup>
          <InputIcon>
            <User size={20} />
          </InputIcon>
          <Input
            type="text"
            name="username"
            placeholder="Имя пользователя"
            value={formData.username}
            onChange={handleChange}
            required
          />
        </InputGroup>

        <InputGroup>
          <InputIcon>
            <Lock size={20} />
          </InputIcon>
          <Input
            type={showPassword ? 'text' : 'password'}
            name="password"
            placeholder="Пароль"
            value={formData.password}
            onChange={handleChange}
            required
          />
          <PasswordToggle
            type="button"
            onClick={() => setShowPassword(!showPassword)}
          >
            {showPassword ? <EyeOff size={20} /> : <Eye size={20} />}
          </PasswordToggle>
          <PasswordRequirements>
            Минимум 8 символов
          </PasswordRequirements>
        </InputGroup>

        <Button type="submit" disabled={loading}>
          {loading ? 'Регистрация...' : 'Зарегистрироваться'}
        </Button>

        <LinkText>
          Уже есть аккаунт? <Link to="/login">Войти</Link>
        </LinkText>
      </RegisterForm>
    </RegisterContainer>
  );
};

export default Register;
