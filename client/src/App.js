import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { Toaster } from 'react-hot-toast';
import styled from 'styled-components';

import { AuthProvider, useAuth } from './contexts/AuthContext';
import { ChatProvider } from './contexts/ChatContext';
import Login from './components/auth/Login';
import Register from './components/auth/Register';
import ChatLayout from './components/layout/ChatLayout';
import ProtectedRoute from './components/auth/ProtectedRoute';

const AppContainer = styled.div`
  height: 100vh;
  width: 100vw;
  display: flex;
  flex-direction: column;
`;

function AppRoutes() {
  const { user, loading } = useAuth();

  if (loading) {
    return (
      <LoadingContainer>
        <LoadingSpinner>
          <div></div>
          <div></div>
          <div></div>
        </LoadingSpinner>
        <p>Загрузка...</p>
      </LoadingContainer>
    );
  }

  return (
    <Router>
      <Routes>
        <Route 
          path="/login" 
          element={user ? <Navigate to="/" replace /> : <Login />} 
        />
        <Route 
          path="/register" 
          element={user ? <Navigate to="/" replace /> : <Register />} 
        />
        <Route 
          path="/" 
          element={
            <ProtectedRoute>
              <ChatProvider>
                <ChatLayout />
              </ChatProvider>
            </ProtectedRoute>
          } 
        />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </Router>
  );
}

const LoadingContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
  gap: 16px;
`;

const LoadingSpinner = styled.div`
  display: inline-block;
  position: relative;
  width: 64px;
  height: 64px;

  div {
    box-sizing: border-box;
    display: block;
    position: absolute;
    width: 51px;
    height: 51px;
    margin: 6px;
    border: 6px solid var(--primary-color);
    border-radius: 50%;
    animation: spin 1.2s cubic-bezier(0.5, 0, 0.5, 1) infinite;
    border-color: var(--primary-color) transparent transparent transparent;
  }

  div:nth-child(1) { animation-delay: -0.45s; }
  div:nth-child(2) { animation-delay: -0.3s; }
  div:nth-child(3) { animation-delay: -0.15s; }

  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }
`;

function App() {
  return (
    <AppContainer>
      <AuthProvider>
        <AppRoutes />
        <Toaster 
          position="top-right"
          toastOptions={{
            duration: 4000,
            style: {
              background: '#363636',
              color: '#fff',
              borderRadius: '8px',
              padding: '16px',
            },
            success: {
              iconTheme: {
                primary: '#10b981',
                secondary: '#fff',
              },
            },
            error: {
              iconTheme: {
                primary: '#ef4444',
                secondary: '#fff',
              },
            },
          }}
        />
      </AuthProvider>
    </AppContainer>
  );
}

export default App;
