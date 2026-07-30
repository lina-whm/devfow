'use client';

import * as React from 'react';
import { apiClient } from '@/lib/api-client';
import { storeTokens, clearTokens, getAccessToken } from '@/lib/auth';
import type { User, AuthResponse, LoginRequest, RegisterRequest, APIError } from '@/types/api';

interface AuthContextValue {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (data: LoginRequest) => Promise<void>;
  register: (data: RegisterRequest) => Promise<void>;
  logout: () => void;
}

const AuthContext = React.createContext<AuthContextValue | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = React.useState<User | null>(null);
  const [isLoading, setIsLoading] = React.useState(true);

  React.useEffect(() => {
    const token = getAccessToken();
    console.debug(`[AuthProvider] mount token=${token ? token.slice(0,20)+'...' : 'null'}`);
    if (!token) {
      setUser(null);
      setIsLoading(false);
      return;
    }
    apiClient.get<User>('/auth/me').then((u) => {
      console.debug('[AuthProvider] /auth/me OK', u.email);
      setUser(u);
    }).catch((e) => {
      console.debug('[AuthProvider] /auth/me FAILED', e?.response?.status, e?.message);
      clearTokens();
      setUser(null);
    }).finally(() => setIsLoading(false));
  }, []);

  const login = async (data: LoginRequest) => {
    const response = await apiClient.post<AuthResponse>('/auth/login', data);
    storeTokens(response.accessToken, response.refreshToken);
    const userData = await apiClient.get<User>('/auth/me');
    setUser(userData);
  };

  const register = async (data: RegisterRequest) => {
    const response = await apiClient.post<AuthResponse>('/auth/register', data);
    storeTokens(response.accessToken, response.refreshToken);
    const userData = await apiClient.get<User>('/auth/me');
    setUser(userData);
  };

  const logout = React.useCallback(() => {
    clearTokens();
    setUser(null);
  }, []);

  return (
    <AuthContext.Provider
      value={{ user, isLoading, isAuthenticated: !!user, login, register, logout }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth(): AuthContextValue {
  const context = React.useContext(AuthContext);
  if (!context) throw new Error('useAuth must be used within an AuthProvider');
  return context;
}
