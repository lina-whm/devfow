'use client';

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '@/lib/api-client';
import type { Notification, PaginatedResponse, APIError } from '@/types/api';

export function useNotifications(page = 1, pageSize = 20) {
  return useQuery<PaginatedResponse<Notification>>({
    queryKey: ['notifications', page, pageSize],
    queryFn: () => apiClient.get(`/notifications?page=${page}&pageSize=${pageSize}`),
  });
}

export function useUnreadCount() {
  return useQuery<number>({
    queryKey: ['notifications', 'unread-count'],
    queryFn: () => apiClient.get('/notifications/unread-count'),
    refetchInterval: 30 * 1000,
  });
}

export function useMarkAsRead() {
  const queryClient = useQueryClient();

  return useMutation<void, APIError, string>({
    mutationFn: (notificationId) =>
      apiClient.patch(`/notifications/${notificationId}/read`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notifications'] });
    },
  });
}

export function useMarkAllAsRead() {
  const queryClient = useQueryClient();

  return useMutation<void, APIError, void>({
    mutationFn: () => apiClient.patch('/notifications/read-all'),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notifications'] });
    },
  });
}
