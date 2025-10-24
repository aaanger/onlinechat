import { useState, useEffect, useRef, useCallback } from 'react';
import { useAuth } from '../contexts/AuthContext';

const useWebSocket = (chatId, onMessage) => {
  const [socket, setSocket] = useState(null);
  const [isConnected, setIsConnected] = useState(false);
  const [error, setError] = useState(null);
  const { token } = useAuth();
  const reconnectTimeoutRef = useRef(null);
  const reconnectAttempts = useRef(0);
  const maxReconnectAttempts = 5;

  const connect = useCallback(() => {
    if (!chatId || !token) return;

    // Close existing connection
    if (socket) {
      socket.close();
    }

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.hostname}:8080/chats/${chatId}/ws`;
    
    try {
      const ws = new WebSocket(`${wsUrl}?token=${encodeURIComponent(token)}`);

      ws.onopen = () => {
        console.log('WebSocket connected');
        setIsConnected(true);
        setError(null);
        reconnectAttempts.current = 0;
      };

      ws.onmessage = (event) => {
        console.log('WebSocket message received:', event.data);
        if (onMessage) {
          onMessage(event);
        }
      };

      ws.onclose = (event) => {
        console.log('WebSocket disconnected:', event.code, event.reason);
        setIsConnected(false);
        
        // Attempt to reconnect if it wasn't a manual close
        if (event.code !== 1000 && reconnectAttempts.current < maxReconnectAttempts) {
          const delay = Math.min(1000 * Math.pow(2, reconnectAttempts.current), 30000);
          reconnectAttempts.current++;
          
          reconnectTimeoutRef.current = setTimeout(() => {
            console.log(`Reconnecting... attempt ${reconnectAttempts.current}`);
            connect();
          }, delay);
        }
      };

      ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        setError('Connection error');
        setIsConnected(false);
      };

      setSocket(ws);
    } catch (err) {
      console.error('Failed to create WebSocket connection:', err);
      setError('Failed to connect');
    }
  }, [chatId, token]);

  const disconnect = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
      reconnectTimeoutRef.current = null;
    }
    
    if (socket) {
      socket.close(1000, 'Manual disconnect');
      setSocket(null);
    }
    setIsConnected(false);
  }, [socket]);

  const sendMessage = useCallback((messageData) => {
    if (socket && isConnected) {
      try {
        socket.send(JSON.stringify(messageData));
        return true;
      } catch (err) {
        console.error('Failed to send message:', err);
        return false;
      }
    }
    return false;
  }, [socket, isConnected]);

  // Connect when chatId or token changes
  useEffect(() => {
    connect();
    return () => disconnect();
  }, [connect, disconnect]);

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
    };
  }, []);

  return {
    socket,
    isConnected,
    error,
    sendMessage,
    connect,
    disconnect
  };
};

export { useWebSocket };