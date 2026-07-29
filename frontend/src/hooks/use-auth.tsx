'use client';

import * as React from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '@/lib/api-client';
import { storeTokens, clearTokens, getAccessToken } from '@/lib/auth';
import type {
  User,
  AuthResponse,
  LoginRequest,
  RegisterRequest,
  APIError,
} from '@/types/api';

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
  const queryClient = useQueryClient();

  const { data: user, isLoading } = useQuery<User | null>({
    queryKey: ['current-user'],
    queryFn: async () => {
      const token = getAccessToken();
      if (!token) return null;
      try {
        return await apiClient.get<User>('/auth/me');
      } catch {
        clearTokens();
        return null;
      }
    },
    staleTime: 5 * 60 * 1000,
    retry: false,
  });

  const loginMutation = useMutation<AuthResponse, APIError, LoginRequest>({
    mutationFn: (data) => apiClient.post('/auth/login', data),
    onSuccess: (response) => {
      storeTokens(response.access_token, response.refresh_token);
      queryClient.invalidateQueries({ queryKey: ['current-user'] });
    },
  });

  const registerMutation = useMutation<AuthResponse, APIError, RegisterRequest>({
    mutationFn: (data) => apiClient.post('/auth/register', data),
    onSuccess: (response) => {
      storeTokens(response.access_token, response.refresh_token);
      queryClient.invalidateQueries({ queryKey: ['current-user'] });
    },
  });

  const logout = React.useCallback(() => {
    clearTokens();
    queryClient.setQueryData(['current-user'], null);
    queryClient.clear();
    window.location.href = '/login';
  }, [queryClient]);

  const value: AuthContextValue = {
    user: user ?? null,
    isLoading,
    isAuthenticated: !!user,
    login: async (data) => {
      await loginMutation.mutateAsync(data);
    },
    register: async (data) => {
      await registerMutation.mutateAsync(data);
    },
    logout,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextValue {
  const context = React.useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}

export function useLoginMutation() {
  return useMutation<AuthResponse, APIError, LoginRequest>({
    mutationFn: (data) => apiClient.post('/auth/login', data),
  });
}

export function useRegisterMutation() {
  return useMutation<AuthResponse, APIError, RegisterRequest>({
    mutationFn: (data) => apiClient.post('/auth/register', data),
  });
}

export function useCurrentUser() {
  return useQuery<User | null>({
    queryKey: ['current-user'],
    queryFn: async () => {
      const token = getAccessToken();
      if (!token) return null;
      try {
        return await apiClient.get<User>('/auth/me');
      } catch {
        clearTokens();
        return null;
      }
    },
    staleTime: 5 * 60 * 1000,
    retry: false,
  });
}
